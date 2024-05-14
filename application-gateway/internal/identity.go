package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/pkg/errors"
	"time"
)

// GetUserId returns a unique ID associated with the invoking identity.
func (a *Application) GetUserId() (string, error) {
	return getUserId(a.CertPath)
}

// RegisterAsPatient This type of transaction would typically only be run once by an application the first time it was started after its
// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func (a *Application) RegisterAsPatient(userName, dob string) error {
	contract := a.network.GetContract(healthContract)
	fmt.Printf("\n--> Submit Transaction: RegisterAsPatient, function registers himself on the ledger \n")

	_, err := contract.SubmitTransaction("RegisterAsPatient", userName, dob)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	fmt.Printf("*** Transaction committed successfully\n")
	return nil
}

func getUserId(certPath string) (string, error) {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		return "", err
	}

	id := fmt.Sprintf("x509::%s::%s", getDN(&certificate.Subject), getDN(&certificate.Issuer))
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}

// Evaluate a transaction to query ledger state.

func (a *Application) ListenBlockEvents() error {
	blocks, err := a.network.BlockEvents(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("--- Listening for events ---")

	go func() {
		for block := range blocks {
			ct := time.Now().UTC()
			envelope, envPErr := GetEnvelopeFromBlock(block.Data.Data[0])
			if envPErr != nil {
				fmt.Println("error: ", envPErr)
				return
			}

			payload := &common.Payload{}
			err = proto.Unmarshal(envelope.Payload, payload)
			if err != nil {
				fmt.Println(err, "unmarshaling Payload error: ")
				return
			}

			channelHeader := &common.ChannelHeader{}
			err := proto.Unmarshal(payload.Header.ChannelHeader, channelHeader)
			if err != nil {
				fmt.Println("unmarshaling Channel Header error: ")
				return
			}

			t := time.Unix(channelHeader.Timestamp.Seconds, int64(channelHeader.Timestamp.Nanos)).UTC()
			fmt.Printf("Now: %v, Then %v\n", ct.Format(time.RFC3339), t.Format(time.RFC3339))
		}

	}()

	time.Sleep(time.Minute * 1000)

	return nil
}

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {

	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	return env, nil
}

package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/pkg/errors"
	"math/rand"
	"os"
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

func (a *Application) CreateDummyDiagnosis(userName string) error {
	fmt.Println("User: ", userName)
	if userName == "" {
		return fmt.Errorf("invalid user name")
	}

	// Read the JSON file
	data, err := os.ReadFile("../application-gateway/diagnosis.json")
	if err != nil {
		return err
	}

	var diagnosisList []map[string]interface{}

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(data, &diagnosisList)
	if err != nil {
		return err
	}

	userId, err := getUserId(a.CertPath)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	diagnosis := diagnosisList[rand.Intn(len(diagnosisList))]

	jsonData, err := json.Marshal(diagnosis)
	if err != nil {
		return fmt.Errorf("Error marshaling to JSON:", err)
	}

	// Print the JSON string
	//str := string(jsonData)
	//dataStr := strings.Replace(str, `"`, `\"`, -1)

	contract := a.network.GetContract(healthContract)
	fmt.Printf("\n--> Submit Transaction: CreateDummyDiagnosis, creates diagnosis %s\n", time.Now().UTC().Format(time.RFC3339))

	_, err = contract.SubmitTransaction("CreateDiagnosis", userId, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	fmt.Printf("*** Transaction committed successfully\n")
	return nil
}

func (a *Application) InsertDiagnosisFromPimaDiabetesDataset(
	pregnancies,
	glucose,
	bloodPressure,
	skinThickness,
	insulin,
	BMI int64,
	DiabetesPedigreeFunction float64,
	Age,
	Outcome int64,
) error {
	userId, err := getUserId(a.CertPath)
	if err != nil {
		return err
	}
	fmt.Println(userId)

	rand.Seed(time.Now().UnixNano())

	var diagnosis = map[string]interface{}{
		"pregnancies":                pregnancies,
		"glucose":                    glucose,
		"blood_pressure":             bloodPressure,
		"skin_thickness":             skinThickness,
		"insulin":                    insulin,
		"bmi":                        BMI,
		"diabetes_pedigree_function": DiabetesPedigreeFunction,
		"age":                        Age,
		"outcome":                    Outcome,
	}

	jsonData, err := json.Marshal(diagnosis)
	if err != nil {
		return fmt.Errorf("Error marshaling to JSON:", err)
	}

	contract := a.network.GetContract(healthContract)
	fmt.Printf("\n--> Submit Transaction: CreateDummyDiagnosis, creates diagnosis \n")

	_, err = contract.SubmitTransaction("CreateDiagnosis", userId, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	fmt.Printf("*** Transaction committed successfully\n")
	return nil
}

func (a *Application) ReadUserData(userName string) error {
	fmt.Println(userName)
	if userName == "" {
		return fmt.Errorf("invalid user name")
	}

	userId, err := getUserId(a.CertPath)
	if err != nil {
		return err
	}

	contract := a.network.GetContract(healthContract)

	evaluateResult, err := contract.EvaluateTransaction("ReadUserData", userId)
	if err != nil {
		return fmt.Errorf("failed to call EvaluateTransaction: %w", err)
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
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
func readUserData(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("ReadUserData", "")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

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

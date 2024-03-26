package internal

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"time"
)

const (
	mspID         = "Org1MSP"
	cryptoPath    = "../test-network/organizations/peerOrganizations/org1.example.com"
	alphaCertPath = cryptoPath + "/users/alpha@org1.example.com/msp/signcerts/cert.pem"
	alphaKeyPath  = cryptoPath + "/users/alpha@org1.example.com/msp/keystore/"

	betaCertPath = cryptoPath + "/users/beta@org1.example.com/msp/signcerts/cert.pem"
	betaKeyPath  = cryptoPath + "/users/beta@org1.example.com/msp/keystore/"

	gamaCertPath = cryptoPath + "/users/gama@org1.example.com/msp/signcerts/cert.pem"
	gamaKeyPath  = cryptoPath + "/users/gama@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"

	alphaUser = "alpha"
	betaUser  = "beta"
	gamaUser  = "gama"

	policyContract = "policy-contract"
	accessContract = "access-contract"
	healthContract = "health-contract"
)

var (
	App *Application
)

type Application struct {
	CertPath string
	KeyPath  string

	network *client.Network
}

func certFilePathByUserName(userName string) string {
	return cryptoPath + fmt.Sprintf("/users/%s@org1.example.com/msp/signcerts/cert.pem", userName)
}

func keyFilePathByUserName(userName string) string {
	return cryptoPath + fmt.Sprintf("/users/%s@org1.example.com/msp/keystore/", userName)
}

func NewApplication(userName string) (*Application, error) {
	if userName == "" {
		return nil, fmt.Errorf("invalid user name")
	}
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()

	var keyPath, certPath string

	certPath = certFilePathByUserName(userName)
	keyPath = keyFilePathByUserName(userName)

	var app Application
	app.KeyPath = keyPath
	app.CertPath = certPath

	id := newIdentity(certPath)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}

	channelName := "mychannel"
	app.network = gw.GetNetwork(channelName)

	return &app, nil
}

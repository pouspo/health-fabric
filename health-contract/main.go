package main

import (
	"github.com/pouspo/health-contract/chaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	healthChaincode, err := contractapi.NewChaincode(&chaincode.HealthContract{})
	if err != nil {
		log.Panicf("Error creating health-contract chaincode: %v", err)
	}

	if err := healthChaincode.Start(); err != nil {
		log.Panicf("Error starting health-contract chaincode: %v", err)
	}
}

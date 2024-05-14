package main

import (
	"github.com/pouspo/access-contract/chaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	accessChaincode, err := contractapi.NewChaincode(&chaincode.AccessContract{})
	if err != nil {
		log.Panicf("Error creating access-contract chaincode: %v", err)
	}

	if err := accessChaincode.Start(); err != nil {
		log.Panicf("Error starting access-contract chaincode: %v", err)
	}
}

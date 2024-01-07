package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pouspo/policy-contract/chaincode"
	"log"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.PolicyContract{})
	if err != nil {
		log.Panicf("Error creating policy-contract chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting policy-contract chaincode: %v", err)
	}
}

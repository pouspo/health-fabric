package chaincode

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AccessContract struct {
	contractapi.Contract
}

type Access struct {
	Allowed bool
	Read    []string
	Write   []string
}

func (a *AccessContract) Test(ctx contractapi.TransactionContextInterface, mode string) (Access, error) {
	fmt.Println(mode)

	return Access{
		Allowed: true,
		Read:    []string{"field_1", "field_2"},
		Write:   []string{},
	}, nil
}

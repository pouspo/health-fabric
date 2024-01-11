package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

type AccessContract struct {
	contractapi.Contract
}

type Access struct {
	Allowed bool     `json:"allowed"`
	Read    []string `json:"read"`
	Write   []string `json:"write"`
}

func (a *AccessContract) AccessList(ctx contractapi.TransactionContextInterface, userId string) (Access, error) {
	cid := ctx.GetClientIdentity()
	attrString, found, err := cid.GetAttributeValue("groups")
	if err != nil {
		return Access{}, err
	}

	if !found {
		return Access{
			Allowed: false,
		}, nil
	}

	fmt.Println(attrString)

	groups := strings.Split(attrString, "-")

	params := []string{"UserSpecificAccess", userId, groups[0]}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}
	response := ctx.GetStub().InvokeChaincode("policy-contract", queryArgs, "mychannel")
	if response.Status != shim.OK {
		return Access{}, fmt.Errorf("failed to query chaincode. Got error: %s", response.Message)
	}

	type AccessPolicy struct {
		Read  []string `json:"read,omitempty" metadata:",optional"`
		Write []string `json:"write,omitempty" metadata:",optional"`
	}

	var accessResp AccessPolicy
	if err := json.Unmarshal(response.Payload, &accessResp); err != nil {
		return Access{}, err
	}

	return Access{
		Allowed: true,
		Read:    accessResp.Read,
		Write:   accessResp.Write,
	}, nil
}

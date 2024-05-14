package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pouspo/access-contract/pkg"
	"strings"
)

type AccessContract struct {
	contractapi.Contract
}

type Policy struct {
	Id        string            `json:"id,omitempty"`
	PolicyMap map[string]Access `json:"policy_map,omitempty"`
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
	var policies []Policy
	var readAccessList []string
	var writeAccessList []string

	for _, group := range groups {
		params := []string{"ReadPolicy", group}
		queryArgs := make([][]byte, len(params))
		for i, arg := range params {
			queryArgs[i] = []byte(arg)
		}
		response := ctx.GetStub().InvokeChaincode("policy-contract", queryArgs, "mychannel")
		if response.Status != shim.OK {
			// return Access{}, fmt.Errorf("failed to query chaincode. Got error: %s", response.Message)
			fmt.Printf("failed to query chaincode. Got error: %s\n", response.Message)
			continue
		}

		var accessResp Policy
		if err := json.Unmarshal(response.Payload, &accessResp); err != nil {
			fmt.Printf("failed to unmarshal data into Policy Struct. Got error: %s\n", err.Error())
			continue
		}

		policies = append(policies, accessResp)
	}

	for _, policy := range policies {
		if access, ok := policy.PolicyMap[userId]; ok {
			readAccessList = append(readAccessList, access.Read...)
			writeAccessList = append(readAccessList, access.Write...)
		}
	}

	return Access{
		Allowed: true,
		Read:    pkg.UniqueStrings(readAccessList),
		Write:   pkg.UniqueStrings(writeAccessList),
	}, nil
}

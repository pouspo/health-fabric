package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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
	requestUserId, err := a.getSubmittingClientIdentity(ctx)
	if err != nil {
		return Access{}, err
	}

	groups, err := getGroupsFromId(userId)
	if err != nil {
		return Access{}, err
	}
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
		if access, ok := policy.PolicyMap[requestUserId]; ok {
			readAccessList = append(readAccessList, access.Read...)
			writeAccessList = append(writeAccessList, access.Write...)
		}
	}

	return Access{
		Allowed: true,
		Read:    UniqueStrings(readAccessList),
		Write:   UniqueStrings(writeAccessList),
	}, nil
}

func (a *AccessContract) getSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	certificate, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return "", err
	}

	groupAttr, err := getGroupAttr(certificate)
	if err != nil {
		return "", err
	}

	id := fmt.Sprintf("x509::%s::%s::%s", getDN(&certificate.Subject), getDN(&certificate.Issuer), groupAttr)
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}

func getGroupsFromId(id string) ([]string, error) {
	decodeByte, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}
	decodedString := string(decodeByte)

	targetUserAttrs := strings.Split(decodedString, "::")
	if len(targetUserAttrs) == 0 {
		return nil, errors.New("no attr found in the userid")
	}
	targetUserAttr := targetUserAttrs[len(targetUserAttrs)-1]

	groups := strings.Split(targetUserAttr, "-")
	return groups, nil
}

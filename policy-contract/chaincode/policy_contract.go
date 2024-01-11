package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

// PolicyContract provides function for managing access contract policies.
//
// Responsibility of this module is straight forward.
// - CRUD of access control data
// - Response to some basic queries, like
//   - User 'A' can access 'field' of user `B`
type PolicyContract struct {
	contractapi.Contract
}

var (
	ErrNotFound = errors.New("not found")
)

func (p *PolicyContract) ReadPolicy(ctx contractapi.TransactionContextInterface, group string) (*Policy, error) {
	fmt.Println("ReadPolicy: ", group)
	policyJSON, err := ctx.GetStub().GetState(group)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy from world state: %v", err)
	}
	if policyJSON == nil {
		return nil, ErrNotFound
	}
	var policy Policy
	err = json.Unmarshal(policyJSON, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// AddPolicy add policy to PolicyContract data. It creates a policy if not exist
func (p *PolicyContract) AddPolicy(ctx contractapi.TransactionContextInterface, group, userId, mode string, field string) error {
	policyExist, err := p.policyExist(ctx, group)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	var policy *Policy
	if policyExist {
		policy, err = p.ReadPolicy(ctx, group)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return err
		}
	} else {
		policy = &Policy{
			Id:        group,
			PolicyMap: make(map[string]Access),
		}
	}

	accessDetail, _ := policy.PolicyMap[userId]
	accessDetail.add(mode, []string{field})

	policy.PolicyMap[userId] = accessDetail
	err = policy.write(ctx)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return err
}

// RemovePolicy remove policy from PolicyContract data. It creates a policy if not exist
func (p *PolicyContract) RemovePolicy(ctx contractapi.TransactionContextInterface, group, userId, mode string, field string) error {
	policy, err := p.ReadPolicy(ctx, group)
	if err != nil {
		return err
	}

	accessDetail, _ := policy.PolicyMap[userId]
	accessDetail.remove(mode, []string{field})

	policy.PolicyMap[userId] = accessDetail
	return policy.write(ctx)
}

// UserSpecificAccess finds the access list for a particular user for the groups.
// Argument group is a string - you can pass multiple groups seperated by comma.
func (p *PolicyContract) UserSpecificAccess(ctx contractapi.TransactionContextInterface, userID string, group string) (*Access, error) {
	groups := strings.Split(group, ",")
	if len(groups) == 0 {
		return nil, nil
	}

	var policies []*Policy
	for _, group := range groups {
		policy, err := p.ReadPolicy(ctx, group)
		if err != nil && !errors.Is(err, ErrNotFound) {
			return nil, err
		}

		if policy != nil {
			policies = append(policies, policy)
		}
	}

	if len(policies) == 0 {
		return nil, nil
	}

	var access Access
	for _, policy := range policies {
		if _access, ok := policy.PolicyMap[userID]; ok {
			access.Read = append(access.Read, _access.Read...)
			access.Write = append(access.Write, _access.Write...)
		}
	}

	return &access, nil
}

// policyExist returns true when policy with given ID exists in world state
func (p *PolicyContract) policyExist(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

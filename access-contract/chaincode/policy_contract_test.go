package chaincode_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pouspo/policy-contract/chaincode"
	"github.com/pouspo/policy-contract/chaincode/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestPolicyContract(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}

	transactionContext.GetStubReturns(chaincodeStub)

	policyContract := chaincode.PolicyContract{}

	chaincodeStub.GetStateReturns(nil, nil)

	err := policyContract.AddPolicy(transactionContext, "admin", "user_1", "read", "field_1")
	require.NoError(t, err)

	policy := chaincode.Policy{
		Id: "admin",
		PolicyMap: map[string]chaincode.Access{
			"user_01": chaincode.Access{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
		},
	}

	bt, _ := json.Marshal(policy)
	chaincodeStub.GetStateReturns(bt, nil)
	p, err := policyContract.ReadPolicy(transactionContext, "admin")
	require.NoError(t, err)
	require.True(t, cmp.Equal(policy, *p))
}

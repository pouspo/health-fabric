package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Policy struct {
	Id        string            `json:"id,omitempty"`
	PolicyMap map[string]Access `json:"policy_map,omitempty"`
}

func (p *Policy) write(ctx contractapi.TransactionContextInterface) error {
	policyJSON, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(p.Id, policyJSON)
}

type Access struct {
	Read  []string `json:"read,omitempty"`
	Write []string `json:"write,omitempty" metadata:",optional"`
}

func (a *Access) add(mode string, fields []string) {
	switch mode {
	case "read":
		a.Read = append(a.Read, fields...)
	case "write":
		a.Write = append(a.Write, fields...)
	}
}

func (a *Access) remove(mode string, fields []string) {
	var mp = make(map[string]bool)

	switch mode {
	case "read":
		for _, s := range a.Read {
			mp[s] = true
		}
		a.Read = []string{}
	case "write":
		for _, s := range a.Write {
			mp[s] = true
		}
		a.Write = []string{}
	}

	for _, field := range fields {
		delete(mp, field)
	}

	switch mode {
	case "read":
		for key, _ := range mp {
			a.Read = append(a.Read, key)
		}
	case "write":
		for key, _ := range mp {
			a.Write = append(a.Write, key)
		}
	}
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

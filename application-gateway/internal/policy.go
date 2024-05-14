package internal

import (
	"encoding/json"
	"fmt"
)

func (a *Application) ReadPolicy(userName string) error {
	if userName == "" {
		return fmt.Errorf("invalid user name")
	}

	userId, err := getUserId(certFilePathByUserName(userName))
	if err != nil {
		return err
	}

	contract := a.network.GetContract(accessContract)

	evaluateResult, err := contract.EvaluateTransaction("AccessList", userId)
	if err != nil {
		return fmt.Errorf("failed to call EvaluateTransaction: %w", err)
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
	return nil
}

func (a *Application) InsertDummyPolicy() error {
	contract := a.network.GetContract(policyContract)
	// AddPolicy(ctx contractapi.TransactionContextInterface, group, userId, mode string, field string)

	alphaUserId, err := getUserId(alphaCertPath)
	if err != nil {
		return err
	}

	betaUserId, err := getUserId(betaCertPath)
	if err != nil {
		return err
	}

	gamaUserId, err := getUserId(gamaCertPath)
	if err != nil {
		return err
	}

	// ----- Alpha, Group1 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "read", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "read", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "write", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "write", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Beta, Group1 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "read", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "read", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "write", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "write", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Gama, Group1 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group1", gamaUserId, "read", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", gamaUserId, "read", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", gamaUserId, "write", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", gamaUserId, "write", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Alpha, Group2 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group2", alphaUserId, "read", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group2", alphaUserId, "read", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group2", alphaUserId, "write", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Beta, Group2 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group2", betaUserId, "read", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group2", betaUserId, "read", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group2", betaUserId, "write", "sugar_level")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	fmt.Println("Dummy policy has been set")

	type Access struct {
		Read  []string `json:"read"`
		Write []string `json:"write"`
	}

	type Policy struct {
		Id        string            `json:"id"`
		PolicyMap map[string]Access `json:"policy_map"`
	}

	p1 := Policy{
		Id: "group1",
		PolicyMap: map[string]Access{
			"alpha": {
				Read:  []string{"age", "sugar_level"},
				Write: []string{"age", "sugar_level"},
			}, "beta": {
				Read:  []string{"age", "sugar_level"},
				Write: []string{"age", "sugar_level"},
			}, "gama": {
				Read:  []string{"age", "sugar_level"},
				Write: []string{"age", "sugar_level"},
			},
		},
	}

	p2 := Policy{
		Id: "group2",
		PolicyMap: map[string]Access{
			"alpha": {
				Read:  []string{"age", "sugar_level"},
				Write: []string{"sugar_level"},
			}, "beta": {
				Read:  []string{"age", "sugar_level"},
				Write: []string{"sugar_level"},
			},
		},
	}

	fmt.Println("Policy1")

	p1bt, _ := json.Marshal(p1)
	fmt.Println(formatJSON(p1bt))

	p2bt, _ := json.Marshal(p2)
	fmt.Println(formatJSON(p2bt))

	return nil
}

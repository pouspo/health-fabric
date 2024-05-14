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
	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "read", "glucose")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "read", "bmi")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "read", "age")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "write", "glucose")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", alphaUserId, "write", "bmi")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Beta, Group1 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "read", "glucose")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "read", "bmi")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "write", "glucose")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", betaUserId, "write", "bmi")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Gama, Group1 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group1", gamaUserId, "read", "glucose")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group1", gamaUserId, "read", "bmi")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Alpha, Group2 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group2", alphaUserId, "read", "skin_thickness")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group2", alphaUserId, "write", "skin_thickness")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// ----- Beta, Group2 ----- //
	_, err = contract.SubmitTransaction("AddPolicy", "group2", betaUserId, "read", "skin_thickness")
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	_, err = contract.SubmitTransaction("AddPolicy", "group2", betaUserId, "write", "skin_thickness")
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
				Read:  []string{"glucose", "bmi", "age"},
				Write: []string{"glucose", "bmi"},
			}, "beta": {
				Read:  []string{"glucose", "bmi"},
				Write: []string{"glucose", "bmi"},
			}, "gama": {
				Read:  []string{"glucose", "bmi"},
				Write: []string{},
			},
		},
	}

	p2 := Policy{
		Id: "group2",
		PolicyMap: map[string]Access{
			"alpha": {
				Read:  []string{"skin_thickness"},
				Write: []string{"skin_thickness"},
			}, "beta": {
				Read:  []string{"skin_thickness"},
				Write: []string{"skin_thickness"},
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

package internal

import "fmt"

func (a *Application) ReadUserData(userName string) error {
	if userName == "" {
		return fmt.Errorf("invalid user name")
	}

	userId, err := getUserId(certFilePathByUserName(userName))
	if err != nil {
		return err
	}

	contract := a.network.GetContract(healthContract)

	evaluateResult, err := contract.EvaluateTransaction("ReadUserData", userId)
	if err != nil {
		return fmt.Errorf("failed to call EvaluateTransaction: %w", err)
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
	return nil
}

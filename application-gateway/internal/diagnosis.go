package internal

import (
	"encoding/json"
	"fmt"
)

func (a *Application) InsertDiagnosisData(param ...string) error {
	if len(param) <= 1 {
		return fmt.Errorf("invalid parameter, pass them as pair of key value like [key1 value1 key2 value2 ... keyn valuen]")
	}

	var (
		userId string
		err    error
		i      int
	)

	if len(param)%2 == 1 {
		userId, err = getUserId(certFilePathByUserName(param[0]))
		if err != nil {
			return err
		}
		i = 1
	} else {
		userId, err = getUserId(a.CertPath)
		if err != nil {
			return err
		}
	}

	var diagnosis = map[string]interface{}{}
	for ; i < len(param); i += 2 {
		diagnosis[param[i]] = param[i+1]
	}

	jsonData, err := json.Marshal(diagnosis)
	if err != nil {
		return fmt.Errorf("error marshaling to json: %v", err)
	}

	contract := a.network.GetContract(healthContract)
	fmt.Printf("\n--> Submit Transaction: CreateDiagnosis, creates diagnosis \n")

	_, err = contract.SubmitTransaction("CreateDiagnosis", userId, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to submit CreateDiagnosis transaction: %w", err)
	}

	fmt.Printf("*** CreateDiagnosis, Transaction committed successfully\n")

	return nil
}

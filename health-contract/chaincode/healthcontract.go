package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

// HealthContract provides functions for managing an Asset
type HealthContract struct {
	contractapi.Contract
}

type UserData struct {
	UserId         string                   `json:"user_id,omitempty"`
	UserName       string                   `json:"user_name,omitempty" metadata:",optional"`
	Dob            string                   `json:"dob,omitempty" metadata:",optional"`
	DigagnosisList []map[string]interface{} `json:"digagnosis_list,omitempty" metadata:",optional"`
}

// InitHealthDataLedger adds a base set of assets to the ledger
func (s *HealthContract) InitHealthDataLedger(ctx contractapi.TransactionContextInterface) error {
	dataList := []UserData{
		{
			UserId:   "user_1",
			UserName: "user_name_1",
			Dob:      "01-12-1997",
			DigagnosisList: []map[string]interface{}{
				{
					"created_at": time.Now().UTC().Unix(),
					"diabetes":   false,
					"bp":         "70/110",
					"cancer":     true,
					"HBsAg":      "-ve",
				},
				{
					"created_at": time.Now().UTC().Unix(),
					"diabetes":   true,
					"bp":         "90/110",
					"cancer":     true,
					"HBsAg":      "+ve",
				},
			},
		},
	}

	for _, data := range dataList {
		assetJSON, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(data.UserId, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *HealthContract) RegisterAsPatient(ctx contractapi.TransactionContextInterface, name, dob string) error {
	userId, err := s.getSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	userData := UserData{
		UserId:         userId,
		UserName:       name,
		Dob:            dob,
		DigagnosisList: nil,
	}

	assetJSON, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(userData.UserId, assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// ReadUserData returns the user data stored in the world state with given id.
func (s *HealthContract) ReadUserData(ctx contractapi.TransactionContextInterface, userId string) (*UserData, error) {
	userData, err := s.readUserData(ctx, userId)
	if err != nil {
		return nil, err
	}

	if userData == nil {
		return nil, fmt.Errorf("the user %s does not exist", userId)
	}

	return userData, nil
}

func (s *HealthContract) readUserData(ctx contractapi.TransactionContextInterface, userId string) (*UserData, error) {
	assetJSON, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	fmt.Println(assetJSON)

	if assetJSON == nil {
		return nil, nil
	}

	var userData UserData
	err = json.Unmarshal(assetJSON, &userData)
	if err != nil {
		return nil, err
	}

	return &userData, nil
}

// UserDataExists returns true when asset with given ID exists in world state
func (s *HealthContract) UserDataExists(ctx contractapi.TransactionContextInterface, userId string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// UpdateUserData updates an existing user data in the world state with provided parameters.
func (s *HealthContract) UpdateUserData(ctx contractapi.TransactionContextInterface, userData UserData) error {
	existingUserData, err := s.readUserData(ctx, userData.UserId)
	if err != nil {
		return err
	}

	if existingUserData == nil {
		return fmt.Errorf("the user %s does not exist", userData.UserId)
	}

	dataJSON, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(userData.UserId, dataJSON)
}

// CreateDiagnosis issues a new diagnosis to the world state with given details.
func (s *HealthContract) CreateDiagnosis(ctx contractapi.TransactionContextInterface, userId string, data map[string]interface{}) error {
	userData, err := s.ReadUserData(ctx, userId)
	if err != nil {
		return err
	}

	userData.DigagnosisList = append(userData.DigagnosisList, data)

	return s.UpdateUserData(ctx, *userData)
}

// GetAllUserData returns all user's data found in world state
func (s *HealthContract) GetAllUserData(ctx contractapi.TransactionContextInterface) ([]*UserData, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var userDataList []*UserData
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var userData UserData
		err = json.Unmarshal(queryResponse.Value, &userData)
		if err != nil {
			return nil, err
		}
		userDataList = append(userDataList, &userData)
	}

	return userDataList, nil
}

// GetSubmittingClientIdentity returns the name and issuer of the identity that
// invokes the smart contract. This function base64 decodes the identity string
// before returning the value to the client or smart contract.
func (s *HealthContract) getSubmittingClientIdentityFull(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := s.getSubmittingClientIdentity(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func (s *HealthContract) getSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	return ctx.GetClientIdentity().GetID()
}

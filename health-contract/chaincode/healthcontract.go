package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// HealthContract provides functions for managing an Asset
type HealthContract struct {
	contractapi.Contract
}

type Access struct {
	Allowed bool     `json:"allowed"`
	Read    []string `json:"read"`
	Write   []string `json:"write"`
}

type UserData struct {
	UserId         string                   `json:"user_id,omitempty"`
	UserName       string                   `json:"user_name,omitempty" metadata:",optional"`
	Dob            string                   `json:"dob,omitempty" metadata:",optional"`
	DigagnosisList []map[string]interface{} `json:"digagnosis_list,omitempty" metadata:",optional"`
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
// If you pass empty string, this will look for your userdata
func (s *HealthContract) ReadUserData(ctx contractapi.TransactionContextInterface, userId string) (*UserData, error) {
	requestUserId, err := s.getSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if userId == "" {
		userId = requestUserId
	}

	userData, err := s.readUserData(ctx, userId)
	if err != nil {
		return nil, err
	}

	if userData == nil {
		return nil, fmt.Errorf("user data not found")
	}

	if userId == requestUserId {
		return userData, nil
	}

	access, err := s.getAccessDetails(ctx, userId)
	if err != nil {
		return nil, err
	}

	if access == nil {
		return nil, fmt.Errorf("you do not have access to this dataset")
	}

	for _, diagnostic := range userData.DigagnosisList {
		for key, _ := range diagnostic {
			if !inArray(access.Read, key) {
				diagnostic[key] = nil
			}
		}
	}

	return userData, nil
}

func (s *HealthContract) getAccessDetails(ctx contractapi.TransactionContextInterface, userId string) (*Access, error) {
	// AccessList
	params := []string{"AccessList", userId}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}
	response := ctx.GetStub().InvokeChaincode("access-contract", queryArgs, "mychannel")
	if response.Status != shim.OK {
		return nil, fmt.Errorf("failed to query chaincode. Got error: %s", response.Payload)
	}

	var access Access
	if err := json.Unmarshal(response.Payload, &access); err != nil {
		return nil, err
	}

	if !access.Allowed {
		return nil, fmt.Errorf("you do not have access to this dataset")
	}

	return &access, nil
}

func (s *HealthContract) readUserData(ctx contractapi.TransactionContextInterface, userId string) (*UserData, error) {
	assetJSON, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

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
	dataJSON, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(userData.UserId, dataJSON)
}

// CreateDiagnosis issues a new diagnosis to the world state with given details.
func (s *HealthContract) CreateDiagnosis(ctx contractapi.TransactionContextInterface, userId string, data map[string]interface{}) error {
	fmt.Println(data)

	requestUserId, err := s.getSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if userId == "" {
		userId = requestUserId
	}

	userData, err := s.readUserData(ctx, userId)
	if err != nil {
		return err
	}

	fmt.Println("v1 ", userData)

	fmt.Println("v2 ", data)

	if userData == nil {
		return fmt.Errorf("user data not found")
	}

	if requestUserId == userId {
		userData.DigagnosisList = append(userData.DigagnosisList, data)

		return s.UpdateUserData(ctx, *userData)
	}

	access, err := s.getAccessDetails(ctx, userId)
	if err != nil {
		return err
	}

	if access != nil {
		return fmt.Errorf("you do not have access to this dataset")
	}

	fmt.Println("v3 ", access)

	for key, _ := range data {
		if !inArray(access.Write, key) {
			data[key] = nil
		}
	}

	fmt.Println("v4 ", data)

	userData.DigagnosisList = append(userData.DigagnosisList, data)

	fmt.Println("v5 ", userData)

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

func (s *HealthContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	return ctx.GetClientIdentity().GetID()
}

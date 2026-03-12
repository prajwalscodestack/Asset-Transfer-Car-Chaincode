package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract structure
type SmartContract struct {
	contractapi.Contract
}

// Car asset structure
type Car struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Color string `json:"color"`
	Owner string `json:"owner"`
}

// CreateCar adds a new car to the ledger
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface,
	id string, make string, model string, color string, owner string) error {

	car := Car{
		ID:    id,
		Make:  make,
		Model: model,
		Color: color,
		Owner: owner,
	}

	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}

// ReadCar returns car from ledger
func (s *SmartContract) ReadCar(ctx contractapi.TransactionContextInterface, id string) (*Car, error) {

	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if carJSON == nil {
		return nil, fmt.Errorf("car %s does not exist", id)
	}

	var car Car
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

// TransferCar changes the owner of the car
func (s *SmartContract) TransferCar(ctx contractapi.TransactionContextInterface,
	id string, newOwner string) error {

	car, err := s.ReadCar(ctx, id)
	if err != nil {
		return err
	}

	car.Owner = newOwner

	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}

// GetAllCars returns all cars stored in ledger
func (s *SmartContract) GetAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var cars []*Car

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car Car
		err = json.Unmarshal(queryResponse.Value, &car)
		if err != nil {
			return nil, err
		}

		cars = append(cars, &car)
	}

	return cars, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting chaincode: %v", err))
	}
}

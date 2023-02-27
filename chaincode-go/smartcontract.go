package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a poll
type SmartContract struct {
	contractapi.Contract
}

// Poll describes specified details of what makes up a poll
type Poll struct {
	ID 		        string    	`json:"ID"`
	Name          	string 		`json:"Name"`
	Researcher      string    	`json:"Researcher"`
	Description     string 		`json:"Description"`
	Status          string      `json:"Status"`
}

// InitLedger adds the live testing poll into the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	polls := []Poll{
			{ID: "poll1", Name: "e", Researcher: "e", Description: "e", Status: "e"},
		}
		
	for _, poll := range polls {
		pollJSON, err := json.Marshal(poll)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(poll.ID, pollJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// GetAllPolls returns all polls found in the world state aka current ledger.
func (s *SmartContract) GetAllPolls(ctx contractapi.TransactionContextInterface) ([]*Poll, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var polls []*Poll
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var poll Poll
		err = json.Unmarshal(queryResponse.Value, &poll)
		if err != nil {
			return nil, err
		}
		polls = append(polls, &poll)
	}

	return polls, nil
}

// Create the chaincode & start it. Catch errors.
func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create e-voting chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting e-voting chaincode: %s", err.Error())
	}
}

package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a poll.
type SmartContract struct {
	contractapi.Contract
}

// Poll describes specified details of what makes up a poll.
type Poll struct {
	ID 		        string    	`json:"ID"`
	Name          	string 		`json:"Name"`
	Researcher      string    	`json:"Researcher"`
	Description     string 		`json:"Description"`
	Status          string      `json:"Status"`
}

// InitLedgerPoll adds the live testing poll into the ledger.
func (s *SmartContract) InitLedgerPoll(ctx contractapi.TransactionContextInterface) error {
	polls := []Poll{
			{ID: "1", Name: "Does blockchain increase participation in polls for academic research?", Researcher: "UTAR", Description: "Polling is used by sociologists for academic research. \nHowever, the participation rate has decreased over the years due to lack of privacy, ease of use & accessibility. \nFrom recent research, using blockchain technology addresses these aforementioned issues. \nThis survey gathers public opinion to test this hypothesis.", Status: "Ongoing"},
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

// CreatePoll issues a new poll to the world state with given details
func (s *SmartContract) CreatePoll(ctx contractapi.TransactionContextInterface, id string, name string, researcher string, description string, status string) error {
	exists, err := s.PollExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the poll %s already exists", id)
	}

	poll := Poll{
		ID:            id,
		Name:          name,
		Researcher:    researcher,
		Description:   description,
		Status:        status,
	}
	pollJSON, err := json.Marshal(poll)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, pollJSON)
}

// ReadPoll returns the poll stored in the world state with given id.
func (s *SmartContract) ReadPoll(ctx contractapi.TransactionContextInterface, id string) (*Poll, error) {
	pollJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if pollJSON == nil {
		return nil, fmt.Errorf("the poll %s does not exist", id)
	}

	var poll Poll
	err = json.Unmarshal(pollJSON, &poll)
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

// UpdatePoll updates an existing poll in the world state with provided parameters.
func (s *SmartContract) UpdatePoll(ctx contractapi.TransactionContextInterface, id string, name string, researcher string, description string, status string) error {
	exists, err := s.PollExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the poll %s does not exist", id)
	}

	// overwriting original poll with new poll
	poll := Poll{
		ID:             id,
		Name:           name,
		Researcher:     researcher,
		Description:    description,
		Status: 		status,
	}
	pollJSON, err := json.Marshal(poll)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, pollJSON)
}

// DeletePoll deletes an given asset from the world state.
func (s *SmartContract) DeletePoll(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.PollExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the poll %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// PollExists returns true when poll with given ID exists in world state
func (s *SmartContract) PollExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	pollJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return pollJSON != nil, nil
}

// GetAllPolls returns all polls found in the world state.
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

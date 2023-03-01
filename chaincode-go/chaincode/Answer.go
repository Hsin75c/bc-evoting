package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a answer.
type SmartContract struct {
	contractapi.Contract
}

// Answer describes specified details of what makes up an answer.
type Answer struct {
	ID 		        string    	`json:"ID"`
	Answer      	string    	`json:"Answer"`
}

// InitLedgerAnswer adds answers to the live testing answer into the ledger.
func (s *SmartContract) InitLedgerAnswer(ctx contractapi.TransactionContextInterface) error {
	answers := []Answer{
			{ID: "1-1-1", Answer: "4"},
			{ID: "1-2-1", Answer: "4"},
			{ID: "1-3-1", Answer: "5"},
			{ID: "1-4-1", Answer: "4"},
			{ID: "1-5-1", Answer: "4"},
			{ID: "1-6-1", Answer: "5"},
		}
		
	for _, answer := range answers {
		answerJSON, err := json.Marshal(answer)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(answer.ID, answerJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAnswer issues a new answer to the world state with given details
func (s *SmartContract) CreateAnswer(ctx contractapi.TransactionContextInterface, id string, answer string) error {
	exists, err := s.AnswerExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the answer %s already exists", id)
	}

	answer := Answer{
		ID:            id,
		Answer:        answer,
	}
	answerJSON, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, answerJSON)
}

// Readanswer returns the answer stored in the world state with given id.
func (s *SmartContract) ReadAnswer(ctx contractapi.TransactionContextInterface, id string) (*answer, error) {
	answerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if answerJSON == nil {
		return nil, fmt.Errorf("the answer %s does not exist", id)
	}

	var answer Answer
	err = json.Unmarshal(answerJSON, &answer)
	if err != nil {
		return nil, err
	}

	return &answer, nil
}

// UpdateAnswer updates an existing answer in the world state with provided parameters.
func (s *SmartContract) UpdateAnswer(ctx contractapi.TransactionContextInterface, id string, answer string) error {
	exists, err := s.AnswerExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the answer %s does not exist", id)
	}

	// overwriting original answer with new answer
	answer := Answer{
		ID:             id,
		Answer:           answer,
	}
	answerJSON, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, answerJSON)
}

// DeleteAnswer deletes an given answer from the world state.
func (s *SmartContract) DeleteAnswer(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AnswerExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the answer %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AnswerExists returns true when answer with given ID exists in world state
func (s *SmartContract) AnswerExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	answerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return answerJSON != nil, nil
}

// GetAllAnswers returns all answers found in the world state.
func (s *SmartContract) GetAllAnswers(ctx contractapi.TransactionContextInterface) ([]*Answer, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var answers []*Answer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var answer Answer
		err = json.Unmarshal(queryResponse.Value, &answer)
		if err != nil {
			return nil, err
		}
		answers = append(answers, &answer)
	}

	return answers, nil
}

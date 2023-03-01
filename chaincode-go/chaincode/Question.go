package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a question.
type SmartContract struct {
	contractapi.Contract
}

// Question describes specified details of what makes up a question.
type Question struct {
	ID 	      	    string    	`json:"ID"`
	Question      	string    	`json:"Question"`
}

// InitLedgerQuestion adds the live testing poll questions into the ledger.
func (s *SmartContract) InitLedgerQuestion(ctx contractapi.TransactionContextInterface) error {
	questions := []Question{
			{pollID: "1-1", Question: "How likely are you to participate in polling research?"},
			{pollID: "1-2", Question: "Rate the standard of privacy compared to other polling methods."},
			{pollID: "1-3", Question: "Rate the ease of use compared to other polling methods."},
			{pollID: "1-4", Question: "Rate the accessibilty compared to other polling methods."},
			{pollID: "1-5", Question: "Do you prefer to use this blockchain app over other polling methods?"},
			{pollID: "1-6", Question: "Does this application increase the likelyhood of you participating in polling research?"},
		}
		
	for _, question := range questions {
		pollJSON, err := json.Marshal(question)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(question.ID, questionJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateQuestion issues a new question to the world state with given details
func (s *SmartContract) CreateQuestion(ctx contractapi.TransactionContextInterface, id string, question string) error {
	exists, err := s.QuestionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the question %s already exists", id)
	}

	question := Question{
		ID:        	    id,
		Question:       question,
	}
	questionJSON, err := json.Marshal(question)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(pollID, pollJSON)
}

// ReadQuestion returns the question stored in the world state with given id.
func (s *SmartContract) ReadQuestion(ctx contractapi.TransactionContextInterface, id string) (*Poll, error) {
	questionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if questionJSON == nil {
		return nil, fmt.Errorf("the question %s does not exist", id)
	}

	var question Question
	err = json.Unmarshal(questionJSON, &question)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

// UpdateQuestion updates an existing question in the world state with provided parameters.
func (s *SmartContract) UpdateQuestion(ctx contractapi.TransactionContextInterface, id string, question string) error {
	exists, err := s.QuestionExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the question %s does not exist", id)
	}

	// overwriting original question with new question
	question := Question{
		ID:             id,
		Question:       question,
	}
	questionJSON, err := json.Marshal(question)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, questionJSON)
}

// DeleteQuestion deletes an given question from the world state.
func (s *SmartContract) DeleteQuestion(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.PollExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the question %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// QuestionExists returns true when question with given ID exists in world state
func (s *SmartContract) QuestionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	questionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return questionJSON != nil, nil
}

// GetAllQuestions returns all questions found in the world state.
func (s *SmartContract) GetAllQuestions(ctx contractapi.TransactionContextInterface) ([]*Poll, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var questions []*Questions
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var question Questions
		err = json.Unmarshal(queryResponse.Value, &question)
		if err != nil {
			return nil, err
		}
		questions = append(questions, &questions)
	}

	return questions, nil
}

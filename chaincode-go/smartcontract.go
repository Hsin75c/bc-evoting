package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a poll
type SmartContract struct {
	contractapi.Contract
}

// Poll describes specified details of what makes up a poll
type Poll struct {
	poll_id 		int    		`json:"poll_id"`
	name          	string 		`json:"name"`
	researcher_id   int     	`json:"researcher_id"`
	description     string 		`json:"description"`
	startDate       string      `json:"startDate"`
	endDate			string      `json:"endDate"`
	status          string      `json:"status"`
}

// InitLedger adds the live testing poll into the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	polls := []Poll{
		Poll{poll_id: 1, 
			 name: "Does blockchain increase participation in polls for academic research?",
		 	 researcher_id: 1,
		 	 description: "Polling is used by sociologists for academic research. \nHowever, the participation rate has decreased over the years due to lack of privacy, ease of use & accessibility.  \nFrom recent research, using blockchain technology addresses these aforementioned issues.  \nThis survey gathers public opinion to test this hypothesis.",
			 startDate: "2023-01-02",
		 	 endDate: "2023-06-07",
		 	 status: "Ongoing"},
		}
		
	for _, poll := range polls {
		pollJSON, err := json.Marshal(poll)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState("Poll"+strconv.Itoa(poll.poll_id), pollJSON)
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

package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Poll struct {
	poll_id 		int    		`json:"poll_id"`
	name          	string 		`json:"name"`
	researcher_id   int     	`json:"researcher_id"`
	description     string 		`json:"description"`
	startDate       time.Time   `json:"startDate"`
	endDate			time.Time   `json:"endDate"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	polls := []Poll{
		{poll_id: 1, 
		 name: "Does blockchain increase participation in polls for academic research?",
		 researcher_id: 1,
		 description: "Polling is used by sociologists for academic research. 
				However, the participation rate has decreased over the years due to lack of privacy, ease of use & accessibility.
				From recent research, using blockchain technology addresses these aforementioned issues.
				This survey gathers public opinion to test this hypothesis.",
		 startDate: "2023-01-02",
		 endDate: "2023-06-07"},
	}

	for _, poll := range polls {
		pollJSON, err := json.Marshal(poll)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(poll.poll_id, pollJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

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
		polls = append(poll, &poll)
	}

	return polls, nil
}
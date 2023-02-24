package main

import (
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/Hsin75c/bc-evoting/chaincode-go/chaincode"
	
)

func main() {
	voteChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating the chaincode: %v", err)
	}

	if err := voteChaincode.Start(); err != nil {
		log.Panicf("Error starting the chaincode: %v", err)
	}
}
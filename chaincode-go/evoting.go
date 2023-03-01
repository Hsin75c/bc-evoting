package main

import (
	"log"

	"github.com/Hsin75c/bc-evoting/chaincode-go/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Create the chaincode & start it. Catch errors.
func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		log.Panicf("Error create e-voting chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting e-voting chaincode: %s", err.Error())
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Problem struct to hold the problem information
type Problem struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Inputs      []string `json:"inputs"`
	Outputs     []string `json:"outputs"`
}

// SimpleContract contract for handling writing and reading problem information
type SimpleContract struct {
	contractapi.Contract
}

// UploadProblem adds a new problem to the world state
func (sc *SimpleContract) UploadProblem(ctx contractapi.TransactionContextInterface, id string, description string, inputs []string, outputs []string) error {
	// Check if the problem ID already exists
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return errors.New("unable to interact with world state")
	}

	if existing != nil {
		return fmt.Errorf("cannot upload problem with ID %s. Problem already exists", id)
	}

	// Ensure that inputs and outputs have the same length
	if len(inputs) != len(outputs) {
		return fmt.Errorf("the number of inputs must match the number of outputs")
	}

	// Create a new problem object
	problem := Problem{
		ID:          id,
		Description: description,
		Inputs:      inputs,
		Outputs:     outputs,
	}

	// Convert the problem struct to JSON
	problemJSON, err := json.Marshal(problem)
	if err != nil {
		return fmt.Errorf("failed to marshal problem: %v", err)
	}

	// Save the problem information to the world state
	err = ctx.GetStub().PutState(id, problemJSON)
	if err != nil {
		return errors.New("unable to interact with world state")
	}

	return nil
}

// ViewProblem retrieves the problem by ID from the world state
func (sc *SimpleContract) ViewProblem(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	// Get the problem from the world state
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return "", errors.New("unable to interact with world state")
	}

	if existing == nil {
		return "", fmt.Errorf("problem with ID %s does not exist", id)
	}

	// Return the problem as a string (JSON format)
	return string(existing), nil
}

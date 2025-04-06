package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Problem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	PassCount   int    `json:"pass_count"`
}

type SampleIO struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// 上传题目
func (s *SmartContract) CreateProblem(ctx contractapi.TransactionContextInterface, id string, description string) error {
	problem := Problem{ID: id, Description: description, PassCount: 0}
	problemBytes, _ := json.Marshal(problem)

	// 写入公共数据
	if err := ctx.GetStub().PutState(id, problemBytes); err != nil {
		return err
	}

	// 从 transient 获取私有数据
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return err
	}
	sampleBytes := transientMap["samples"]

	// 写入私有数据
	return ctx.GetStub().PutPrivateData("collectionSamples", id, sampleBytes)
}

// 查询题目（包含私有样例）
func (s *SmartContract) QueryProblem(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	probBytes, err := ctx.GetStub().GetState(id)
	if err != nil || probBytes == nil {
		return "", fmt.Errorf("problem not found")
	}
	sampleBytes, _ := ctx.GetStub().GetPrivateData("collectionSamples", id)
	return fmt.Sprintf("Public: %s\nPrivate Samples: %s", string(probBytes), string(sampleBytes)), nil
}

// 提交结果
func (s *SmartContract) SubmitAnswer(ctx contractapi.TransactionContextInterface, id string, result string) error {
	probBytes, err := ctx.GetStub().GetState(id)
	if err != nil || probBytes == nil {
		return fmt.Errorf("problem not found")
	}
	var prob Problem
	_ = json.Unmarshal(probBytes, &prob)

	if result == "pass" {
		prob.PassCount++
	}

	newBytes, _ := json.Marshal(prob)
	return ctx.GetStub().PutState(id, newBytes)
}

func main() {
	cc, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type SampleIO struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func main() {
	// 初始化 wallet
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("failed to create wallet: %v", err)
	}

	// 配置 connection profile 文件路径
	ccpPath := filepath.Join("..", "test-network", "organizations", "peerOrganizations", "org1.example.com", "connection-org1.yaml")

	// 连接 Gateway
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("failed to get network: %v", err)
	}

	// 获取链码实例（此处名称与部署时保持一致）
	contract := network.GetContract("secured")

	// 构造样例数据
	samples := []SampleIO{
		{Input: "1 2", Output: "3"},
		{Input: "5 7", Output: "12"},
	}
	sampleBytes, err := json.Marshal(samples)
	if err != nil {
		log.Fatalf("failed to marshal samples: %v", err)
	}

	// 创建题目时带上 transient 数据（私有数据由链码内部控制）
	transient := map[string][]byte{
		"samples": sampleBytes,
	}

	// 使用 WithTransient 选项传入 transient 数据
	createTx, err := contract.CreateTransaction("CreateProblem", gateway.WithTransient(transient))
	if err != nil {
		log.Fatalf("failed to create transaction: %v", err)
	}

	fmt.Println("=== Creating problem with private samples...")
	_, err = createTx.Submit("Q001", "Add two numbers")
	if err != nil {
		log.Fatalf("failed to submit CreateProblem transaction: %v", err)
	}

	// 查询题目信息（包含私有数据）
	fmt.Println("=== Query problem...")
	result, err := contract.EvaluateTransaction("QueryProblem", "Q001")
	if err != nil {
		log.Fatalf("failed to evaluate QueryProblem transaction: %v", err)
	}
	fmt.Println(string(result))

	// 提交测试结果
	fmt.Println("=== Submit 'pass' result...")
	_, err = contract.SubmitTransaction("SubmitAnswer", "Q001", "pass")
	if err != nil {
		log.Fatalf("failed to submit SubmitAnswer transaction: %v", err)
	}
}

package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	// 修改以下路径为你的身份文件位置
	credPath := filepath.Join("..", "test-network", "organizations", "peerOrganizations", "org1.example.com", "users", "User1@org1.example.com", "msp")
	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	keyDir := filepath.Join(credPath, "keystore")

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatalf("Failed to read certificate: %v", err)
	}

	// 假设 keystore 目录下只有一个文件
	files, err := ioutil.ReadDir(keyDir)
	if err != nil || len(files) == 0 {
		log.Fatalf("Failed to read private key directory: %v", err)
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	err = wallet.Put("appUser", identity)
	if err != nil {
		log.Fatalf("Failed to put identity into wallet: %v", err)
	}

	log.Println("Successfully imported appUser into wallet")
}

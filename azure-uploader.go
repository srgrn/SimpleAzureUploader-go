package main

import (
	// "bytes"
	// "encoding/base64"
	"flag"
	"fmt"
	// "io/ioutil"
	// "math/rand"
	"os"

	"github.com/Azure/azure-sdk-for-go/storage"
)

func main() {
	accountName := flag.String("accountname", "", "Storage account Name")
	accountKey := flag.String("accountkey", "", "Storage account key")
	containerName := flag.String("containername", "", "The name of the container to upload to")
	fileName := flag.String("filename", "", "The name of the file to upload ")
	targetName := flag.String("targetname", "", "The name of the blob")
	flag.Parse()
	if *accountName == "" {
		fmt.Println("Using account name from environment")
		accountName = getEnvVarOrExit("ACCOUNT_NAME")
	}
	if *accountKey == "" {
		fmt.Println("Using account key from environment")
		accountKey = getEnvVarOrExit("ACCOUNT_KEY")
	}
	if *containerName == "" {
		fmt.Println("Using container name from environment")
		accountName = getEnvVarOrExit("CONTAINER_NAME")
	}

	if *fileName == "" {
		fmt.Printf("Missing filename")
		os.Exit(1)
	}
	client, _ := storage.NewBasicClient(*accountName, *accountKey)

	blobClinet := client.GetBlobService()

	container := blobClinet.GetContainerReference(*containerName)

	file, _ := os.Open(*fileName)
	if *targetName == "" {
		targetName = fileName
	}
	blob := container.GetBlobReference(*targetName)
	fmt.Println("Start uploading file")
	blob.CreateBlockBlobFromReader(file, nil)
	fmt.Println("Done uploading file")
}

// getEnvVarOrExit returns the value of specified environment variable or terminates if it's not defined.
func getEnvVarOrExit(varName string) *string {
	value := os.Getenv(varName)
	if value == "" {
		fmt.Printf("Missing environment variable %s\n", varName)
		fmt.Println("Set environment variable or specify on CLI")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return &value
}

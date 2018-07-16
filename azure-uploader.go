package main

import (
	// "bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	// "math/rand"
	"github.com/Azure/azure-sdk-for-go/storage"
	"os"
	// "strings"
	"github.com/lithammer/shortuuid"
)

const (
	KILO = 1024
	MEGA = KILO * 1024
)

func main() {
	accountName := flag.String("accountname", "", "Storage account Name")
	accountKey := flag.String("accountkey", "", "Storage account key")
	containerName := flag.String("containername", "", "The name of the container to upload to")
	fileName := flag.String("filename", "", "The name of the file to upload ")
	targetName := flag.String("targetname", "", "The name of the blob")
	contentType := flag.String("contenttype", "application/octet-stream", "The name of the blob")
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
		containerName = getEnvVarOrExit("CONTAINER_NAME")
	}

	if *fileName == "" {
		fmt.Printf("Missing filename")
		os.Exit(1)
	}
	client, err := storage.NewBasicClient(*accountName, *accountKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	blobClinet := client.GetBlobService()

	container := blobClinet.GetContainerReference(*containerName)
	// handle single file
	err = HandleSingleFile(fileName, targetName, contentType, container)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func HandleSingleFile(fileName, targetName, contentType *string, container *storage.Container) error {
	file, _ := os.Open(*fileName)
	// test for file size
	defer file.Close()
	if *targetName == "" {
		targetName = fileName
	}

	blob := container.GetBlobReference(*targetName)
	blob.Properties.ContentType = *contentType

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	// if bigger than 250mb should upload by blocks
	var bytes int64
	bytes = stat.Size()
	fmt.Println(bytes / MEGA)
	if bytes > 256.0*MEGA {
		blob.CreateBlockBlob(nil)
		fmt.Println("File is bigger than allowed for single block")
		buffer := make([]byte, 100*MEGA)
		var blocks []storage.Block
		for {
			length, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					return err
				}
				// fmt.Println("Should break")
				break
			}
			// fmt.Println(len(buffer), length)
			blockID := base64.StdEncoding.EncodeToString([]byte(shortuuid.New()))
			err = blob.PutBlock(blockID, buffer[0:length], nil)
			if err != nil {
				// fmt.Println("create block failed")
				return err
			}
			blocks = append(blocks, storage.Block{blockID, storage.BlockStatusUncommitted})
		}
		// uncommitted, err := blob.GetBlockList(storage.BlockListTypeUncommitted, nil)
		// fmt.Println(uncommitted)
		// fmt.Println(blocks)
		err = blob.PutBlockList(blocks, nil)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("File can be uploaded as single block")
		err := blob.CreateBlockBlobFromReader(file, nil)
		if err != nil {
			return err
		}
	}

	fmt.Println("Done")
	return nil
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

package main

import (
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
	"os"
)

func main() {

	var args = os.Args[1:]
	var argc = len(args)
	if argc < 2 {
		fmt.Printf("Incorrect number of arguments")
		return
	}

	var bucketName = args[0]
	var fileName = args[1]

	auth, err := aws.SharedAuth()
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(auth, aws.EUWest)
	resp, err := client.ListBuckets()

	if err != nil {
		log.Fatal(err)
	}

	var bucket s3.Bucket

	for key := range resp.Buckets {
		if resp.Buckets[key].Name == bucketName {
			bucket = resp.Buckets[key]
		}
	}

	if bucket.S3 == nil {
		log.Fatal("Can't find bucket")
	}

	data, err := bucket.Get(fileName)

	if err != nil {
		log.Fatal(err)
	}

	fo, err := os.Create(fileName)
	defer fo.Close()
  
	if err != nil {
		log.Fatal(err)
	}

	a, err := fo.Write(data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", a)
}

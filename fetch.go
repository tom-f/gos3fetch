package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"strconv"
	"net/http"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/cheggaaa/pb"
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

    resp2, _ := bucket.GetResponse(fileName)
    defer resp2.Body.Close()

    length := getLength(resp2)

    progressBar := pb.New(length).SetUnits(pb.U_BYTES)
    progressBar.Start()

	fo, err := os.Create(fileName)
	defer fo.Close()
  
	if err != nil {
		log.Fatal(err)
	}

	writer := io.MultiWriter(fo, progressBar)

	io.Copy(writer, resp2.Body)
	progressBar.Finish()

}

func getLength(resp *http.Response) int {
	for key, value := range resp.Header {
        if key == "Content-Length" {
            i, err := strconv.Atoi(value[0])
            if err == nil {
            	return i
            } 
        }
    }

    return 0
}


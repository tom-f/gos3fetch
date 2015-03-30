package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func main() {

	if len(os.Args[1:]) < 2 {
		fmt.Println("Usage is: gos3fetch [bucket] [file]")
		return
	}
	bucketName, fileName := os.Args[1], os.Args[2]

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
	i, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		log.Fatal(err)
	}
	return i
}

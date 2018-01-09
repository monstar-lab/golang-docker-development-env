// s3.go
//Check out mino guide to know more
//https://docs.minio.io/docs/golang-client-quickstart-guide

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/minio/minio-go"
)

type s3Handler struct{}

func (h s3Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "S3 Page\n\n")
	
	endpoint := "s3:9000"
	accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	s3 := newClient(endpoint, accessKeyID, secretAccessKey, false)

	// for requests comes to /s3/triggeraput, puts the s3_upload_test_file.txt as a new abject, adding a date suffix to the file name
	if strings.TrimPrefix(r.URL.Path, "/s3/") == "triggeraput" {
		s3.makeBucket("testbucket", "us-east-1")
		jp, _ := time.LoadLocation("Asia/Tokyo")
		s3.putFile("testbucket", "s3_upload_test_file_"+time.Now().In(jp).Format("2006-01-02 15:04:05")+".txt", "s3_upload_test_file.txt", "text/plain")
		fmt.Fprintf(w, "New object added. \n\n")
	} else {
		fmt.Fprintf(w, "Access /s3/triggeraput to add s3_upload_test_file.txt as a new object. \n\n")
	}

	fmt.Fprintf(w, "Objects2 in the testbucket...\nBrowse buckets on http://localhost:9000\n\n")

	// list current objects in the testbucket
	files := s3.listobjects("testbucket")
	for c, obkey := range files {
		fmt.Fprintln(w, c+1, ". ", obkey)
	}

}

type client struct {
	host string
	s3   *minio.Client
}

func newClient(host, key, secret string, insecure bool) *client {
	if host == "" {
		host = "s3.amazonaws.com"
		insecure = false
	}

	s3Client, err := minio.New(host, key, secret, insecure)
	if err != nil {
		log.Fatalln("minio.New", err)
	}

	return &client{
		host,
		s3Client,
	}
}

// cehck whether a bucket exists
func (c *client) bucketExists(bucket string) bool {
	_, err := c.s3.BucketExists(bucket)
	return err == nil
}

// make a new bucket if not exists
func (c *client) makeBucket(bucketName string, location string) {

	err := c.s3.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already have this bucket
		if c.bucketExists(bucketName) {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	log.Printf("Successfully created %s\n", bucketName)
}

// put a file on object stirage
func (c *client) putFile(bucketName string, objectName string, filePath string, contentType string) {

	// Upload the file with FPutObject
	n, err := c.s3.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

// get the list of objects in a bucket
func (c *client) listobjects(bucketName string) []string {

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	var files []string
	// List all objects from a bucket-name with a matching prefix.
	for object := range c.s3.ListObjects(bucketName, "/", true, doneCh) {
		if object.Err != nil {
			log.Println(object.Err)
		}
		files = append(files, object.Key)
	}
	return files
}

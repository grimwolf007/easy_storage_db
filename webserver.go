package main

import (
	"context"
	"fmt"
	"log"

	// File Uploads
	"os"

	// For configs as json
	"encoding/json"
	"io/ioutil"

	// Webserver Framework
	"net/http"

	"github.com/gin-gonic/gin"

	// MinioSDK
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//global

var minio_client = init_minio_client()

func main() {

	// Minio Settings (from minio/identity/account)
	// put_minio_env(minio_env{Ssl: false})
	// Golang Webserver
	start_webserver()

}

// Does not take params, Just runs a basic web-server
func start_webserver() {
	// Basic Webserver
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB Memory limit for multipart forms. Don't understand this

	//GET pages
	router.GET("/health-check", Healthcheck_page)
	router.GET("/bucket-list", BucketList_page)
	router.GET("/stop-server", StopServer_page)

	//POST pages
	router.POST("/upload", Upload)

	//start server
	router.Run() // Listen/serve 0.0.0.0:8080
}

func init_minio_client() *minio.Client {
	env := pull_minio_env()
	log.Println(env)
	log.Println("Setting up Minio Client")
	minioClient, err := minio.New(env.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.AccessID, env.SecretKey, ""),
		Secure: env.Ssl,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Minio Client Ready")
	return minioClient
}

// GET HealthCheck
func Healthcheck_page(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Health": "OK"})
}

// GET BucketList
func BucketList_page(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Bucket List": list_minio_buckets(minio_client)})
}

// GET Stop webserver page
func StopServer_page(c *gin.Context) {
	os.Exit(200)
}

// GET List Buckets
func list_minio_buckets(minio_client *minio.Client) string {
	client := minio_client
	list := ""
	bucket_list, _ := client.ListBuckets(context.Background())
	for message := range bucket_list {
		list += string(message) + "\n"
	}
	return list
}

// POST Upload
func Upload(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	dst := "./test/"

	//Print each file
	for _, file := range files {
		log.Println(file.Filename)
		log.Println(dst + file.Filename)

		//upload the file to destination
		c.SaveUploadedFile(file, dst+file.Filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

// uploads files from webserver to minio
func upload_to_minio() string {
	log.Print("Not yet setup")
	return "Not Setup Yet"
}

// minio_env Object
type minio_env struct {
	Endpoint, AccessID, SecretKey string
	Ssl                           bool
}

// Pulls the minio environment from .minio_secrets
func pull_minio_env() minio_env {
	file, _ := ioutil.ReadFile(".minio_secrets")
	env := minio_env{}
	_ = json.Unmarshal([]byte(file), &env)
	return env
}

// Saves the minio environment to .minio_secrets
// REQUIRES SSL USAGE!
// To-do: see if there is a new change before putting
func put_minio_env(env minio_env) {
	old_env := pull_minio_env()
	if env.Endpoint == "" {
		env.Endpoint = old_env.Endpoint
	}
	if env.AccessID == "" {
		env.AccessID = old_env.AccessID
	}
	if env.SecretKey == "" {
		env.SecretKey = old_env.SecretKey
	}
	env_json, _ := json.MarshalIndent(env, "", " ")
	err := ioutil.WriteFile(".minio_secrets", env_json, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

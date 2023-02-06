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

func debug(router gin.RouterGroup) {
	router.GET("/stop-server", StopServer_page)
}

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
	router.LoadHTMLGlob("template/*")

	router.MaxMultipartMemory = 8 << 20 // 8 MiB Memory limit for multipart forms. Don't understand this

	//GET pages
	router.GET("/", Index_page)
	router.GET("/health-check", Healthcheck_page)
	router.GET("/list_buckets", BucketList_page)
	router.GET("/bucket/:name", Get_bucket_page)
	router.GET("/upload/server", Upload_server_page)

	//POST pages
	router.POST("/upload/server", Upload_host)  // Must remove or create authenitacion for
	router.POST("/upload/object", Upload_minio) // Create login page for minio?
	router.POST("/create_bucket/:name", Create_bucket)
	router.POST("/remove_bucket/:name", Remove_bucket)

	//Debug pages
	if os.Getenv("DEBUG") == "TRUE" {
		router.GET("/stop-server", StopServer_page)
	}
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

// GET Index
func Index_page(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// GET HealthCheck
func Healthcheck_page(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Health": "OK"})
}

// ListBuckets
func list_minio_buckets(minio_client *minio.Client) []minio.BucketInfo {
	client := minio_client
	bucket_list, err := client.ListBuckets(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	return bucket_list
}

// GET BucketList
func BucketList_page(c *gin.Context) {
	list := list_minio_buckets(minio_client)
	c.JSON(http.StatusOK, list)
}

// GET Upload_server_page
func Upload_server_page(c *gin.Context) {
	c.JSON(http.StatusOK, "hello")
}

// GET Stop webserver page
// Stops the webserver with code 200 [DEBUG]
func StopServer_page(c *gin.Context) {
	os.Exit(200)
}

// GET Get_bucket_page
// Displays the name given in the uri as a json
func Get_bucket_page(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Bucket List": c.Param("name")})
}

// POST Upload_host (File to host)
func Upload_host(c *gin.Context) {
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		log.Fatalln(err)
	}
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

// POST Create_bucket
func Create_bucket(c *gin.Context) {
	location := "Test"
	name := c.Param("name")
	client := minio_client
	ctx := context.Background()
	err := client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, name)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", name)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", name)
	}
}

// POST Remove_bucket
func Remove_bucket(c *gin.Context) {
	client := minio_client
	bucket_name := c.Param("name")
	err := client.RemoveBucket(context.Background(), bucket_name)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

// POST Upload_to_minio uploads files from webserver to minio
func Upload_minio(c *gin.Context) {
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
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded to minio!", len(files)))
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

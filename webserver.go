package main

import (
	"fmt"
	"log"

	// For configs as json
	"encoding/json"
	"io/ioutil"

	// Webserver Framework
	"net/http"

	"github.com/gin-gonic/gin"
	// MinioSDK
)

func main() {

	// Minio Settings (from minio/identity/account)
	//minio_endpoint := "minio"
	//accessKeyID := ""
	//secretAccessKey := ""
	//useSSL := false
	// Put minio_env
	//put_minio_env(minio_env_params{Endpoint: minio_endpoint, AccessID: accessKeyID, SecretKey: secretAccessKey, Ssl: useSSL})
	put_minio_env(minio_env_params{Ssl: false})
	// Pull minio_env
	pull_minio_env()
	//Golang Webserver
	start_webserver()

	//Test upload
	upload_to_minio()
}

// Does not take params, Just runs a basic web-server
func start_webserver() {
	// Basic Webserver
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // Listen/serve 0.0.0.0:8080
}

// uploads files from webserver to minio
func upload_to_minio() {
	fmt.Println("Not yet setup")
}

// Pulls the minio environment from .minio_secrets
func pull_minio_env() minio_env_params {
	file, _ := ioutil.ReadFile(".minio_secrets")
	env := minio_env_params{}
	_ = json.Unmarshal([]byte(file), &env)
	return env
}

type minio_env_params struct {
	Endpoint, AccessID, SecretKey string
	Ssl                           bool
}

// Saves the minio environment to .minio_secrets
// REQUIRES SSL USAGE!
// To-do: see if there is a new change before putting
func put_minio_env(env minio_env_params) {
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

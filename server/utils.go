package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc/credentials"
)

type Secret struct {
	DbUser     string `json:"DB_USER"`
	DbPassword string `json:"DB_PASSWORD"`
}

type ServiceConfig struct {
	GrpcServerAddress  string
	GrpcServerCertPath string
	GrpcServerKeyPath  string
	DbName             string
	DbHost             string
	DbUser             string
	DbPassword         string
}

func strToPair(envVar string) (string, string) {
	keyValue := strings.Split(envVar, "=")
	return keyValue[0], keyValue[1]
}

func getEnvironment() map[string]string {
	values := os.Environ()
	environment := make(map[string]string)

	for _, e := range values {
		k,v := envVarToPair(e)
		environment[k] = v
	}

	return environment
}

func getenv(env map[string]string, key, fallback string) string {
	value := env[key]
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getSecrets(raw []byte) Secret {
	secret := Secret{}
	if err := json.Unmarshal(raw, &secret); err != nil {
		log.Fatalf("[*] Unable to decode secrets json: %v", err)
	}
	return secret
}

func getConfigIO() ServiceConfig {
	env := getEnvironment()
	path := getVar(env, "SECRET_PATH", "config/app_secrets.json")

	raw, err := ioutil.ReadFile(path)
	if err != nil || len(raw) == 0 {
		log.Fatalf("Unable to read secrets file: %v", err)
	}

	return getConfig(env, raw)
}

func getConfig(env map[string]string, secret []byte) ServiceConfig {
	config := ServiceConfig{}
	secrets := getSecrets(secret)

	config.GrpcServerAddress = getVar(env, "GRPC_SERVER_ADDRESS", "localhost:9000")
	config.GrpcServerCertPath = getVar(env, "GRPC_SERVER_CERT_PATH", "../cert/localhost.crt")
	config.GrpcServerKeyPath = getVar(env, "GRPC_SERVER_KEY_PATH", "../cert/localhost.key")
	config.DbName = getVar(env, "DB_NAME", "books")
	config.DbHost = getVar(env, "DB_HOST", "localhost")
	config.DbUser = secrets.DbUser
	config.DbPassword = secrets.DbPassword

	return config
}

func getCredentials(config ServiceConfig) credentials.TransportCredentials {
	creds, err := credentials.NewServerTLSFromFile(
		config.GrpcServerCertPath, config.GrpcServerKeyPath)
	if err != nil {
		log.Fatalf("[*] Failed to setup TLS: %v", err)
	}

	return creds
}

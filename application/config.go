package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Configuration
type Configurations struct {
	Server    ServerConfig
	ChunkSize int64
}

type ServerConfig struct {
	Host string
	Port int
}

func (sc *ServerConfig) ConnectionString() string {
	return strings.Join([]string{sc.Host, ":", fmt.Sprint(sc.Port)}, "")
}

func SetupConfiguration() *Configurations {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	var configuration Configurations

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("WARNING: Error reading config file, %s", err)

		return getConfigurationFromEnv()
	}

	err := viper.Unmarshal(&configuration)

	if err != nil {
		log.Fatalf("FATAL: Unable to decode into struct, %v", err)
	}

	return &configuration
}

func getConfigurationFromEnv() *Configurations {
	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		log.Panic("PORT information not available")
	}

	chunkSize, err := strconv.Atoi(os.Getenv("CHUNK_SIZE"))
	if err != nil {
		chunkSize = 512000
	}

	return &Configurations{
		Server: ServerConfig{
			Host: os.Getenv("HOST"),
			Port: port,
		},
		ChunkSize: int64(chunkSize),
	}
}

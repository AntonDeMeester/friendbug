package main

import (
	"github.com/joho/godotenv"

	// Own packages
	client "friendbug/internal"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	client.ExampleRedis()
	
}

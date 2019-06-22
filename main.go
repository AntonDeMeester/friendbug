package main

import (
	"fmt"
	"github.com/joho/godotenv"

	"friendbug/cmd/friendbug"

)

func main() {
	err := godotenv.Load(".env")
	if (err != nil) {
		panic(err)
	}
	fmt.Println("Starting service")
	friendbug.ContactFriends()
	fmt.Println("Done with service")
}
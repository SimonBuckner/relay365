package main

import (
	"fmt"
	"log"

	"github.com/SimonBuckner/relay365"
	"github.com/SimonBuckner/relay365/graphhelper"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Go Graph Tutorial")
	fmt.Println()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env" + err.Error())
	}

	graphHelper := graphhelper.NewGraphHelper()

	relay365.InitializeGraph(graphHelper)

	relay365.GreetUser(graphHelper)

	var choice int64 = -1

	for {
		fmt.Println("Please choose one of the following options:")
		fmt.Println("0. Exit")
		fmt.Println("1. Display access token")
		fmt.Println("2. List my inbox")
		fmt.Println("3. Send mail")
		fmt.Println("4. Make a Graph call")

		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			choice = -1
		}

		switch choice {
		case 0:
			// Exit the program
			fmt.Println("Goodbye...")
		case 1:
			// Display access token
			relay365.DisplayAccessToken(graphHelper)
		case 2:
			// List emails from user's inbox
			relay365.ListInbox(graphHelper)
		case 3:
			// Send an email message
			relay365.SendMail(graphHelper)
		case 4:
			// Run any Graph code
			relay365.MakeGraphCall(graphHelper)
		default:
			fmt.Println("Invalid choice! Please try again.")
		}

		if choice == 0 {
			break
		}
	}
}

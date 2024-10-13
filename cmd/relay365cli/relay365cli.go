package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/simonbuckner/relay365/graphhelper"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Go Graph App-Only Tutorial")
	fmt.Println()

	// Load .env files
	// .env.local takes precedence (if present)
	godotenv.Load(".env.local")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	graphHelper := graphhelper.NewGraphHelper()

	initializeGraph(graphHelper)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	var choice rune
	// var key keyboard.Key

	for {
		fmt.Println("Please choose one of the following options:")
		fmt.Println("0. Exit")
		fmt.Println("1. Display access token")
		fmt.Println("2. List users")
		fmt.Println("3. Send email")
		fmt.Println("4. List inbox")

		choice, _, err = keyboard.GetSingleKey()
		if err != nil {
			choice = rune('!')
		}

		switch choice {
		case rune('0'):
			// Exit the program
			fmt.Println("Goodbye...")
			os.Exit(0)
		case rune('1'):
			// Display access token
			displayAccessToken(graphHelper)
		case rune('2'):
			// List users
			listUsers(graphHelper)
		case rune('3'):
			// Send Email
			sendMail(graphHelper)
		case rune('4'):
			// Send Email
			listInbox(graphHelper)
		default:
			fmt.Println("Invalid choice! Please try again.")
		}
	}
}

func initializeGraph(graphHelper *graphhelper.GraphHelper) {
	err := graphHelper.InitializeGraphForAppAuth()
	if err != nil {
		log.Panicf("Error initializing Graph for app auth: %v\n", err)
	}
}

func displayAccessToken(graphHelper *graphhelper.GraphHelper) {
	token, err := graphHelper.GetAppToken()
	if err != nil {
		log.Panicf("Error getting user token: %v\n", err)
	}

	fmt.Printf("App-only token: %s", *token)
	fmt.Println()
}

func listUsers(graphHelper *graphhelper.GraphHelper) {
	var nextUrl *string = nil
	var choice rune

	for {
		fmt.Printf("NextURL: %v\n", nextUrl)
		users, err := graphHelper.GetUsers(nextUrl)
		if err != nil {
			log.Panicf("Error getting users: %v", err)
		}

		// Output each user's details
		for _, user := range users.GetValue() {
			fmt.Printf("User: %s\n", *user.GetDisplayName())
			fmt.Printf("  ID: %s\n", *user.GetId())

			noEmail := "NO EMAIL"
			email := user.GetMail()
			if email == nil {
				email = &noEmail
			}
			fmt.Printf("  Email: %s\n", *email)
		}

		// If GetOdataNextLink does not return nil,
		// there are more users available on the server
		nextUrl = users.GetOdataNextLink()

		fmt.Println()
		fmt.Printf("More users available? %t\n", nextUrl != nil)

		if nextUrl == nil {
			break
		}

		fmt.Printf("")
		fmt.Printf("NextURL: %v\n", nextUrl)
		fmt.Printf("Display the next page? Y/N: ")

		choice, _, err = keyboard.GetSingleKey()
		if err != nil {
			continue
		}

		if choice != rune('Y') && choice != rune('y') {
			return
		}
		fmt.Printf("")
	}
}

func sendMail(graphHelper *graphhelper.GraphHelper) {

	from := "simon.buckner@gmal.co.uk"
	subject := "Testing Microsoft Graph"
	body := "Hello world!"
	to := "simonbuckner@hotmail.com"

	err := graphHelper.SendMail(&from, &subject, &body, &to)
	if err != nil {
		log.Panicf("Error sending mail: %v", err)
	}

	fmt.Println("Mail sent.")
	fmt.Println()
}

func listInbox(graphHelper *graphhelper.GraphHelper) {
	messages, err := graphHelper.GetInbox("simon.buckner@gmal.co.uk")
	if err != nil {
		log.Panicf("Error getting user's inbox: %v", err)
	}

	// Load local time zone
	// Dates returned by Graph are in UTC, use this
	// to convert to local
	location, err := time.LoadLocation("Local")
	if err != nil {
		log.Panicf("Error getting local timezone: %v", err)
	}

	// Output each message's details
	for _, message := range messages.GetValue() {
		fmt.Printf("Message: %s\n", *message.GetSubject())
		fmt.Printf("  From: %s\n", *message.GetFrom().GetEmailAddress().GetName())

		status := "Unknown"
		if *message.GetIsRead() {
			status = "Read"
		} else {
			status = "Unread"
		}
		fmt.Printf("  Status: %s\n", status)
		fmt.Printf("  Received: %s\n", (*message.GetReceivedDateTime()).In(location))
	}

	// If GetOdataNextLink does not return nil,
	// there are more messages available on the server
	nextLink := messages.GetOdataNextLink()

	fmt.Println()
	fmt.Printf("More messages available? %t\n", nextLink != nil)
	fmt.Println()
}

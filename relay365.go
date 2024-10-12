package relay365

import (
	"fmt"
	"log"
	"time"

	"github.com/SimonBuckner/relay365/graphhelper"
)

func InitializeGraph(graphHelper *graphhelper.GraphHelper) {
	err := graphHelper.InitializeGraphForUserAuth()
	if err != nil {
		log.Panicf("Error initializing Graph for user auth: %v\n", err)
	}
}

func GreetUser(graphHelper *graphhelper.GraphHelper) {
	user, err := graphHelper.GetUser()
	if err != nil {
		log.Panicf("Error getting user: %v\n", err)
	}

	fmt.Printf("Hello, %s!\n", *user.GetDisplayName())

	// For Work/school accounts, email is in Mail property
	// Personal accounts, email is in UserPrincipalName
	email := user.GetMail()
	if email == nil {
		email = user.GetUserPrincipalName()
	}

	fmt.Printf("Email: %s\n", *email)
	fmt.Println()
}

func DisplayAccessToken(graphHelper *graphhelper.GraphHelper) {
	token, err := graphHelper.GetUserToken()
	if err != nil {
		log.Panicf("Error getting user token: %v\n", err)
	}

	fmt.Printf("User token: %s", *token)
	fmt.Println()
}

func ListInbox(graphHelper *graphhelper.GraphHelper) {
	messages, err := graphHelper.GetInbox()
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
		fmt.Printf(" Email    : %s\n", *message.GetFrom().GetEmailAddress().GetAddress())
		fmt.Printf(" From     : %s\n", *message.GetFrom().GetEmailAddress().GetName())

		status := "Unknown"
		if *message.GetIsRead() {
			status = "Read"
		} else {
			status = "Unread"
		}
		fmt.Printf(" Status   : %s\n", status)
		fmt.Printf(" Received : %s\n", (*message.GetReceivedDateTime()).In(location))
	}

	// If GetOdataNextLink does not return nil,
	// there are more messages available on the server
	nextLink := messages.GetOdataNextLink()

	fmt.Println()
	fmt.Printf("More messages available? %t\n", nextLink != nil)
	fmt.Println()
}

func SendMail(graphHelper *graphhelper.GraphHelper) {
	// Send mail to the signed-in user
	// Get the user for their email address
	user, err := graphHelper.GetUser()
	if err != nil {
		log.Panicf("Error getting user: %v", err)
	}

	// For Work/school accounts, email is in Mail property
	// Personal accounts, email is in UserPrincipalName
	email := user.GetMail()
	if email == nil {
		email = user.GetUserPrincipalName()
	}

	subject := "Testing Microsoft Graph"
	body := "Hello world!"
	err = graphHelper.SendMail(&subject, &body, email)
	if err != nil {
		log.Panicf("Error sending mail: %v", err)
	}

	fmt.Println("Mail sent.")
	fmt.Println()
}

func MakeGraphCall(graphHelper *graphhelper.GraphHelper) {
	// TODO
}

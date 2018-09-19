// Smartling SDK v2 Auth Test Example.
//
// Example shows usage of Smartling authentication API
// https://help.smartling.com/v1.0/reference#authentication-1
//
// This example does nothing except the authentication call.
// Useful for testing your user identifier / token.
//
// `UserID` and `TokenSecret` should be specified in the
// example_credentails_test.go before running that test.

package smartling_test

import (
	"fmt"
	"log"

	smartling "github.com/Smartling/api-sdk-go"
)

func ExampleAuth() {
	log.Printf("Initializing smartling client and performing autorization")

	client := smartling.NewClient(UserID, TokenSecret)

	err := client.Authenticate()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Authentication Successfull")

	// Output: Authentication Successfull
}

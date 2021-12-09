package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/davecgh/go-spew/spew"
	a "github.com/microsoft/kiota/authentication/go/azure"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/microsoft/graph"
	u "github.com/microsoftgraph/msgraph-sdk-go/users"
	ui "github.com/microsoftgraph/msgraph-sdk-go/users/item"
)

func main() {

	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("B2C_TENANT_ID"),
		os.Getenv("B2C_CLIENT_ID"),
		os.Getenv("B2C_CLIENT_SECRET"),
		&azidentity.ClientSecretCredentialOptions{
			ClientOptions: policy.ClientOptions{
				Retry: policy.RetryOptions{
					MaxRetries:    3,
					MaxRetryDelay: time.Duration(30) * time.Second,
				},
			},
		})

	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		fmt.Printf("Error authentication provider: %v\n", err)
		return
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		fmt.Printf("Error creating adapter: %v\n", err)
		return
	}
	client := msgraphsdk.NewGraphServiceClient(adapter)

	// Change your identity.
	//getUser(client, "5fb24546-73b4-437e-acb9-1a5f0975ba8e")
	getUser(client, "e2a12752-ca9e-44ae-bbea-6384c25f58e6")
	listUser(client)
}

func getUser(client *msgraphsdk.GraphServiceClient, id string) {

	options := &ui.UserRequestBuilderGetOptions{
		Q: &ui.UserRequestBuilderGetQueryParameters{
			Expand: []string{"extensions($filter=id eq 'com.example.moris')"},
			Select: []string{"id", "createdDateTime", "accountEnabled", "displayName", "userPrincipalName"},
		},
	}

	user, err := client.UsersById(id).Get(options)
	if err != nil {
		fmt.Printf("Error geuser : %v\n", err)
		return
	}

	printUser(*user)
}

func listUser(client *msgraphsdk.GraphServiceClient) {

	f := "startsWith(mail,'bob')"
	options := &u.UsersRequestBuilderGetOptions{
		Q: &u.UsersRequestBuilderGetQueryParameters{
			Filter: &f,
			Select: []string{"createdDateTime"},
		},
	}

	result, err := client.Users().Get(options)
	if err != nil {
		fmt.Printf("Error list user: %v\n", err)
		return
	}

	for _, user := range result.GetValue() {
		printUser(user)
	}
}

func printUser(user graph.User) {

	spew.Dump(user)

	fmt.Printf("------------------------------------------------------------------------------------\n")
	fmt.Printf("user.GetId(): %v\n", p(user.GetId()))
	fmt.Printf("user.GetUserPrincipalName(): %v\n", p(user.GetUserPrincipalName()))

	fmt.Printf("user.GetDisplayName(): %v\n", p(user.GetDisplayName()))
	fmt.Printf("user.GetSurname(): %v\n", p(user.GetSurname()))
	fmt.Printf("user.GetGivenName(): %v\n", p(user.GetGivenName()))

	fmt.Printf("user.GetMail(): %v\n", p(user.GetMail()))
	fmt.Printf("user.GetCompanyName(): %v\n", p(user.GetCompanyName()))
	fmt.Printf("user.GetCreatedDateTime(): %v\n", user.GetCreatedDateTime())
	fmt.Printf("user.GetAdditionalData(): %v\n", user.GetAdditionalData())

	for _, e := range user.GetExtensions() {
		for k, v := range e.GetAdditionalData() {
			fmt.Printf("k: %v\n", k)
			fmt.Printf("v: %v\n", r(v))
		}
	}
}

func p(p *string) string {
	if p == nil {
		return "<nil>"
	} else {
		return *p
	}
}

func r(value interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(value)).Interface()
}

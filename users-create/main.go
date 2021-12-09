package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	a "github.com/microsoft/kiota/authentication/go/azure"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go/models/microsoft/graph"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/microsoftgraph/msgraph-sdk-go/users/item/extensions"
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

	id, err := createUser(client)

	if err != nil {
		fmt.Printf("Error creating user : %v\n", err)
	}

	err = setExtension(client, id)

	if err != nil {
		fmt.Printf("Error setting extension : %v\n", err)
	}
}

func createUser(client *msgraphsdkgo.GraphServiceClient) (string, error) {

	f := false
	t := true
	tenant := os.Getenv("B2C_TENANT_NAME")
	uid := uuid.New().String()
	upn := uid + "@" + tenant + ".onmicrosoft.com"
	user := graph.NewUser()
	user.SetDisplayName(ref("Shohei Ohtani"))
	user.SetUserPrincipalName(&upn)
	user.SetMailNickname(ref("Shohei"))
	user.SetAccountEnabled(&t)

	prof := graph.NewPasswordProfile()
	prof.SetForceChangePasswordNextSignIn(&f)
	prof.SetPassword(ref(uuid.New().String()))
	user.SetPasswordProfile(prof)
	user.SetPasswordPolicies(ref("DisablePasswordExpiration"))

	d := &users.UsersRequestBuilderPostOptions{
		Body: user,
	}
	result, err := client.Users().Post(d)

	if err != nil {
		fmt.Printf("Error creating new user %v\n", err)
		return "", err
	}

	fmt.Printf("result.GetId(): %v\n", *result.GetId())
	fmt.Printf("result.GetDisplayName(): %v\n", *result.GetDisplayName())
	fmt.Printf("result.GetMailNickname(): %v\n", result.GetMailNickname())

	spew.Dump(result)
	return *result.GetId(), nil
}

func setExtension(client *msgraphsdkgo.GraphServiceClient, id string) error {

	// The following properties cannot be set in the initial POST request
	ext := graph.NewExtension()
	attr := make(map[string]interface{})
	attr["extensionName"] = ref("com.example.moris")
	attr["value"] = ref("you can set some values.")
	attr["@odata.type"] = ref("#microsoft.graph.openTypeExtension")
	ext.SetAdditionalData(attr)

	options := &extensions.ExtensionsRequestBuilderPostOptions{
		Body: ext,
	}

	result, err := client.UsersById(id).Extensions().Post(options)

	if err != nil {
		return err
	}
	spew.Dump(result)
	return nil
}

func ref(str string) *string {
	return &str
}

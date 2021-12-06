package main

import (
	"context"
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	a "github.com/microsoft/kiota/authentication/go/azure"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

func main() {
	fmt.Printf("main\n")

	cred, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID: os.Getenv("B2C_TENANT_ID"),
		ClientID: os.Getenv("B2C_CLIENT_ID"),
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
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

	result, err := client.Me().Get(nil)
	if err != nil {
		fmt.Printf("Error getting the drive: %v\n", err)
	}
	fmt.Printf("Me : %v %v\n",
		result.GetGivenName(),
		result.GetDisplayName())

	fmt.Printf("Me : %s %s\n",
		*result.GetGivenName(),
		*result.GetDisplayName())
}

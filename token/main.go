package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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
				Logging: policy.LogOptions{
					IncludeBody: true,
				},
			},
		})

	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
	}

	// Get Token Directly
	token, _ := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})

	fmt.Printf("https://jwt.ms/#access_token=%s\n", token.Token)
	fmt.Printf("You can decode this token on https://jwt.ms/")
}

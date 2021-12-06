# msgraph-go-samples

MS Graph SDK for golan のサンプル。

[microsoftgraph/msgraph-sdk-go: Microsoft Graph SDK for Go](https://github.com/microsoftgraph/msgraph-sdk-go)

Graph API には、`v1.0` と `beta` があり、`beta`をアクセスするSDKは別途公開されているので、それを利用すること。ただしプロダクション利用はできない。

[microsoftgraph/msgraph-beta-sdk-go: Microsoft Graph Beta Go SDK](https://github.com/microsoftgraph/msgraph-beta-sdk-go)

現時点では、いずれのSDKもPreviewである。

## 認証

Graph APIではないが補足しておく。

### デバイスコード認証

以下は、デバイスコードで認証する方法。

```golang
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
```

よくあるデバイスコード入力メッセージが出力されるので、ブラウザでコードを入力すると認証される。

```
To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code XXXXXXXXX to authenticate.
```

### クライアントクレデンシャル認証

Azure AD B2C にアプリケーションを登録し、テナントID、クライアントID、クライアントシークレットを取得、しかるべきパーミッションを与えておく。

```golang
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
```

## サンプル

TODO

|フォルダ   | 内容   |
|---|---|
| me  | 自身の情報を表示  |
| token |  アクセストークンを直接取得する
| users | user の一覧、取得、作成など |

## Known Issue

- Select クエリが正しく動いていない

以上

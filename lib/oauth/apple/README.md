# Apple 登录验证

### 根据code验证



```go

var (
    code string
    uid string 
)


//验证
secret, _ := apple.GenerateClientSecret(secret, teamID, clientID, keyID)

//
client := apple.NewClient()

vReq := apple.AppValidationTokenRequest{
	ClientID:     clientID,
	ClientSecret: secret,
	Code:         code, // code值
}

var resp apple.ValidationResponse

// 验证
client.VerifyAppToken(context.Background(), vReq, &resp)

// 获取用户信息

unique, _ := apple.GetUniqueID(resp.IDToken)

// 验证是否同一个用户
if unique!=uid{
    return
}


fmt.Println(unique)

```
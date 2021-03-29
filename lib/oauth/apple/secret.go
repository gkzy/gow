package apple

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

/*
GenerateClientSecret generates the client secret used to make requests to the validation server.
The secret expires after 6 months
secret - Private key from Apple obtained by going to the keys section of the developer section
teamID - Your 10-character Team ID
clientID - Your Services ID, e.g. com.aaronparecki.services
keyID - Find the 10-char Key ID value from the portal
*/
func GenerateClientSecret(secret, teamID, clientID, keyID string) (string, error) {
	block, _ := pem.Decode([]byte(secret))
	if block == nil {
		return "", errors.New("empty block after decoding")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Create the Claims
	now := time.Now()
	claims := &jwt.StandardClaims{
		Issuer:    teamID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Hour*24*180 - time.Second).Unix(), // 180 days
		Audience:  "https://appleid.apple.com",
		Subject:   clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["key"] = keyID

	return token.SignedString(privKey)
}
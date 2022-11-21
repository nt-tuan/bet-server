package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type AuthService struct {
	issuer    string
	secretJwt string
}

func NewAuthServive(issuer string, secretJwt string) AuthService {
	return AuthService{
		issuer:    issuer,
		secretJwt: secretJwt,
	}
}

func (a *AuthService) ParseIdToken(token string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	keySet, err := jwk.Fetch(ctx, "https://login.microsoftonline.com/b035d9fd-a388-48f2-9abf-cea8457262ce/discovery/v2.0/keys")
	parsedtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}
		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}
		return publickey, nil
	}, jwt.WithoutClaimsValidation())

	fmt.Printf("%v parser error\n", err)

	if err != nil {
		fmt.Printf("parse error %v", err)
		return "", err
	}

	claims, ok := parsedtoken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("can-not-parse-token")
	}

	var now = time.Now().Unix() - 60*60
	if !claims.VerifyExpiresAt(now, true) {
		return "", errors.New("token-expired")
	}

	if !claims.VerifyIssuer(a.issuer, true) {
		return "", errors.New("invalid-issuer")
	}
	var email = fmt.Sprintf("%v", claims["preferred_username"])
	return email, nil
}

func (a *AuthService) CreateAccessToken(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix() //Token hết hạn sau 3 ngày
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secretJwt))
}

func (a *AuthService) ValidateAccessToken(token string) (string, error) {
	parsedtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.secretJwt), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := parsedtoken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("can-not-parse-token")
	}
	if claims.VerifyExpiresAt(time.Now().UnixMilli(), true) {
		return "", errors.New(("token-expired"))
	}

	var username = fmt.Sprintf("%v", claims["username"])
	return username, nil
}

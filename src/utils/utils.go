package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type HttpsHasuraIoJwtClaims struct {
	XHasuraAllowedRoles []string `json:"x-hasura-allowed-roles"`
	XHasuraDefaultRole  string   `json:"x-hasura-default-role"`
	XHasuraUserid       string   `json:"X-Hasura-User-Id"`
}

type LoginClaims struct {
	HttpsHasuraIoJwtClaims HttpsHasuraIoJwtClaims `json:"https://hasura.io/jwt/claims"`
	jwt.RegisteredClaims
}

var userLoginHttpsHasuraIoJwtClaims HttpsHasuraIoJwtClaims = HttpsHasuraIoJwtClaims{
	XHasuraAllowedRoles: []string{"user"},
	XHasuraDefaultRole:  "user",
}

var userLoginClaims LoginClaims = LoginClaims{
	userLoginHttpsHasuraIoJwtClaims,
	jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	},
}

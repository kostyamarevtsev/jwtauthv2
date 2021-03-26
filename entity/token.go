package entity

import (
	"jwtauthv2"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
)

const AccessTokenTTL = time.Minute * 30
const RefreshTokenTTL = time.Hour * 24

var SECRET = []byte(os.Getenv("SECRET"))

type CustomClaims struct {
	User ID `json:"user"`
	jwt.StandardClaims
}

type TokenPair struct {
	Access, Refresh string
}

func IssueToken(user *ID, ttl time.Duration) (string, error) {

	claims := CustomClaims{
		*user,
		jwt.StandardClaims{
			IssuedAt:  jwt.Now(),
			ExpiresAt: jwt.At(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(SECRET)

	if err != nil {
		return "", &jwtauthv2.Error{
			Op:  "IssueToken",
			Err: err,
		}
	}

	return ss, nil
}

func ParseToken(ss string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(ss, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	if err != nil {
		if msg, ok := err.(*jwt.TokenExpiredError); ok {
			return nil, &jwtauthv2.Error{
				Code:    jwtauthv2.EPERMISSION,
				Message: msg.Error(),
			}
		}

		return nil, &jwtauthv2.Error{
			Code:    jwtauthv2.EINVALID,
			Message: err.Error(),
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, &jwtauthv2.Error{
		Code:    jwtauthv2.EINVALID,
		Message: "Token error",
	}

}

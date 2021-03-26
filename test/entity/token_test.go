package entity_test

import (
	"jwtauthv2"
	"jwtauthv2/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseToken(t *testing.T) {
	user := entity.NewID()

	validToken, err := entity.IssueToken(&user, entity.AccessTokenTTL)
	require.Nil(t, err)

	expiredToken, err := entity.IssueToken(&user, 0)
	require.Nil(t, err)

	testTable := []struct {
		name               string
		token              string
		expectedErrCode    string
		expectedClaimValue entity.ID
	}{
		{"Valid Token", validToken, "", user},
		{"Expired Token", expiredToken, jwtauthv2.EPERMISSION, uuid.Nil},
		{"Invalid Token", "123", jwtauthv2.EINVALID, uuid.Nil},
	}

	for _, c := range testTable {
		t.Run(c.name, func(t *testing.T) {
			claims, err := entity.ParseToken(c.token)

			if err != nil {
				code := jwtauthv2.ErrorCode(err)
				require.Equal(t, code, c.expectedErrCode)
			}

			if claims != nil {
				require.Equal(t, claims.User, c.expectedClaimValue)
			}

		})
	}

}

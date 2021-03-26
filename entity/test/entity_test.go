package entity_test

import (
	"jwtauthv2"
	"jwtauthv2/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {

	type input struct {
		name, password string
	}

	testTable := []struct {
		name            string
		input           *input
		expectedValue   *entity.User
		expectedErrCode string
	}{
		{
			name:            "Short name",
			input:           &input{"I", "123456"},
			expectedValue:   &entity.User{Name: "", Password: ""},
			expectedErrCode: jwtauthv2.EINVALID,
		},

		{
			name:            "Short password",
			input:           &input{"Ivan", "1234"},
			expectedValue:   &entity.User{Name: "Ivan", Password: ""},
			expectedErrCode: jwtauthv2.EINVALID,
		},

		{
			name:            "Long name",
			input:           &input{"abckajfhutofdncdligfnwdudjdfnejfeifeenf", "1234"},
			expectedValue:   &entity.User{Name: "", Password: ""},
			expectedErrCode: jwtauthv2.EINVALID,
		},

		{
			name:            "Long password",
			input:           &input{"Ivan", "abckajfhutofdncdligfnwdudjdfnejfeifeenfasdasdasdsdfsdf"},
			expectedValue:   &entity.User{Name: "Ivan", Password: ""},
			expectedErrCode: jwtauthv2.EINVALID,
		},

		{
			name:            "Valid User",
			input:           &input{"Ivan", "abckaj"},
			expectedValue:   &entity.User{Name: "Ivan", Password: entity.GetHashString("abckaj")},
			expectedErrCode: "",
		},
	}

	for _, c := range testTable {
		t.Run(c.name, func(t *testing.T) {
			user, err := entity.NewUser(c.input.name, c.input.password)

			if user != nil && (user.Name != c.expectedValue.Name || user.Password != c.expectedValue.Password) {
				t.Errorf("Error in '%v' case (expectedValue)", c.name)
			}

			if code := jwtauthv2.ErrorCode(err); code != c.expectedErrCode {
				t.Errorf("Error in '%v' case (expectedErrCode)", c.name)
			}

		})
	}
}

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

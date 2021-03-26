package entity_test

import (
	"jwtauthv2"
	"jwtauthv2/entity"
	"testing"
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

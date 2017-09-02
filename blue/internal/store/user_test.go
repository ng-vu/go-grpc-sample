package store

import "testing"

func Test_User_HashPassword(t *testing.T) {
	secrets := []string{
		"hello",
		"123@asd#",
		"this-is-a-long-string",
	}

	t.Run("Test hash password", func(t *testing.T) {
		for _, secret := range secrets {
			s := []byte(secret)
			hashpwd := hexa(saltedHashPassword(s))
			if !isPasswordMatch(dehexa(hashpwd), s) {
				t.Error("Expect password match", s, hashpwd)
			}
			if isPasswordMatch(dehexa(hashpwd), []byte("invalid")) {
				t.Error("Expect password not match", "invalid", hashpwd)
			}
		}
	})
}

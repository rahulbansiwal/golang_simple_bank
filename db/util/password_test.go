package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedpassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedpassword1)

	err = CheckPassword(password, hashedpassword1)
	require.NoError(t, err)

	wrongPasword := RandomString(7)
	err = CheckPassword(wrongPasword, hashedpassword1)
	require.Error(t, err)

	hashedpassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedpassword2)
	require.NotEqual(t,hashedpassword1,hashedpassword2)
}

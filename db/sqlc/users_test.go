package sqlc

import (
	"context"
	"fmt"
	"simple_bank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := CreateRandomUser(t)
	getuser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, getuser)
	require.Equal(t, getuser.Username, user.Username)
	require.Equal(t, getuser.FullName, user.FullName)
	require.Equal(t, getuser.Email, user.Email)
	require.Equal(t, getuser.HashedPassword, user.HashedPassword)
	require.NotEmpty(t, user.CreatedAt)

}

func CreateRandomUser(t *testing.T) User {
	hashedPassword,err := util.HashPassword(util.RandomString(6))
	require.NoError(t,err)
	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), args)
	fmt.Println(user)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, args.Username)
	require.Equal(t, user.Email, args.Email)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.CreatedAt)
	return user
}

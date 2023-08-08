package token

import (
	"fmt"
	"simple_bank/db/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreatePasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expired_at := issuedAt.Add(duration)
	token,payload, err := maker.CreateToken(username, duration)
	fmt.Println(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t,payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expired_at, time.Second)
}

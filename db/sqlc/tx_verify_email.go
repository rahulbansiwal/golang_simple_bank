package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (s *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			fmt.Printf("from update verify email:%s", err)
			return err
		}
		result.User, err = q.UdpateUser(ctx, UdpateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		fmt.Printf("from update user:%s", err)
		return err

	})
	return result, err
}

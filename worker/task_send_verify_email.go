package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/db/util"

	"github.com/hibiken/asynq"
)

const TaskSendAndVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (r *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPaylaod, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error is:%w", err)
	}
	task := asynq.NewTask(TaskSendAndVerifyEmail, jsonPaylaod, opts...)
	_, err = r.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("error is:%w", err)
	}
	return nil
}

func (processer *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("error while unmarshaling payload:%w", asynq.SkipRetry)
	}

	user, err := processer.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("username doesnt exist:%w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user:%w", err)
	}
	arg := db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}
	verifyEmail, err := processer.store.CreateVerifyEmail(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create verify email:%w", err)
	}

	subject := "Welcome to Simple Bank"
	verifyURL := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	 Thank you for registring with us!<br/>
	 Please <a href="%s">click here</a> to verify you email address. <br/>
	 `, user.FullName, verifyURL)
	to := []string{user.Email}

	err = processer.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email:%s", err)
	}

	fmt.Print(user)
	return nil
}

package gapi

import (
	"context"
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/db/util"
	"simple_bank/pb"
	"simple_bank/val"
	"simple_bank/worker"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidateCreatUserRequest(req)
	if violations != nil {
		badrequest := &errdetails.BadRequest{
			FieldViolations: violations,
		}
		return nil, status.Errorf(codes.InvalidArgument, badrequest.String())
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to Hash password: %s", err)
	}
	args := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			Email:          req.GetEmail(),
			FullName:       req.GetFullName(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
			}
			return s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)

		},
	}

	txResult, err := s.store.CreateUserTx(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exsist for %s", req.GetUsername())
			}

		}
		return nil, status.Errorf(codes.Internal, "Failed to create user")
	}

	//TODO: use db transaction

	res := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}
	return res, nil
}

func ValidateCreatUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: err.Error(),
		})
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "email",
			Description: err.Error(),
		})
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "password",
			Description: err.Error(),
		})
	}
	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "full_name",
			Description: err.Error(),
		})
	}
	fmt.Println(violations)
	return violations
}

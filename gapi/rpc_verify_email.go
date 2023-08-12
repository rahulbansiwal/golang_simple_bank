package gapi

import (
	"context"
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/pb"
	"simple_bank/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := ValidateVerifyEmailRequest(req)
	if violations != nil {
		badrequest := &errdetails.BadRequest{
			FieldViolations: violations,
		}
		return nil, status.Errorf(codes.InvalidArgument, badrequest.String())
	}
	txResult,err:= s.store.VerifyEmailTx(ctx,db.VerifyEmailTxParams{
		EmailId: req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil{
		fmt.Print(err)
		return nil,status.Errorf(codes.Internal, "failed to verify email")
	}

	res := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	return res, nil
}

func ValidateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "email_id",
			Description: err.Error(),
		})
	}
	if err := val.ValidateSecretCode(req.SecretCode); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "secret_code",
			Description: err.Error(),
		})
	}
	return violations
}

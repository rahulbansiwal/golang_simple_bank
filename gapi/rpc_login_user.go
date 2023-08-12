package gapi

import (
	"context"
	"database/sql"
	db "simple_bank/db/sqlc"
	"simple_bank/db/util"
	"simple_bank/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "username not exsist")
		}
		return nil, status.Errorf(codes.Internal, "could not get user details for: %s", req.Username)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "password check failed")
	}

	accesstoken, accessPayload, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating token: %s", err)
	}
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(req.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating token: %s", err)
	}

	mtdt := s.extractMetaData(ctx)
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating session in the session table: %s", err)
	}

	res := &pb.LoginUserResponse{
		User:                  convertUser(user),
		Sessionid:             session.ID.String(),
		AccessToken:           accesstoken,
		AccessTokenExpiredAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return res, nil
}

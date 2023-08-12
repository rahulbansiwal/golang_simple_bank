package gapi

import (
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/db/util"
	"simple_bank/pb"
	"simple_bank/token"
	"simple_bank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	taskDistributor worker.TaskDistributor
	tokenMaker      token.Maker
}

func NewServer(config util.Config, store db.Store,taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmerticKey)
	if err != nil {
		return nil, fmt.Errorf("cant create token maker %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/db/util"
	"simple_bank/gapi"
	"simple_bank/mail"
	"simple_bank/pb"
	"simple_bank/worker"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Connection to DB can't be established", err)
	}
	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, store)
	go runGrpcServer(config, store, taskDistributor)
	runGatewayServer(config, store, taskDistributor)
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("Server can't be start", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listner, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Grpc listner can't be start", err)

	}
	log.Printf("Starting server at address %v", config.GRPCServerAddress)
	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("Grpc server can't be start", err)
	}
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskprocessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Printf("starting task process for: %s", taskprocessor)
	err := taskprocessor.Start()
	if err != nil {
		fmt.Printf("error while processing task:%s", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("Server can't be start", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Not able to start GRPC Gateway")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listner, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Grpc listner can't be start", err)

	}
	log.Printf("Starting HTTP GW server at address %s", config.HTTPServerAddress)
	err = http.Serve(listner, mux)
	if err != nil {
		log.Fatal("HTTP GW server can't be start", err)
	}
}

func runGinServer(config util.Config, store db.Store) {

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Server can't be start", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Server can't be start", err)
	}
}

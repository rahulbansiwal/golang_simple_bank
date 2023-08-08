package api

import (
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/db/util"
	"simple_bank/token"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmerticKey)
	if err != nil {
		return nil, fmt.Errorf("cant create token maker %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("supportedcurrency", validcurrency)
	}
	server.setupRoutes()

	return server, nil
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setupRoutes() {
	router := gin.Default()
	router.POST("/user", server.CreateUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/token/renew_access",server.renewAccessToken)
	
	authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccountFromId)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/transfers", server.CreateTransfer)

	server.router = router
}

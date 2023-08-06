package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required,supportedcurrency"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := s.validAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}
	authPaylaod := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPaylaod.Username != fromAccount.Owner{
		err := errors.New("from account dosent belong to logged in user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err)) 
		return
	}
	_, valid = s.validAccount(ctx, req.ToAccountId, req.Currency)
	if !valid {
			return
	}

	args := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}
	result, err2 := s.store.TransferTx(ctx, args)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err2))
		return
	}
	ctx.JSON(http.StatusOK, result)

}

func (s *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true

}

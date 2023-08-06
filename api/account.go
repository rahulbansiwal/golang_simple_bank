package api

import (
	"errors"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,supportedcurrency"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}
	acc, err2 := s.store.CreateAccount(ctx, args)
	if err2 != nil {
		if pqErr, ok := err2.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err2))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err2))
		return
	}
	ctx.JSON(http.StatusOK, acc)

}

type getAccountFromIdParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccountFromId(ctx *gin.Context) {
	var req getAccountFromIdParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	acc, err2 := s.store.GetAccount(ctx, req.ID)
	if err2 != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err2))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if acc.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to login user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type listAccountParams struct {
	PageId   int32 `form:"page_id"`
	PageSize int32 `form:"page_size"`
}

func (s *Server) listAccount(ctx *gin.Context) {
	var req listAccountParams
	var param db.ListAccountsParams
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if ctx.ShouldBindQuery(&req); req.PageId == 0 {
		param.Owner = authPayload.Username
		param.Limit = 5
		param.Offset = 0
	} else {
		param.Owner = authPayload.Username
		param.Limit = req.PageSize
		param.Offset = req.PageSize * (req.PageId - 1)
	}

	list, err := s.store.ListAccounts(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, list)

}

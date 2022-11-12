package handlers // для аккаунта

import (
	"HTTP-REST-API/internal/domain"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handlerImpl) GetUserBalance(ctx *gin.Context) {
	// знаем URL и контекст, можем обратиться к БД
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Errorf("parser error : %s", err)
	}
	account, err := h.db.GetBalance(id)
	if err != nil {
		fmt.Errorf("db error: %s", err)
	}
	ctx.IndentedJSON(http.StatusOK, account)
}

func (h *handlerImpl) AddToUserBalance(ctx *gin.Context) {
	account := domain.Account{
		Id:      0,
		Balance: 0,
	}
	err := ctx.BindJSON(&account)
	if err != nil {
		errors.New("wrong account formatting")
	}
	status, err := h.db.AddToBalance(account.Id, int(account.Balance))
	if err != nil {
		fmt.Errorf("bad request")
	}

	if status == domain.Deposit {
		ctx.IndentedJSON(http.StatusOK, account)
	} else {
		ctx.IndentedJSON(http.StatusBadRequest, nil)
	}

}

func (h *handlerImpl) ReserveUsersAmount(ctx *gin.Context) {

}

func (h *handlerImpl) AdmitPurchase(ctx *gin.Context) {

}

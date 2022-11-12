package handlers // Package handlers для аккаунта

import (
	"HTTP-REST-API/internal/domain"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (handler *handlerImpl) GetAccountBalance(ctx *gin.Context) {
	// знаем URL и контекст, можем обратиться к БД
	id, err := strconv.Atoi(ctx.Param("id"))
	fmt.Printf("%T", id)
	if err != nil {
		fmt.Errorf("parser error : %s", err)
	}
	account, err := handler.db.GetBalance(id)
	if err != nil {
		fmt.Errorf("db error: %s", err)
	}
	ctx.IndentedJSON(http.StatusOK, *account)
}

func (handler *handlerImpl) AddToAccountBalance(ctx *gin.Context) {
	account := domain.AccountDto{}
	err := ctx.BindJSON(&account) // TODO не обрабатывается неверный формат реквеста
	if err != nil {
		errors.New("wrong account formatting")
	}
	status, err := handler.db.AddToBalance(account.Id, int(account.BalanceAdded))
	if err != nil {
		fmt.Errorf("bad request")
	}

	if status == domain.Deposit {
		ctx.IndentedJSON(http.StatusOK, account)
	} else {
		ctx.IndentedJSON(http.StatusBadRequest, nil)
	}

}

func (handler *handlerImpl) ReserveUsersAmount(ctx *gin.Context) {

}

func (handler *handlerImpl) AdmitPurchase(ctx *gin.Context) {

}

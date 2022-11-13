package handlers // Package handlers для аккаунта

import (
	"HTTP-REST-API/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (handler *handlerImpl) GetAccountBalance(ctx *gin.Context) {
	// знаем URL и контекст, можем обратиться к БД
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	account, err := handler.db.GetBalance(id)

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, *account)
}

func (handler *handlerImpl) AddToAccountBalance(ctx *gin.Context) {
	account := domain.AccountDto{}
	err := ctx.BindJSON(&account)

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	if account.Id == 0 && account.BalanceAdded == 0 {
		err := errors.New("wrong json id format or adding zero balance")
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	status, err := handler.db.AddToBalance(account.Id, int(account.BalanceAdded))
	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err)
	}

	if status == domain.Deposit {
		ctx.IndentedJSON(http.StatusOK, account)
	} else {
		ctx.IndentedJSON(http.StatusBadRequest, nil)
	}

}

func (handler *handlerImpl) ReserveUsersAmount(ctx *gin.Context) {

	reserve := domain.ReservationDto{}
	checker := domain.ReservationDto{}

	err := ctx.BindJSON(&reserve)

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	if reserve == checker {
		err := errors.New("wrong json id format")
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	status, err := handler.db.ReserveAmount(reserve.AccountId, reserve.ServiceId, reserve.OrderId, int(reserve.Amount))

	if err != nil {
		ctx.String(http.StatusBadRequest, "Reservation status: %s\nSome error %s\n", status, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, reserve)

}

func (handler *handlerImpl) AdmitPurchase(ctx *gin.Context) {
	report := domain.ReportDto{}
	checker := domain.ReportDto{}
	err := ctx.BindJSON(&report)

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	if report == checker {
		ctx.String(http.StatusBadRequest, "%s", err)
		return
	}

	status, err := handler.db.Admit(report.AccountId, report.OrderId, report.ServiceId, int(report.Amount))

	if err != nil {
		ctx.String(http.StatusBadRequest, "Transaction status: %s\n Error:%s", status, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, report)
}

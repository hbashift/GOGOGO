package handlers // Package handlers для аккаунта

import (
	"HTTP-REST-API/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handlerImpl) GetAccountBalance(ctx *gin.Context) {
	// знаем URL и контекст, можем обратиться к БД
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.Error(err)
		return
	}

	account, err := h.service.GetBalance(id)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, account)
}

func (h *handlerImpl) AddToAccountBalance(ctx *gin.Context) {
	account := domain.AccountDto{}
	err := ctx.BindJSON(&account)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if account.Id < 0 || account.BalanceAdded == 0 {
		err := errors.New("wrong json id format or adding zero balance")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = h.service.AddToBalance(account.Id, int(account.BalanceAdded))

	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}

	ctx.Status(http.StatusOK)
}

func (h *handlerImpl) ReserveUsersAmount(ctx *gin.Context) {

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

	status, err := h.service.ReserveAmount(reserve.AccountId, reserve.ServiceId, reserve.OrderId, int(reserve.Amount))

	if err != nil {
		ctx.Error(err)
		err = errors.New("reservation status: " + status.String())
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx.Status(http.StatusCreated)

}

func (h *handlerImpl) AdmitPurchase(ctx *gin.Context) {
	report := domain.ReportDto{}
	checker := domain.ReportDto{}
	err := ctx.BindJSON(&report)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if report == checker {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	status, err := h.service.Admit(report.AccountId, report.OrderId, report.ServiceId, int(report.Amount))

	if err != nil {
		ctx.Error(err)
		err = errors.New("transaction status: " + status.String())
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx.Status(http.StatusOK)
}

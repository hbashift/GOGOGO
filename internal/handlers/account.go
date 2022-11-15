package handlers

import (
	"HTTP-REST-API/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAccountBalance returns domain.Account in JSON with Context if input data is valid
// else aborts operation with specified error
func (h *Handler) GetAccountBalance(ctx *gin.Context) {
	// знаем URL и контекст, можем обратиться к БД
	account := &domain.Account{}
	err := ctx.BindJSON(account)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := account.Id

	if id < 0 {
		er := errors.New("negative id")
		ctx.AbortWithError(http.StatusBadRequest, er)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	account, err = h.service.GetBalance(id)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, account)
}

// AddToAccountBalance returns http.StatusOK if operation was executed
// else aborts it with specified error
func (h *Handler) AddToAccountBalance(ctx *gin.Context) {
	account := domain.AccountDto{}
	err := ctx.BindJSON(&account)

	if account.Id < 0 || account.BalanceAdded < 0 {
		er := errors.New("negative account_id or balance_added")
		ctx.AbortWithError(http.StatusBadRequest, er)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if account.BalanceAdded == 0 {
		err := errors.New("adding zero balance")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = h.service.AddToBalance(account.Id, int(account.BalanceAdded))

	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// ReserveAccountsBalance returns http.StatusCreated if operation was executed
// else aborts it with specified error and http.Status
func (h *Handler) ReserveAccountsBalance(ctx *gin.Context) {

	reserve := domain.ReservationDto{}
	checker := domain.ReservationDto{}

	err := ctx.BindJSON(&reserve)

	if reserve.AccountId < 0 || reserve.Amount < 0 || reserve.OrderId < 0 || reserve.ServiceId < 0 {
		er := errors.New("negative account_id | amount | oreser_id | service_id")
		ctx.AbortWithError(http.StatusBadRequest, er)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if reserve == checker {
		err := errors.New("wrong json id format")
		ctx.AbortWithError(http.StatusBadRequest, err)
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

// AdmitPurchase returns http.StatusAccepted if operation was executed
// else aborts it with specified error and http.Status
func (h *Handler) AdmitPurchase(ctx *gin.Context) {
	report := domain.ReservationDto{}
	checker := domain.ReservationDto{}
	err := ctx.BindJSON(&report)

	if report.AccountId < 0 || report.Amount < 0 || report.OrderId < 0 || report.ServiceId < 0 {
		er := errors.New("negative account_id | amount | oreser_id | service_id")
		ctx.AbortWithError(http.StatusBadRequest, er)
		return
	}

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

	ctx.Status(http.StatusAccepted)
}

// Decline returns http.StatusAccepted if operation was executed
// else aborts it with specified error and http.Status
func (h *Handler) Decline(ctx *gin.Context) {
	decline := domain.DeclineDto{}
	checker := domain.DeclineDto{}

	err := ctx.BindJSON(&decline)

	if decline.OrderId < 0 || decline.ServiceId < 0 || decline.AccountId < 0 {
		er := errors.New("negative account_id | amount | order_id | service_id")
		ctx.AbortWithError(http.StatusBadRequest, er)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if decline == checker {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	status, err := h.service.Decline(decline.AccountId, decline.OrderId, decline.ServiceId)

	if err != nil {
		ctx.Error(err)
		err = errors.New("transaction status: " + status.String())
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx.Status(http.StatusAccepted)
}

// Transfer returns http.StatusAccepted if operation was executed
// else aborts it with specified error and http.Status
func (h *Handler) Transfer(ctx *gin.Context) {
	transfer := domain.Transfer{}
	err := ctx.BindJSON(&transfer)

	if transfer.Sender < 0 || transfer.Receiver < 0 || transfer.Amount < 0 {
		err = errors.New("negative sender_id | receiver_id | amount")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	status, err := h.service.Transfer(transfer.Sender, transfer.Receiver, int(transfer.Amount))

	if err != nil {
		ctx.Error(err)
		err = errors.New("transaction status: " + status.String())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusAccepted)
}

package handlers

import (
	"HTTP-REST-API/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
)

type Handler struct { // implements Handler
	service service.Service
}

func NewHandler(service *service.Service) (*Handler, error) {
	if &service == nil {
		return nil, errors.New("service was not created")
	}

	return &Handler{
		service: *service,
	}, nil
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/account/balance/:id", h.GetAccountBalance)
	router.POST("/account/balance/add", h.AddToAccountBalance)
	router.POST("/account/reservation", h.ReserveAccountsBalance)
	router.POST("/account/reservation/admit", h.AdmitPurchase)
	router.POST("/account/transfer", h.Transfer)
	router.POST("/reservation/decline", h.Decline)

	return router
}

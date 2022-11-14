package handlers

import (
	"HTTP-REST-API/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

/*func GetAccountBalance(w http.ResponseWriter, r *http.Request) {

}*/

type Handler interface {

	// TODO methods

	GetAccountBalance(ctx *gin.Context)
	AddToAccountBalance(ctx *gin.Context)
	ReserveUsersAmount(ctx *gin.Context)
	AdmitPurchase(ctx *gin.Context)
}

type handlerImpl struct { // implements Handler
	service *service.Service
}

func NewHandler(service *service.Service) (Handler, error) {
	if service == nil {
		return nil, fmt.Errorf("empty repository")
	}

	return &handlerImpl{
		service: service,
	}, nil
}

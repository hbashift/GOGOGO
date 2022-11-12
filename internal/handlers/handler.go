package handlers

import (
	"HTTP-REST-API/internal/domain/repository"
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

// TODO прочитать про Context
// TODO можно сделать разделение handler'ов

type handlerImpl struct { // implements Handler
	db repository.Repository
}

func NewHandler(repo repository.Repository) (Handler, error) {
	if repo == nil {
		return nil, fmt.Errorf("empty repository")
	}

	return &handlerImpl{
		db: repo,
	}, nil
}

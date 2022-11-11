package handlers

import (
	"HTTP-REST-API/internal/storage"
	"fmt"
	"github.com/gin-gonic/gin"
)

/*func GetUserBalance(w http.ResponseWriter, r *http.Request) {

}*/

type Handler interface {

	// TODO methods

	GetUserBalance(ctx *gin.Context)
	AddToUserBalance(ctx *gin.Context)
	ReserveUsersAmount(ctx *gin.Context)
	AdmitPurchase(ctx *gin.Context)
}

// TODO прочитать про Context
// TODO можно сделать разделение handler'ов

type handlerImpl struct { // implements Handler
	acсountDB     *storage.AccountDb
	reserveDB     *storage.ReserveDb
	transactionDB *storage.TransactionDb
}

func InitHandler(acc *storage.AccountDb, res *storage.ReserveDb) (Handler, error) {
	if acc == nil {
		return nil, fmt.Errorf("some error")
	}

	if res == nil {
		return nil, fmt.Errorf("some error")
	}

	return &handlerImpl{
		acсountDB:     acc,
		reserveDB:     res,
		transactionDB: nil,
	}, nil
}

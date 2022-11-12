package middleware

import (
	"HTTP-REST-API/internal/domain/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Middleware struct {
}

func GetBalance(ctx *gin.Context, repository repository.Repository) {
	id := 1
	account, err := repository.GetBalance(id)
	if err != nil {
		fmt.Errorf("db error: %s", err)
	}

	ctx.IndentedJSON(http.StatusOK, account)
}

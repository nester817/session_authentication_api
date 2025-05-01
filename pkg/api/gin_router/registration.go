package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	customer "github.com/nester817/session_authentication_api.git/pkg/storage"
)

func (h *Handler) AddProfile(ctx *gin.Context) {
	var customer customer.Customer
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(customer.Password) < 6 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "incorrect password",
		})
	}

	hashPassword, err := HashPassword([]byte(customer.Password))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	customer.Password = string(hashPassword)

	if err := h.sqlDb.InsertCustomer(ctx.Request.Context(), &customer); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer)
}

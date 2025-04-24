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

	if err := h.db.InsertCustomer(ctx.Request.Context(), customer); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer)
}

func (h *Handler) DeleteProfile(ctx *gin.Context) {
	var customerConfiguration customer.Customer
	if err := ctx.ShouldBindJSON(&customerConfiguration); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.db.DeleteCustomer(
		ctx.Request.Context(),
		customerConfiguration.Email,
		customerConfiguration.Password,
	); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, customerConfiguration.Email)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/internal/service"
)

type TransactionHandler struct {
	transactionSvc service.TransactionService
}

func NewTransactionHandler(
	transactionSvc service.TransactionService,
) *TransactionHandler {
	return &TransactionHandler{
		transactionSvc: transactionSvc,
	}
}

// @Id listTransactions
// @Tags transaction
// @Summary list transactions
// @Description list transactions
// @Router /v1/transactions [get]
// @Param user_id query string true "user uuid"
// @Param from query string false "from wallet"
// @Param to query string false "to wallet"
// @Param start query string false "start time"
// @Param end query string false "end time"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @version 1.0
// @Success 200 {object} model.TransactionPage
func (h *TransactionHandler) List(c *gin.Context) {
	userID := c.Query("user_id")
	from := c.Query("from")
	to := c.Query("to")
	start := c.Query("start")
	end := c.Query("end")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	var fromWallet, toWallet *string
	if len(from) != 0 {
		fromWallet = &from
	}
	if len(to) != 0 {
		toWallet = &to
	}
	var startTime, endTime *int64
	if len(start) != 0 {
		if v, err := strconv.ParseInt(start, 10, 64); err == nil {
			startTime = &v
		}
	}
	if len(end) != 0 {
		if v, err := strconv.ParseInt(start, 10, 64); err == nil {
			endTime = &v
		}
	}

	page, err := h.transactionSvc.List(c.Request.Context(), userID, model.TransactionSearchCondition{
		From:   fromWallet,
		To:     toWallet,
		Start:  startTime,
		End:    endTime,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, page)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/internal/service"
)

type WalletHandler struct {
	walletSvc      service.WalletService
	transactionSvc service.TransactionService
}

func NewWalletHandler(
	walletSvc service.WalletService,
	transactionSvc service.TransactionService,
) *WalletHandler {
	return &WalletHandler{
		walletSvc:      walletSvc,
		transactionSvc: transactionSvc,
	}
}

// @Id createWallet
// @Tags wallet
// @Summary create wallet
// @Description create wallet
// @Router /v1/wallets [post]
// @Param walletInfo body handler.Create.createWalletDto true "wallet information"
// @version 1.0
// @Success 201 {object} model.Wallet
func (h *WalletHandler) Create(c *gin.Context) {
	type createWalletDto struct {
		UserID  string  `json:"userID"`
		Balance float64 `json:"balance"`
	}
	var dto createWalletDto
	if err := c.BindJSON(&dto); err != nil {
		_ = c.Error(err)
		return
	}

	wallet, err := h.walletSvc.Create(c.Request.Context(), model.Wallet{
		UserID:  dto.UserID,
		Balance: dto.Balance,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// @Id listWallet
// @Tags wallet
// @Summary list wallet
// @Description list wallet
// @Router /v1/wallets [get]
// @Param user_id query string true "user uuid"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @version 1.0
// @Success 200 {object} model.WalletPage
func (h *WalletHandler) List(c *gin.Context) {
	userID := c.Query("user_id")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	page, err := h.walletSvc.List(c.Request.Context(), userID, limit, offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Id transferFund
// @Tags wallet
// @Summary transfer fund
// @Description transfer fund
// @Router /v1/wallets/transfer [post]
// @Param transferRequestBody body handler.Transfer.transferDto true "transfer information"
// @version 1.0
// @Success 201 {object} model.Transaction
func (h *WalletHandler) Transfer(c *gin.Context) {
	type transferDto struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	var dto transferDto
	if err := c.BindJSON(&dto); err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := h.transactionSvc.Transfer(c.Request.Context(), dto.From, dto.To, dto.Amount)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

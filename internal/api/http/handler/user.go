package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/internal/service"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(
	userSvc service.UserService,
) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
	}
}

// @Id createUser
// @Tags user
// @Summary create user
// @Description create user
// @Router /v1/users [post]
// @Param userInfo body handler.Create.createUserDto true "user information"
// @version 1.0
// @Success 201 {object} model.User
func (h *UserHandler) Create(c *gin.Context) {
	type createUserDto struct {
		Name string `json:"name"`
	}
	var dto createUserDto
	if err := c.BindJSON(&dto); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := h.userSvc.Create(c.Request.Context(), model.User{
		Name: dto.Name,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

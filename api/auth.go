package api

import (
	"fmt"
	"gb-detection/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.services.GetUserByID(user_id)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, Response{"success", "", fmt.Sprintf("%b", u)})
}

func (h *Handler) Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	u := model.User{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := h.services.LoginCheck(u.Username, u.Password)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Username or password is incorrect.")
		return
	}

	c.JSON(http.StatusOK, Response{"", token, ""})

}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	u := model.User{}

	u.Username = input.Username
	u.Password = input.Password

	_, err := h.services.SaveUser(&u)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, Response{"Registration success", "", fmt.Sprintf("%b", u)})
}

type Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Data    string `json:"data"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, Response{message, "", ""})
}

func newSuccessResponse(method string, name string) {
	if name == "" {
		logrus.Printf("Succesful request for %s", method)
	} else {
		logrus.Printf("Succesful request for %s - %s", method, name)
	}
}

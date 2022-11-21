package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createChatRequest struct {
	Token string `json:"token" binding:"required"`
}

func (s *Server) createChat(ctx *gin.Context) {
	var req createChatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("token can not be empty")))
		return
	}

	payload := struct {
		Token string `json:"token"`
	}{
		Token: req.Token,
	}

	err := s.pushToQueue(payload, "applications.chat.create")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := JsonResponse{
		Status:  true,
		Message: "Request Received",
	}

	ctx.JSON(http.StatusOK, resp)
}

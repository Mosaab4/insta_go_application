package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createMessageRequest struct {
	Token      string `json:"token"`
	ChatNumber int64  `json:"chat_number"`
	Body       string `json:"body"`
}

func (s *Server) createMessage(ctx *gin.Context) {
	var req createMessageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			errorResponse(errors.New("token, chat number and body can not be empty")),
		)
		return
	}

	payload := struct {
		Token      string `json:"token"`
		ChatNumber int64  `json:"chat_number"`
		Body       string `json:"body"`
	}{
		Token:      req.Token,
		ChatNumber: req.ChatNumber,
		Body:       req.Body,
	}

	err := s.pushToQueue(payload, "applications.message.create")

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

package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"instabug/go/event"
	"log"
)

type Server struct {
	router *gin.Engine
	Rabbit *amqp.Connection
}

type JsonResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func NewServer(rabbit *amqp.Connection) *Server {
	server := &Server{
		Rabbit: rabbit,
	}
	router := gin.Default()

	routes := router.Group("/api")
	{
		routes.POST("/chats", server.createChat)
		routes.POST("/messages", server.createMessage)
	}

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (s *Server) pushToQueue(payload any, severity string) error {
	emitter, err := event.NewEventEmitter(s.Rabbit, "applications")

	if err != nil {
		return err
	}

	j, _ := json.Marshal(payload)

	err = emitter.Push(string(j), severity)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

package api

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"fmt"
	_ "log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	db "github.com/junimslage10/gofinance-backend-user/db/sqlc"
	// util "github.com/junimslage10/gofinance-backend-user/util"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	// errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	// if errOnValiteToken != nil {
	// 	return
	// }
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	hashedInput := sha512.Sum512_256([]byte(req.Password))
	trimmedHash := bytes.Trim(hashedInput[:], "\x00")
	preparedPassword := string(trimmedHash)
	passwordHashInBytes, err := bcrypt.GenerateFromPassword([]byte(preparedPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	var passwordHashed = string(passwordHashInBytes)
	arg := db.CreateUserParams{
		Username: req.Username,
		Password: passwordHashed,
		Email:    req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	// Apache Kafka ---
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka-hostname:29092"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "create-user"
	messageText := &db.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	messageTextJson, _ := json.Marshal(messageText)
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(string(messageTextJson)),
	}, nil)

	// Wait for message deliveries before shutting down
	// Flush and close the producer and the events channel
	for p.Flush(1000) > 0 {
		fmt.Print("Still waiting to flush outstanding messages\n", p)
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required"`
}

func (server *Server) getUser(ctx *gin.Context) {
	// errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	// if errOnValiteToken != nil {
	// 	return
	// }
	var req getUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Apache Kafka ---
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka-hostname:29092"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "get-user-username"
	messageText := &db.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	messageTextJson, _ := json.Marshal(messageText)
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(string(messageTextJson)),
	}, nil)

	// Wait for message deliveries before shutting down
	// Flush and close the producer and the events channel
	for p.Flush(1000) > 0 {
		fmt.Print("Still waiting to flush outstanding messages\n", p)
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserByIdRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) getUserById(ctx *gin.Context) {
	// errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	// if errOnValiteToken != nil {
	// 	return
	// }
	var req getUserByIdRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUserById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Apache Kafka ---
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka-hostname:29092"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "get-user-id"
	messageText := &db.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	messageTextJson, _ := json.Marshal(messageText)
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(string(messageTextJson)),
	}, nil)

	// Wait for message deliveries before shutting down
	// Flush and close the producer and the events channel
	for p.Flush(1000) > 0 {
		fmt.Print("Still waiting to flush outstanding messages\n", p)
	}
	ctx.JSON(http.StatusOK, user)

}

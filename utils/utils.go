package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"github.com/golang-jwt/jwt/v5"
)

type ApiResponderType struct {
	Data interface{}
	StatusCode *int
	Error interface{}
}

func ErrorHandler(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ApiResponder(input ApiResponderType) map[string]interface{}{
	result := make(map[string]interface{})

	// default values
	data := "Successfully processed request"
	statusCode := 200
	result["data"] = data
	result["statusCode"] = &statusCode
	result["error"] = "null"

	// setting values if sent and returning response
	if input.Data != nil { result["data"] = input.Data }
	if input.StatusCode != nil { result["statusCode"] = *input.StatusCode }
	if input.Error != nil { result["error"] = input.Error }
	return result
}

func StrToInt(value string) int {
	convertedNumber, err := strconv.Atoi(value)
	ErrorHandler(err)
	return convertedNumber
}

func GetLimitAndOffset(value string) (int, int) {
	page := StrToInt(value)
	limit := 20
	offset := (page - 1) * limit

	if page < 1 { return limit, 0 }

	return limit, offset
}

type UserToken struct {
	Username string
	Email string
}

func NewUserToken(username, email string) *UserToken {
	return &UserToken{
		Username: username,
		Email: email,
	}
}

func CreateToken(user UserToken) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"email": user.Email,
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("Invalid Token")
	}

	return nil
}

func RabbitMQ(receiverChanMap *map[string]chan string) (*amqp.Channel, amqp.Queue, amqp.Queue) {
	// connect to rabbit mq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ErrorHandler(err)
	//defer conn.Close()

	// creating channel
	ch, err := conn.Channel()
	ErrorHandler(err)
	//defer ch.Close()

	// creating queues
	reqQ, err := ch.QueueDeclare("data_request_queue", false, false, false, false, nil)
	ErrorHandler(err)
	resQ, err := ch.QueueDeclare("data_response_queue", false, false, false, false, nil)
	ErrorHandler(err)

	// setting up consumer
	go func() {
		msgs, err := ch.Consume(resQ.Name, "", false, false, false, false, nil)
		ErrorHandler(err)

		for msg := range msgs {
			msg.Ack(false)

			myReceiverChanMap := *receiverChanMap
			myReceiverChanMap[msg.CorrelationId] <- string(msg.Body)
		}
	}()

	return ch, reqQ, resQ
}

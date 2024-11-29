package controllers

import (
	"encoding/json"
	"fmt"
	"myapp/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

type ControllerRequest struct {
	Ch *amqp.Channel
	ResQ *amqp.Queue
	ReqQ *amqp.Queue
	ReceiverChanMap map[string]chan string
	Validate *validator.Validate
}

type Payload struct {
	Entity string
	Operation string 
	Payload interface{}
}

func GetAndSendRequest(cr *ControllerRequest, payload string) string {
	corrId := GenerateCorrId()
	cr.ReceiverChanMap[string(corrId)] = make(chan string)

	err := cr.Ch.Publish("", cr.ReqQ.Name, false, false, amqp.Publishing {
		ContentType: "text/plain",
		Body: []byte(payload),
		ReplyTo: cr.ResQ.Name,
		CorrelationId: corrId,
	}) 
	utils.ErrorHandler(err)
	var result string = <-cr.ReceiverChanMap[corrId]
	delete(cr.ReceiverChanMap, corrId)

	return result
}

func GenerateCorrId() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func GeneratePayload(payload Payload) string {
	if payload.Payload != nil {
		data, err := json.Marshal(payload.Payload)
		utils.ErrorHandler(err)
		return fmt.Sprintf(`{"entity":"%s","operation":"%s","payload":%s}`, payload.Entity, payload.Operation, data)
	}

	return fmt.Sprintf(`{"entity":"%s","operation":"%s"}`, payload.Entity, payload.Operation)
}

func GenerateJsonResponse(data string, isArray bool) interface{} {

	// converting response to hashmap for proper json response
	var response interface{}
	if isArray {
		response = []map[string]interface{}{}
	} else {
		response = map[string]interface{}{}
	}
	err := json.Unmarshal([]byte(data), &response)
	utils.ErrorHandler(err)
	return response
}

func SendValidationError(c echo.Context) error {
	errCode := 400
	return c.JSON(errCode, utils.ApiResponder(utils.ApiResponderType{Data: "null", Error: "Validation error occurred", StatusCode: &errCode}))
}

func SendSuccessResponse(c echo.Context) error {
	statusCode := 200
	return c.JSON(statusCode, utils.ApiResponder(utils.ApiResponderType{}))
}

func SendCustomSuccessResponse(c echo.Context, data interface{}) error {
	statusCode := 200
	return c.JSON(statusCode, utils.ApiResponder(utils.ApiResponderType{Data: data, StatusCode: &statusCode}))
}

func MapRequest(c echo.Context, reqBody interface{}) {
	err := c.Bind(reqBody)
	utils.ErrorHandler(err)
}


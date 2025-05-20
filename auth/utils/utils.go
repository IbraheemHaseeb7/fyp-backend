package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"

	"github.com/IbraheemHaseeb7/pubsub"
)

var Requests map[string](chan pubsub.PubsubMessage)

func StructToMap(data interface{}) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    v := reflect.ValueOf(data)
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        // Get the JSON tag
        jsonTag := fieldType.Tag.Get("json")
        if jsonTag != "" && jsonTag != "-" { // Exclude fields with "-" JSON tag
            result[jsonTag] = field.Interface()
        }
    }

    return result, nil
}

// InternalApiRequest is a struct that represents the request to the internal API
type InternalApiRequest struct {
	Endpoint string
	Method   string
	Body     map[string]any
	Headers  map[string]string
}

func NewInternalApiRequest(endpoint string, method string, body map[string]any, headers map[string]string) *InternalApiRequest {
	return &InternalApiRequest{
		Endpoint: endpoint,
		Method:   method,
		Body:     body,
		Headers:  headers,
	}
}

var client = &http.Client{}
func (r *InternalApiRequest) InternalCall() error {
	
	// Make json body if provided
	var reqBody []byte
	if r.Body != nil {
		jsonBody, err := json.Marshal(r.Body)
		if err != nil {
			return err
		}

		reqBody = jsonBody
	}
	bodyReader := bytes.NewReader(reqBody)

	// Create a new request
	environment := os.Getenv("ENVIRONMENT")
	if environment == "local" {
		r.Endpoint = "http://localhost:8000" + r.Endpoint
	} else if environment == "staging" {
		r.Endpoint = "http://" + os.Getenv("BASE_ADDRESS") + ":" + os.Getenv("PORT") + r.Endpoint
	} else if environment == "production" {
		r.Endpoint = "https://" + os.Getenv("BASE_ADDRESS") + r.Endpoint
	} else {
		r.Endpoint = os.Getenv("BASE_ADDRESS") + r.Endpoint
	}
	req, err := http.NewRequest(r.Method, r.Endpoint, bodyReader)
	if err != nil {
		return err
	}

	// Set the request body if provided
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	// Set the request body if provided
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("Error: " + resp.Status)
	}

	return nil
}


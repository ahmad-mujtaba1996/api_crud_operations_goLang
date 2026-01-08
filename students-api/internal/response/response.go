package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(err validator.ValidationErrors) Response {

	var erroMsg []string
	for _, e := range err {
		switch e.ActualTag() {
		case "required":
			erroMsg = append(erroMsg, fmt.Sprintf("field %s is required", e.Field()))
		case "email":
			erroMsg = append(erroMsg, e.Field()+" must be a valid email")
		case "gte":
			erroMsg = append(erroMsg, e.Field()+" must be greater than or equal to "+e.Param())
		case "lte":
			erroMsg = append(erroMsg, e.Field()+" must be less than or equal to "+e.Param())
		default:
			erroMsg = append(erroMsg, e.Field()+" is not valid")
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(erroMsg, ", "),
	}
}

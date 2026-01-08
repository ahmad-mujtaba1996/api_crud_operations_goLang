package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/response"
	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/storage"
	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/types"
	"github.com/go-playground/validator/v10"
)

// Storage is passed to the handler to interact with the database. It's called dependency injection.
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a Student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body"))) //fmt.Errorf("empty body") for creating error with custom message
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrors := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrors))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User created with id ", slog.String("id", fmt.Sprintf("%d", lastId)))

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "OK", "id": fmt.Sprintf("%d", lastId)}) //fmt.Sprintf to convert int64 to string
	}
}

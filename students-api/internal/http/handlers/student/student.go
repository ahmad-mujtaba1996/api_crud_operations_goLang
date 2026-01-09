package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

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

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting student by id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}
		student, error := storage.GetStudentById(intId)
		if error != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(error))
			return
		}
		response.WriteJSON(w, http.StatusOK, student)

	}
}

func GetStudentsList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting list of students")
		studentList, err := storage.GetStudentsList()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJSON(w, http.StatusOK, studentList)
	}
}

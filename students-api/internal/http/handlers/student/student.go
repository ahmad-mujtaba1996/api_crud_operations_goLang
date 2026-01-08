package student

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/types"
)

func Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		json.NewDecoder(r.Body).Decode(&student)

		slog.Info("Creating a Student")

		w.Write([]byte("Welcome to students api"))
	}
}

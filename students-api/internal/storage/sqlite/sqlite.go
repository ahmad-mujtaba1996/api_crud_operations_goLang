package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/config"
	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT,
   email TEXT,
   age INTEGER 
   )`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) CreateStudent(name, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ? , ?)") // Adding question marks to avoid SQL injection attacks.
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err // 0 is empty id in case of error.
	}

	return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with id %d not found", fmt.Sprint(id)) // Checking if no rows were returned.
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

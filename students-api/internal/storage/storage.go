package storage

import "github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error) // returns the ID and error of the created student
	GetStudentById(id int64) (types.Student, error)
	GetStudentsList() ([]types.Student, error)
	UpdateStudent(id int64, name string, email string, age int) (types.Student, error)
}

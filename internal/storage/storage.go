package storage

import (
	"context"

	"github.com/pchintan243/golang/internal/types"
)

type Storage interface {
	CreateStudent(ctx context.Context, name string, email string, age int) (int64, error)
	GetStudentById(ctx context.Context, id int64) (types.Student, error)
	GetStudents(ctx context.Context) ([]types.Student, error)
	DeleteStudentById(ctx context.Context, id int64) (string, error)
	UpdateStudent(ctx context.Context, id int64, name string, email string, age int) (types.Student, error)
}

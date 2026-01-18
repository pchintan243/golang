package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pchintan243/golang/internal/config"
	"github.com/pchintan243/golang/internal/types"
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
	age INTEGER)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, err
}

func (s *Sqlite) CreateStudent(ctx context.Context, name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(ctx, name, email, age)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(ctx context.Context, id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRowContext(ctx, id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents(ctx context.Context) ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students")

	if err != nil {
		slog.Error("Error1")
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil

}

func (s *Sqlite) DeleteStudentById(ctx context.Context, id int64) (string, error) {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return "error occurred", fmt.Errorf("failed to prepare delete: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(ctx, id)

	if err != nil {
		return "error occurred", fmt.Errorf("failed to execute delete: %w", err)
	}

	// 3. Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "error occurred", fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return "no student found", fmt.Errorf("no student found with id %d", id)
	}

	return "delete successfully", nil
}

func (s *Sqlite) UpdateStudent(ctx context.Context, id int64, name string, email string, age int) (types.Student, error) {
	stmt, err := s.Db.PrepareContext(ctx, "UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")

	if err != nil {
		return types.Student{}, err
	}

	// When use prepare must need to close it
	defer stmt.Close()

	// short hand prop
	// No defer needed! Go handles the statement lifecycle automatically.
	// result, err := s.Db.Exec("UPDATE students SET name = ? WHERE id = ?", name, id)

	result, err := stmt.ExecContext(ctx, name, email, age, id)

	if err != nil {
		return types.Student{}, err
	}
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id %d", id)
	}

	return types.Student{
		Id:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}, nil

}

package database

import (
	"database/sql"
	"github.com/google/uuid"
)

type CourseDb struct {
	db *sql.DB
	ID string
	Name string
	Description string
	CategoryID string
}

func NewCourseDb(db *sql.DB) *CourseDb {
	return &CourseDb{db: db}
}

func (c *CourseDb) Create(name, description, categoryID string) (*CourseDb, error) {
	id := uuid.New().String()
	query := `INSERT INTO courses (id, name, description, category_id) VALUES ($1, $2, $3, $4)`
	_, err := c.db.Exec(query, id, name, description, categoryID)
	if err != nil {
		return nil, err
	}
	return &CourseDb{
		ID: id,
		Name: name,
		Description: description,
		CategoryID: categoryID,
	}, nil
}

func (c *CourseDb) FindAll() ([]CourseDb, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []CourseDb{}
	for rows.Next() {
		var id, name, descrpition, categoryID string
		if err := rows.Scan(&id, &name, &descrpition, &categoryID); err != nil {
			return nil, err
		}
		courses = append(courses, CourseDb{ID: id, Name: name, Description: descrpition,  CategoryID: categoryID})
	}
	return courses, nil
}

func (c* CourseDb) FindByCategoryID(categoryID string) ([]CourseDb, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses WHERE category_id = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	courses := []CourseDb{}
	for rows.Next() {
		var id, name, description, categoryID string
		if err := rows.Scan(&id, &name, &description, &categoryID); err != nil {
			return nil, err
		}
		courses = append(courses, CourseDb{
			ID: id,
			Name: name,
			CategoryID: categoryID,
			Description: description,
		})
	}
	return courses, nil
}
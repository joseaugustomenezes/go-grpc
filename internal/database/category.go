package database

import (
	"database/sql"
	"github.com/google/uuid"
)

type CategoryDb struct {
	db *sql.DB
	ID string
	Name string
	Description string
}

func NewCategoryDb(db *sql.DB) *CategoryDb {
	return &CategoryDb{db: db}
}

func (c *CategoryDb) Create(name string, description string) (*CategoryDb, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO categories(id, name, description) VALUES ($1, $2, $3)", id, name, description)
	if err != nil {
		return &CategoryDb{}, err
	}
	return &CategoryDb{ID: id, Name: name, Description: description}, nil
}

func (c *CategoryDb) FindAll() ([]CategoryDb, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []CategoryDb{}
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			return nil, err
		}
		categories = append(categories, CategoryDb{ID: id, Name: name, Description: description})
	}
	return categories, nil
}

func (c *CategoryDb) FindByCourseID(courseID string) (*CategoryDb, error) {
	var id, name, description string
	err := c.db.QueryRow("SELECT c.id, c.name, c.description FROM categories c JOIN courses co ON c.id = co.category_id WHERE co.id = $1", courseID).Scan(&id, &name, &description)
	if err != nil {
		return nil, err
	}
	return &CategoryDb{
		ID: id,
		Name: name,
		Description: description,
	}, nil
}

func (c *CategoryDb) FindByID(id string) (*CategoryDb, error) {
	var name, description string
	err := c.db.QueryRow("SELECT name, description FROM categories WHERE id = $1", id).Scan(&name, &description)
	if err != nil {
		return nil, err
	}
	return &CategoryDb{
		ID: id,
		Name: name,
		Description: description,
	}, nil
}
package store

import (
	models "github.com/charlieroth/godo/internal/domain"
)

type Store interface {
	Add(task *models.Task) error
	Update(task *models.Task) error
	Delete(id int) error
	List() ([]*models.Task, error)
	Get(id int) (*models.Task, error)
	NextID() int
	Load(dbPath string) error
	Save(dbPath string) error
}

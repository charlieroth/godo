package store

import (
	"encoding/json"
	"os"

	models "github.com/charlieroth/godo/internal/domain"
)

type JsonStore struct {
	dbPath string
	tasks  []*models.Task
}

func NewJsonStore(dbPath string) *JsonStore {
	return &JsonStore{
		dbPath: dbPath,
		tasks:  []*models.Task{},
	}
}

func (s *JsonStore) Load() error {
	jsonFile, err := os.Open(s.dbPath)
	if os.IsNotExist(err) {
		_, err := os.Create(s.dbPath)
		if err != nil {
			return err
		}
		return nil
	}

	defer jsonFile.Close()

	byteValue, err := os.ReadFile(s.dbPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &s.tasks)
	if err != nil {
		return err
	}

	return nil
}

func (s *JsonStore) Save() error {
	jsonFile, err := os.Create(s.dbPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	json.NewEncoder(jsonFile).Encode(s.tasks)

	return nil
}

func (s *JsonStore) Add(task *models.Task) error {
	s.tasks = append(s.tasks, task)
	return nil
}

func (s *JsonStore) Update(task *models.Task) error {
	for i, t := range s.tasks {
		if t.ID == task.ID {
			s.tasks[i] = task
			return nil
		}
	}

	return nil
}

func (s *JsonStore) Delete(id int) error {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}

	return nil
}

func (s *JsonStore) List() ([]*models.Task, error) {
	return s.tasks, nil
}

func (s *JsonStore) Get(id int) (*models.Task, error) {
	for _, t := range s.tasks {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, nil
}

func (s *JsonStore) NextID() int {
	return len(s.tasks) + 1
}

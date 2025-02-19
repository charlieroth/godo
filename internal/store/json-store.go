package store

import (
	"encoding/json"
	"fmt"
	"os"

	models "github.com/charlieroth/godo/internal/domain"
)

type JsonStore struct {
	dbPath string
	tasks  map[int]models.Task
	nextID int
}

func NewJsonStore(dbPath string) *JsonStore {
	return &JsonStore{
		dbPath: dbPath,
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *JsonStore) Load() error {
	jsonFile, err := os.Open(s.dbPath)
	if os.IsNotExist(err) {
		file, err := os.Create(s.dbPath)
		if err != nil {
			return err
		}

		file.Write([]byte("{}"))
		file.Close()

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

func (s *JsonStore) Add(task models.Task) error {
	s.tasks[s.nextID] = task
	s.nextID++
	return nil
}

func (s *JsonStore) Update(task models.Task) error {
	if _, ok := s.tasks[task.ID]; !ok {
		return fmt.Errorf("task with id %d not found", task.ID)
	}

	s.tasks[task.ID] = task
	return nil
}

func (s *JsonStore) Delete(id int) error {
	delete(s.tasks, id)
	return nil
}

func (s *JsonStore) List() ([]models.Task, error) {
	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *JsonStore) Get(id int) (models.Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, fmt.Errorf("task with id %d not found", id)
	}
	return task, nil
}

func (s *JsonStore) NextID() int {
	return s.nextID
}

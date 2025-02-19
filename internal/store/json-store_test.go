package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/charlieroth/godo/internal/domain"
)

func TestLoadNonExistentFile(t *testing.T) {
	// Create a temporary directory and construct a path to a file that does not exist.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "nonexistent.json")

	store := NewJsonStore(dbPath)
	err := store.Load()
	if err != nil {
		t.Fatalf("expected no error when loading a nonexistent file, got: %v", err)
	}

	// Check that the file was created by Load.
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("expected file to be created, but it does not exist")
	}

	// After loading from a non-existent file, the tasks map should be empty.
	tasks, err := store.List()
	if err != nil {
		t.Fatalf("error listing tasks: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestLoadCorruptFile(t *testing.T) {
	// Create a temporary file and write invalid JSON into it.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db.json")
	corruptData := []byte("this is not valid json")
	err := os.WriteFile(dbPath, corruptData, 0644)
	if err != nil {
		t.Fatalf("error writing corrupt file: %v", err)
	}

	store := NewJsonStore(dbPath)
	err = store.Load()
	if err == nil {
		t.Fatalf("expected an error loading corrupt file, got nil")
	}
}

func TestAddAndGet(t *testing.T) {
	// Create a temporary store.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db.json")
	store := NewJsonStore(dbPath)

	// Add a task.
	task := domain.Task{
		ID:    store.NextID(),
		Title: "Test task",
		Done:  false,
	}
	if err := store.Add(task); err != nil {
		t.Fatalf("unexpected error adding task: %v", err)
	}

	// Verify nextID is incremented.
	if got, want := store.NextID(), 2; got != want {
		t.Errorf("expected nextID to be %d, got %d", want, got)
	}

	// Retrieve the task.
	gotTask, err := store.Get(1)
	if err != nil {
		t.Fatalf("unexpected error getting task with id 1: %v", err)
	}
	if gotTask.Title != "Test task" {
		t.Errorf("expected task title 'Test task', got '%s'", gotTask.Title)
	}
}

func TestUpdate(t *testing.T) {
	// Create a temporary store and add a task.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db.json")
	store := NewJsonStore(dbPath)

	task := domain.Task{
		ID:    store.NextID(),
		Title: "Original Task",
		Done:  false,
	}
	store.Add(task)

	// Update the task.
	task.Title = "Updated Task"
	task.Done = true
	if err := store.Update(task); err != nil {
		t.Fatalf("unexpected error updating task: %v", err)
	}

	// Verify that the update took effect.
	updatedTask, err := store.Get(task.ID)
	if err != nil {
		t.Fatalf("unexpected error retrieving updated task: %v", err)
	}
	if updatedTask.Title != "Updated Task" || !updatedTask.Done {
		t.Errorf("task not updated correctly. Got %+v", updatedTask)
	}

	// Try updating a task that doesn't exist.
	nonExistentTask := domain.Task{
		ID:    999,
		Title: "Non-existent",
		Done:  false,
	}
	if err := store.Update(nonExistentTask); err == nil {
		t.Errorf("expected error when updating non-existent task, got nil")
	}
}

func TestDelete(t *testing.T) {
	// Create a temporary store and add two tasks.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db.json")
	store := NewJsonStore(dbPath)

	task1 := domain.Task{
		ID:    store.NextID(),
		Title: "Task 1",
		Done:  false,
	}
	store.Add(task1)
	task2 := domain.Task{
		ID:    store.NextID(),
		Title: "Task 2",
		Done:  false,
	}
	store.Add(task2)

	// Delete the first task.
	if err := store.Delete(task1.ID); err != nil {
		t.Fatalf("error deleting task: %v", err)
	}

	// Verify task1 is deleted.
	if _, err := store.Get(task1.ID); err == nil {
		t.Errorf("expected an error when retrieving deleted task, got nil")
	}

	// Verify task2 still exists.
	gotTask, err := store.Get(task2.ID)
	if err != nil {
		t.Fatalf("unexpected error retrieving task2: %v", err)
	}
	if gotTask.Title != "Task 2" {
		t.Errorf("expected task2 title 'Task 2', got '%s'", gotTask.Title)
	}
}

func TestList(t *testing.T) {
	// Create a temporary store.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db.json")
	store := NewJsonStore(dbPath)

	// Initially, the task list should be empty.
	list, err := store.List()
	if err != nil {
		t.Fatalf("error listing tasks: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d tasks", len(list))
	}

	// Add tasks.
	task1 := domain.Task{
		ID:    store.NextID(),
		Title: "Task 1",
		Done:  false,
	}
	store.Add(task1)
	task2 := domain.Task{
		ID:    store.NextID(),
		Title: "Task 2",
		Done:  true,
	}
	store.Add(task2)

	// Verify that the list contains both tasks.
	list, err = store.List()
	if err != nil {
		t.Fatalf("error listing tasks after adding: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 tasks in list, got %d", len(list))
	}
	foundTask1, foundTask2 := false, false
	for _, task := range list {
		switch task.ID {
		case task1.ID:
			foundTask1 = true
		case task2.ID:
			foundTask2 = true
		}
	}
	if !foundTask1 || !foundTask2 {
		t.Errorf("list missing tasks: foundTask1=%v, foundTask2=%v", foundTask1, foundTask2)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory and file.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db.json")

	// Create a store, add tasks, then save.
	store := NewJsonStore(dbPath)
	task1 := domain.Task{
		ID:    store.NextID(),
		Title: "Persisted Task 1",
		Done:  false,
	}
	store.Add(task1)
	task2 := domain.Task{
		ID:    store.NextID(),
		Title: "Persisted Task 2",
		Done:  true,
	}
	store.Add(task2)
	if err := store.Save(); err != nil {
		t.Fatalf("error saving store: %v", err)
	}

	// Create a new store instance and load from the same file.
	newStore := NewJsonStore(dbPath)
	if err := newStore.Load(); err != nil {
		t.Fatalf("error loading store: %v", err)
	}

	// The loaded store should contain the two tasks.
	list, err := newStore.List()
	if err != nil {
		t.Fatalf("error listing tasks from loaded store: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 tasks after load, got %d", len(list))
	}

	// Check that each task is loaded correctly.
	loadedTask1, err := newStore.Get(task1.ID)
	if err != nil {
		t.Fatalf("error retrieving task1 after load: %v", err)
	}
	if loadedTask1.Title != task1.Title || loadedTask1.Done != task1.Done {
		t.Errorf("task1 not loaded correctly, got %+v", loadedTask1)
	}

	loadedTask2, err := newStore.Get(task2.ID)
	if err != nil {
		t.Fatalf("error retrieving task2 after load: %v", err)
	}
	if loadedTask2.Title != task2.Title || loadedTask2.Done != task2.Done {
		t.Errorf("task2 not loaded correctly, got %+v", loadedTask2)
	}

	// Note: The NextID field is not persisted and will be reset in a new store.
	// If this behavior is not desired, you may want to adjust the implementation.
}

func TestJSONEncodingConsistency(t *testing.T) {
	// Test that saving the store produces valid JSON that can be decoded
	// into the same map structure.

	// Create a store and add a couple of tasks.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db.json")
	store := NewJsonStore(dbPath)
	task := domain.Task{
		ID:    store.NextID(),
		Title: "Encoding Test Task",
		Done:  false,
	}
	store.Add(task)
	if err := store.Save(); err != nil {
		t.Fatalf("error saving store: %v", err)
	}

	// Read the file contents.
	data, err := os.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("error reading saved file: %v", err)
	}

	// Unmarshal into a new map.
	var decodedTasks map[string]domain.Task
	err = json.Unmarshal(data, &decodedTasks)
	if err != nil {
		t.Fatalf("error unmarshaling JSON: %v", err)
	}

	// Check that the decoded tasks contain our task.
	if len(decodedTasks) != 1 {
		t.Errorf("expected 1 task in saved JSON, got %d", len(decodedTasks))
	}
	for _, dt := range decodedTasks {
		if dt.Title != task.Title || dt.Done != task.Done {
			t.Errorf("decoded task does not match, got %+v", dt)
		}
	}
}

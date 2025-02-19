package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charlieroth/godo/internal/domain"
	"github.com/charlieroth/godo/internal/store"
)

type AddCmd struct {
	Args    []string
	FlagSet *flag.FlagSet
}

type ListCmd struct {
	Args    []string
	FlagSet *flag.FlagSet
}

type DoCmd struct {
	Args    []string
	FlagSet *flag.FlagSet
}

type UndoCmd struct {
	Args    []string
	FlagSet *flag.FlagSet
}

type DeleteCmd struct {
	Args    []string
	FlagSet *flag.FlagSet
}

type App struct {
	Store *store.JsonStore
	AddCmd
	ListCmd
	DoCmd
	UndoCmd
	DeleteCmd
}

func NewApp(store *store.JsonStore) *App {
	addCmdFlagSet := flag.NewFlagSet("add", flag.ExitOnError)
	addCmdFlagSet.String("title", "Unnamed Task", "Title of task to add.")

	listCmdFlagSet := flag.NewFlagSet("list", flag.ExitOnError)

	doCmdFlagSet := flag.NewFlagSet("do", flag.ExitOnError)
	doCmdFlagSet.Int("id", 0, "ID of task to mark as done.")

	undoCmdFlagSet := flag.NewFlagSet("undo", flag.ExitOnError)
	undoCmdFlagSet.Int("id", 0, "ID of task to mark as undone.")

	deleteCmdFlagSet := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmdFlagSet.Int("id", 0, "ID of task to delete.")

	return &App{
		Store: store,
		AddCmd: AddCmd{
			FlagSet: addCmdFlagSet,
		},
		ListCmd: ListCmd{
			FlagSet: listCmdFlagSet,
		},
		DoCmd: DoCmd{
			FlagSet: doCmdFlagSet,
		},
		UndoCmd: UndoCmd{
			FlagSet: undoCmdFlagSet,
		},
		DeleteCmd: DeleteCmd{
			FlagSet: deleteCmdFlagSet,
		},
	}
}

func (a *App) Run(args []string) error {
	if len(args) < 2 {
		Help()
		os.Exit(1)
	}

	switch args[1] {
	case "add":
		err := a.Add(args)
		if err != nil {
			log.Fatal(err)
		}
	case "list":
		err := a.List(args)
		if err != nil {
			log.Fatal(err)
		}
	case "do":
		err := a.Do(args)
		if err != nil {
			log.Fatal(err)
		}
	case "undo":
		err := a.Undo(args)
		if err != nil {
			log.Fatal(err)
		}
	case "delete":
		err := a.Delete(args)
		if err != nil {
			log.Fatal(err)
		}
	case "help":
		Help()
		os.Exit(0)
	default:
		Help()
		os.Exit(1)
	}

	return nil
}

func (a *App) Add(args []string) error {
	if err := a.AddCmd.FlagSet.Parse(args[2:]); err != nil {
		return fmt.Errorf("error parsing add command: %w", err)
	}

	title := a.AddCmd.FlagSet.Lookup("title").Value.String()

	task := domain.Task{
		ID:    a.Store.NextID(),
		Title: title,
		Done:  false,
	}

	err := a.Store.Add(task)
	if err != nil {
		return err
	}

	err = a.Store.Save()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) List(args []string) error {
	if err := a.ListCmd.FlagSet.Parse(args[2:]); err != nil {
		return fmt.Errorf("error parsing list command: %w", err)
	}

	tasks, err := a.Store.List()
	if err != nil {
		return err
	}
	fmt.Printf("Tasks: %v\n", tasks)

	for _, task := range tasks {
		if task.Done {
			fmt.Printf("[x] %d. %s\n", task.ID, task.Title)
		} else {
			fmt.Printf("[ ] %d. %s\n", task.ID, task.Title)
		}
	}

	return nil
}

func (a *App) Do(args []string) error {
	if err := a.DoCmd.FlagSet.Parse(args[2:]); err != nil {
		return fmt.Errorf("error parsing do command: %w", err)
	}

	taskIDString := a.DoCmd.FlagSet.Lookup("id").Value.String()
	taskID, err := strconv.Atoi(taskIDString)
	if err != nil {
		return fmt.Errorf("error converting task ID to int: %w", err)
	}

	task, err := a.Store.Get(taskID)
	if err != nil {
		return err
	}

	task.Done = true
	err = a.Store.Update(task)
	if err != nil {
		return err
	}

	err = a.Store.Save()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Undo(args []string) error {
	if err := a.UndoCmd.FlagSet.Parse(args[2:]); err != nil {
		return fmt.Errorf("error parsing undo command: %w", err)
	}

	taskIDString := a.UndoCmd.FlagSet.Lookup("id").Value.String()
	taskID, err := strconv.Atoi(taskIDString)
	if err != nil {
		return fmt.Errorf("error converting task ID to int: %w", err)
	}

	task, err := a.Store.Get(taskID)
	if err != nil {
		return err
	}

	task.Done = false
	err = a.Store.Update(task)
	if err != nil {
		return err
	}

	err = a.Store.Save()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Delete(args []string) error {
	if err := a.DeleteCmd.FlagSet.Parse(args[2:]); err != nil {
		return fmt.Errorf("error parsing delete command: %w", err)
	}

	taskIDString := a.DeleteCmd.FlagSet.Lookup("id").Value.String()
	taskID, err := strconv.Atoi(taskIDString)
	if err != nil {
		return fmt.Errorf("error converting task ID to int: %w", err)
	}

	err = a.Store.Delete(taskID)
	if err != nil {
		return err
	}

	err = a.Store.Save()
	if err != nil {
		return err
	}

	return nil
}

func Help() {
	fmt.Println("Usage: godo <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add     Add a new task")
	fmt.Println("  list    List all tasks")
	fmt.Println("  do      Mark a task as done")
	fmt.Println("  undo    Mark a task as undone")
	fmt.Println("  delete  Delete a task")
	fmt.Println("  help    Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  godo add -title \"Buy groceries\"")
	fmt.Println("  godo list")
	fmt.Println("  godo do -id 1")
	fmt.Println("  godo undo -id 1")
	fmt.Println("  godo delete -id 1")
}

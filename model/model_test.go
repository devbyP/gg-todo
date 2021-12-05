package model

import (
	"fmt"
	"testing"
)

func TestTodo(t *testing.T) {
	testName := "testing testing"
	_, err := AddTodo(testName)
	if err != nil {
		t.Error("error at addTodo.")
	}
	defer deleteTodoByName(testName)
	_, err = getTodos()
	if err != nil {
		t.Error("error getting todos")
	}
}

func TestTagSearch(t *testing.T) {
	tags, err := SearchTags("G")
	fmt.Println((*tags[0]).Name)
	fmt.Println(len(tags))
	if err != nil {
		t.Error("cannot search")
	}
	if !(len(tags) > 0) || tags[0].Name != "GO" {
		t.Error("not get correct tag.")
	}
}

func TestSelectFunc(t *testing.T) {
	todo, err := selectMultipleTodo("SELECT id, name, created_on FROM todos WHERE checked_on IS NULL")
	if err != nil {
		t.Errorf("Error in selectMultipleTodo: %v\n", err)
	}
	if len(todo) <= 0 {
		t.Errorf("Error get no rows selectMultipleTodo sql")
	}
}

func TestTagAndColor(t *testing.T) {
	todo, err := getTodos()
	if err != nil {
		t.Errorf("Error get todo: %v\n", err)
	}
	if len(todo) <= 0 {
		t.Error("no todo found.")
	}
	id := todo[0].ID
	_, err = getTagsAndHighlightInTodo(id)
	if err != nil {
		t.Errorf("no todo tags: %v\n", err)
	}
	_, err = getTagById(id)
	if err != nil {
		t.Error("error get Tag")
	}
	_, err = getColorById(id)
	if err != nil {
		t.Error("error get Color")
	}
}

func TestFullTodos(t *testing.T) {
	todo, _ := getTodos()
	_, err := getFullTodo(todo[0].ID)
	if err != nil {
		t.Errorf("Error get full todolist: %v\n", err)
	}
}

package model

import (
	"fmt"
	"testing"
)

func TestTodo(t *testing.T) {
	testName := "testing testing"
	_, err := addTodo(testName)
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

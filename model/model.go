package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	db = conDatabase()
}

const (
	dbname   = "todo"
	user     = "postgres"
	password = "postgres"
	port     = "5432"
	host     = "localhost"
)

func conDatabase() *sql.DB {
	dbInfo := fmt.Sprintf(`host=%s port=%s user=%s 
		password=%s dbname=%s sslmode=disable`,
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("database connected.")
	return db
}

type fullScanner interface {
	fullScan(scanner) error
}

type scanner interface {
	Scan(...interface{}) error
}

type Todo struct {
	ID        int
	Name      string
	CreatedOn time.Time
}

func addTodo(todo string) (int64, error) {
	result, err := db.Exec("INSERT INTO todos(name) VALUES($1)",
		todo)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getTodos() ([]*Todo, error) {
	rows, err := db.Query("SELECT id, name, created_on FROM todos")
	todoList := make([]*Todo, 0)
	if err != nil {
		return todoList, err
	}
	defer rows.Close()
	for rows.Next() {
		todo := new(Todo)
		todo.fullScan(rows)
		todoList = append(todoList, todo)
	}
	return todoList, nil
}

func getTodo(id int) (*Todo, error) {
	todo := new(Todo)
	row := db.QueryRow("SELECT id, name, created_on FROM todos WHERE id=$1", id)
	err := todo.fullScan(row)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func deleteTodoByName(name string) {
	result, err := db.Exec("DELETE FROM todos WHERE name=$1", name)
	if err != nil {
		log.Fatal(err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		log.Fatal("no row deleted.")
	}
}

func (t *Todo) fullScan(rows scanner) error {
	return rows.Scan(&t.ID, &t.Name, &t.CreatedOn)
}

type Tag struct {
	ID   int
	Name string
}

func SearchTags(searchText string) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	rows, err := db.Query("SELECT id, name FROM tags WHERE name like '%' || $1 || '%'", searchText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tag := new(Tag)
		tag.fullScan(rows)
		tags = append(tags, tag)
	}
	return tags, nil
}

func AddTag(name string) (int64, error) {
	result, err := db.Exec("INSERT INTO tags(name) VALUES($1)", name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getTagId(name string) (int, error) {
	var id int
	row := db.QueryRow("SELECT id FROM tags WHERE name=$1", name)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (t *Tag) fullScan(rows scanner) error {
	return rows.Scan(&t.ID, &t.Name)
}

type HighlightColor struct {
	ID   int
	Name string
	Hex  string
}

func SelectHighlight() ([]*HighlightColor, error) {
	var colors []*HighlightColor = make([]*HighlightColor, 0)
	rows, err := db.Query("SELECT id, name, hex FROM highlight")
	if err != nil {
		return colors, err
	}
	defer rows.Close()
	for rows.Next() {
		color := new(HighlightColor)
		color.fullScan(rows)
		colors = append(colors, color)
	}
	return colors, nil
}

func (hl *HighlightColor) fullScan(rows scanner) error {
	return rows.Scan(&hl.ID, &hl.Name, &hl.Hex)
}

type TodoTag struct {
	Todo      int
	Tag       int
	Highligth int
}

func getTagsAndHighlightInTodo(todoID int) ([]*TodoTag, error) {
	todoTags := make([]*TodoTag, 0)
	rows, err := db.Query("SELECT todo_id, tag_id, highlight_id FROM todotags WHERE todo_id=$1", todoID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		tt := new(TodoTag)
		tt.fullScan(rows)
		todoTags = append(todoTags, tt)
	}
	return todoTags, nil
}

func (tt *TodoTag) fullScan(rows scanner) error {
	return rows.Scan(&tt.Todo, &tt.Tag, &tt.Highligth)
}

type FullTodo map[string]interface{}

func NewFullTodo(todo *Todo, tags []*TodoTag) FullTodo {
	return FullTodo{
		"id":         todo.ID,
		"name":       todo.Name,
		"created_on": todo.CreatedOn,
		"tags":       tags,
	}
}

func GetFullTodo(todoID int) (FullTodo, error) {
	todo, err := getTodo(todoID)
	if err != nil {
		return nil, err
	}
	tags, err := getTagsAndHighlightInTodo(todo.ID)
	if err != nil {
		return nil, err
	}
	todoItem := NewFullTodo(todo, tags)
	return todoItem, nil
}

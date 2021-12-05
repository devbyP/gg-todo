package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"
  "os"
  "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
  godotenv.Load(".env")
  connectionInfo := dbConnectionInfo{
    dbName: os.Getenv("DB_NAME"),
    user: os.Getenv("DB_USER"),
    password: os.Getenv("DB_PASSWORD"),
    port: os.Getenv("DB_PORT"),
    host: os.Getenv("HOSTNAME"),
  }
	db = conDatabase(connectionInfo)
}

type dbConnectionInfo struct {
  dbName string
  user string
  password string
  port string
  host string
}

func conDatabase(info dbConnectionInfo) *sql.DB {
	dbInfo := fmt.Sprintf(
    `host=%s
    port=%s
    user=%s 
		password=%s
    dbname=%s
    sslmode=disable`,
		info.host,
    info.port,
    info.user,
    info.password,
    info.dbName,
  )
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
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"created_on"`
	CheckedOn time.Time `json:"checked_on"`
}

func selectMultipleTodo(sql string, args ...interface{}) ([]*Todo, error) {
	todos := make([]*Todo, 0)
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		todo := new(Todo)
		err = todo.fullScan(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func selectOneTodo(sql string, args ...interface{}) (*Todo, error) {
	todo := new(Todo)
	row := db.QueryRow(sql, args...)
	err := todo.fullScan(row)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func AddTodo(todo string) (int, error) {
	row := db.QueryRow("INSERT INTO todos(name) VALUES($1) RETURNING id",
		todo)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getTodos() ([]*Todo, error) {
	return selectMultipleTodo("SELECT id, name, created_on FROM todos WHERE checked_on IS NULL")
}

func getTodo(id int) (*Todo, error) {
	return selectOneTodo("SELECT id, name, created_on FROM todos WHERE id=$1 AND checked_on IS NULL", id)
}

func DeleteTodoById(id int) error {
	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func selectMultipleTag(sql string, args ...interface{}) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tag := new(Tag)
		err = tag.fullScan(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func selectOneTag(sql string, args ...interface{}) (*Tag, error) {
	tag := new(Tag)
	row := db.QueryRow(sql, args...)
	err := tag.fullScan(row)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func GetNumOfTag(limit int) ([]*Tag, error) {
	if limit <= 0 {
		limit = 10
	}
	return selectMultipleTag("SELECT id, name FROM tags LIMIT $1", limit)
}

func SearchTags(searchText string) ([]*Tag, error) {
	return selectMultipleTag("SELECT id, name FROM tags WHERE name like '%' || $1 || '%'", searchText)
}

func AddTag(name string) (int, error) {
	row := db.QueryRow("INSERT INTO tags(name) VALUES($1) RETURNING id", name)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getTagById(tagId int) (*Tag, error) {
	return selectOneTag("SELECT id, name FROM tags WHERE id=$1", tagId)
}

func (t *Tag) fullScan(rows scanner) error {
	return rows.Scan(&t.ID, &t.Name)
}

type HighlightColor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

func selectMultipleColor(sql string, args ...interface{}) ([]*HighlightColor, error) {
	colors := make([]*HighlightColor, 0)
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		color := new(HighlightColor)
		err = color.fullScan(rows)
		if err != nil {
			return nil, err
		}
		colors = append(colors, color)
	}
	return colors, nil
}

func selectOneColor(sql string, args ...interface{}) (*HighlightColor, error) {
	color := new(HighlightColor)
	row := db.QueryRow(sql, args...)
	if err := color.fullScan(row); err != nil {
		return nil, err
	}
	return color, nil
}

func GetAllHighlight() ([]*HighlightColor, error) {
	return selectMultipleColor("SELECT id, name, hex FROM highlightcolor")
}

func getColorById(id int) (*HighlightColor, error) {
	return selectOneColor("SELECT id, name, hex FROM highlightcolor WHERE id=$1", id)
}

func (hl *HighlightColor) fullScan(rows scanner) error {
	return rows.Scan(&hl.ID, &hl.Name, &hl.Hex)
}

type TodoTag struct {
	Todo      int `json:"todo_id"`
	Tag       int `json:"tag_id"`
	Highligth int `json:"highlight_id"`
}

func InsertTodoTag(todo, tag, color int) error {
	var err error
	if color == 0 {
		_, err = db.Exec(`INSERT INTO todotags(todo_id, tag_id) VALUES($1, $2)`, todo, tag)
	} else {
		_, err = db.Exec(`INSERT INTO todotags(
			todo_id, tag_id, highlight_id
		) VALUES($1, $2, $3)`, todo, tag, color)
	}
	if err != nil {
		return err
	}
	return nil
}

func selectMultipleTodoTags(sql string, args ...interface{}) ([]*TodoTag, error) {
	todoTags := make([]*TodoTag, 0)
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		tt := new(TodoTag)
		if err := tt.fullScan(rows); err != nil {
			return nil, err
		}
		todoTags = append(todoTags, tt)
	}
	return todoTags, nil
}

func selectOneTodoTag(sql string, args ...interface{}) (*TodoTag, error) {
	todoTag := new(TodoTag)
	row := db.QueryRow(sql, args...)
	if err := todoTag.fullScan(row); err != nil {
		return nil, err
	}
	return todoTag, nil
}

func getTagsAndHighlightInTodo(todoID int) ([]*TodoTag, error) {
	return selectMultipleTodoTags("SELECT todo_id, tag_id, highlight_id FROM todotags WHERE todo_id=$1", todoID)
}

func getOneTagInTodo(todoId, tagId int) (*TodoTag, error) {
	return selectOneTodoTag("SELECT todo_id, tag_id, highlight_id FROM todotags WHERE todo_id=$1 AND tag_id=$2", todoId, tagId)
}

func getTodoThatIncludeTag(tagId int) ([]*TodoTag, error) {
	return selectMultipleTodoTags("SELECT todo_id, tag_id, highlight_id FROm todotags WHERE tag_id=$1", tagId)
}

func (tt *TodoTag) fullScan(rows scanner) error {
	return rows.Scan(&tt.Todo, &tt.Tag, &tt.Highligth)
}

type FullTodo map[string]interface{}

type fullTag map[string]string

func NewFullTodo(todo *Todo, tags []fullTag) FullTodo {
	return FullTodo{
		"id":         todo.ID,
		"name":       todo.Name,
		"created_on": todo.CreatedOn,
		"tags":       tags,
	}
}

func getFullTodo(todoID int) (FullTodo, error) {
	todo, err := getTodo(todoID)
	if err != nil {
		return nil, err
	}
	todotags, err := getTagsAndHighlightInTodo(todo.ID)
	if err != nil {
		return nil, err
	}
	fTags := make([]fullTag, 0)
	for i := 0; i < len(todotags); i++ {
		t := make(fullTag)
		tag, err := getTagById(todotags[i].Tag)
		if err != nil {
			return nil, err
		}
		color, err := getColorById(todotags[i].Highligth)
		if err != nil {
			return nil, err
		}
		t["tag_name"] = tag.Name
		t["color"] = color.Name
		t["color_hex"] = color.Hex
		fTags = append(fTags, t)
	}
	todoItem := NewFullTodo(todo, fTags)

	return todoItem, nil
}

func GetAllFullTodo() ([]FullTodo, error) {
	var fullTodos []FullTodo
	todos, err := getTodos()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(todos); i++ {
		todo := todos[i]
		fullTodo, err := getFullTodo(todo.ID)
		if err != nil {
			return nil, err
		}
		fullTodos = append(fullTodos, fullTodo)
	}
	return fullTodos, nil
}

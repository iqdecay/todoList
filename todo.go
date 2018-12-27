package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


type Todo struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Creation    time.Time `json:"created_at"`
	Due         time.Time `json:"due_date"`
	// time.Time() type has no null value since it's a Struct type, so we can't use omitempty
}

const filename = "tasklist.txt"

type TodoList []Todo

func (t *TodoList) save() error {
	data, err := json.MarshalIndent(t, "", "	")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	return ioutil.WriteFile(filename, []byte(data), 0600)
}

func loadTodoList() (TodoList) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Reading tasklist failed : %s", err)
	}
	var todos TodoList
	if err = json.Unmarshal(data, &todos); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	return todos
}

func addTodo(list TodoList, t Todo) TodoList {
	list = append(list, t)
	return list
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	todos := loadTodoList()
	renderTemplate(w, "add", &todos)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	todos := loadTodoList()
	renderTemplate(w, "view", &todos)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	creation := time.Now().AddDate(3, 0, 0)
	description := r.FormValue("description")
	todo := Todo{Title: title, Description: description, Creation: creation, Due: time.Now().AddDate(1, 1, 1)}
	todos := loadTodoList()
	todos = addTodo(todos, todo)
	todos.save()
	http.Redirect(w, r, "/view/", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, t *TodoList) {

	templates := template.Must(template.ParseFiles("view.html", "add.html"))
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	a := Todo{Title: "task 1", Description: "", Creation: time.Now()}
	b := Todo{"task 2", "perform task 2", time.Now(), time.Now().AddDate(1, 0, 0)}
	t := TodoList{a, b}
	t.save()
	v := loadTodoList()
	fmt.Println(v)
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}



// TODO :


package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const timeFormat = "2006-01-02"

type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Creation    string `json:"created_at,omitempty"`
	Due         string `json:"due_date,omitempty"`
	id          int    `json:"unique_id"`
}

const filename = "tasklist.txt"

type TodoList struct {
	list []Todo
	maxId int
}

func (t *TodoList) save() error {
	data, err := json.MarshalIndent(t, "", "	")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	return ioutil.WriteFile(filename, []byte(data), 0600)
}

func loadTodoList() TodoList {
	var todos TodoList
	// if the file doesn't exist, the tasklist is empty
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return todos
	} else {
		// otherwise we process the contained data
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Reading tasklist failed : #{err}")
		}
		if err = json.Unmarshal(data, &todos); err != nil {
			log.Fatalf("JSON unmarshaling failed: %s", err)
		}
		return todos
	}

}

func addTodo(todos TodoList, t Todo) TodoList {
	if t.id > todos.maxId {
		todos.maxId = t.id
	} else if t.id == todos.maxId {
		panic("maxId not coherent")
	}
	todos.list = append(todos.list, t)
	return todos
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	// it is only used to display the form, treatment is made in the saveHandler
	todos := loadTodoList()
	renderTemplate(w, "add", &todos)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// display the current tasklist
	todos := loadTodoList()
	renderTemplate(w, "view", &todos)

}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	nowString := time.Now().Format(timeFormat)
	now, _ := time.Parse(timeFormat, nowString)
	creation := now.Format(timeFormat)
	description := r.FormValue("description")
	dueString := r.FormValue("due")
	log.Println(dueString)
	dueDate, _ := time.Parse(timeFormat, dueString)
	pressedButton := r.FormValue("submit")

	if pressedButton == "returnButton" {
		http.Redirect(w, r, "/view/", http.StatusFound)
		return
	}
	log.Println("before :", dueDate.Before(now))
	if dueDate.Before(now) && pressedButton == "saveButton" { // Check if the dueDate makes sense

		http.Redirect(w, r, "/add/", http.StatusFound)
	} else {
		todos := loadTodoList()
		id := todos.maxId + 1
		todo := Todo{title, description, creation, dueString, id}
		todos = addTodo(todos, todo)
		todos.save()
		http.Redirect(w, r, "/view/", http.StatusFound)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, t *TodoList) {
	templates := template.Must(template.ParseFiles("view.html", "add.html"))
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

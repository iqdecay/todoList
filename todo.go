package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Date struct {
	Day, Month, Year string
}

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

func stringToTask(s string) Todo {
	lastCommaIndex := -1
	fieldIndex := 0
	var fieldValue string
	todo := Todo{}
	for index, char := range s {
		if index == len(s)-1 || char == ';' {
			fieldIndex += 1
			if fieldIndex != 4 {
				fieldValue = s[lastCommaIndex+1 : index]
			} else {
				fieldValue = s[lastCommaIndex+1:]
			}
			// This is messy but I don't have any other ideas
			switch fieldIndex {
			case 1:
				todo.Title = fieldValue
			case 2:
				todo.Description = fieldValue
			case 3:
				todo.Creation = time.Now()
			case 4:
				todo.Due = time.Now()
			}
			lastCommaIndex = index
		}
	}
	return todo
}

func (d *Date) convertToString() string {
	return d.Day + "/" + d.Month + "/" + d.Year
}

func stringToDate(s string) Date {
	lastSlashIndex := -1
	var date []string
	for index, char := range s {
		if char == '/' || index == len(s)-1 {
			if index == len(s)-1 {
				date = append(date, s[lastSlashIndex+1:])
			} else {
				date = append(date, s[lastSlashIndex+1:index])
			}
			lastSlashIndex = index
		}
	}
	d, m, y := date[0], date[1], date[2]
	return Date{d, m, y}
}

func (t TodoList) buildRep() string {
	var b strings.Builder
	for _, todo := range t {
		// Write title
		title := todo.Title + ";"
		(&b).Grow(len(title))
		_, _ = (&b).Write([]byte(title))
		// Write Description
		mission := todo.Description + ";"
		(&b).Grow(len(mission))
		_, _ = (&b).Write([]byte(mission))
		// Write creationDate
		creation := todo.Creation.String() + ";"
		(&b).Grow(len(creation))
		_, _ = (&b).Write([]byte(creation))
		// Write dueDate
		due := todo.Due.String() + "\n"
		(&b).Grow(len(due))
		_, _ = (&b).Write([]byte(due))
	}
	// Write the whole todoList
	return (&b).String()
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

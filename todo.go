package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	// Think about importing time
)

type Date struct {
	Day, Month, Year string
}

type Todo struct {
	Title       string
	Description string
	Creation    Date
	Due         Date
}

const filename = "tasklist"

type TodoList []Todo

func (t *TodoList) save() error {
	content := t.buildRep()
	return ioutil.WriteFile(filename, []byte(content), 0600)
}

func loadTodoList() TodoList {
	file, _ := ioutil.ReadFile(filename)
	reader := string(file)
	lastNewlineIndex := -1
	var todos TodoList
	var todo Todo
	for index, char := range reader {
		if index == len(reader)-1 || char == '\n' {
			if index == len(reader)-1 {
				todo = stringToTask(reader[lastNewlineIndex+1:])
			} else {
				todo = stringToTask(reader[lastNewlineIndex+1 : index])
			}
			lastNewlineIndex = index
			todos = append(todos, todo)
		}
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
				todo.Creation = stringToDate(fieldValue)
			case 4:
				todo.Due = stringToDate(fieldValue)
			}
			lastCommaIndex = index
		}
	}
	return todo
}

func (d Date) convertToString() string {
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
		creation := todo.Creation.convertToString() + ";"
		(&b).Grow(len(creation))
		_, _ = (&b).Write([]byte(creation))
		// Write dueDate
		due := todo.Due.convertToString() + "\n"
		(&b).Grow(len(due))
		_, _ = (&b).Write([]byte(due))
	}
	// Write the whole todo
	return (&b).String()
}

func addTodo(list TodoList, t Todo) TodoList {
	list = append(list, t)
	return list
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	newTodo := r.FormValue("body")
	fmt.Println(newTodo)
	d1 := stringToDate("01/05/1997")
	d2 := stringToDate("07/12/2018")
	a := Todo{"task 3", "perform task 3", d1, d2}
	e := TodoList{}
	e = addTodo(e, a)
	e.save()
	renderTemplate(w, "edit", &e)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	todos := loadTodoList()
	renderTemplate(w, "view", &todos)
}

func renderTemplate(w http.ResponseWriter, tmpl string, t *TodoList) {

	templates := template.Must(template.ParseFiles("view.html","add.html"))
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	d1 := stringToDate("01/05/1997")
	d2 := stringToDate("07/12/2018")
	a := Todo{"task 1", "perform task 1", d1, d2}
	b := Todo{"task 2", "perform task 2", d2, d1}
	t := TodoList{a, b}
	t.save()
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

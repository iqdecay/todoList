package main

import (
    "io/ioutil"
    "strings"
    "strconv"
//    "reflect"
    "fmt"
    "time"
    )

type Date struct {
    Day, Month, Year int
}
type Todo struct {
	Title       string
	TimeLeft    int // Numbers of day available for completion
	Description string
    Creation time.time
}

const filename = "tasklist"

type TodoList []Todo

func (t *TodoList) save() error {
    content := t.buildRep()
	return ioutil.WriteFile(filename, []byte(content), 0600)
}

func loadTodoList() (TodoList) {
    file, _ := ioutil.ReadFile(filename)
    fileAsString := string(file)
    fmt.Println(fileAsString)

    return TodoList{}

}


func (t TodoList) buildRep() string {
    var b strings.Builder
	for _, todo := range t {
        title := todo.Title + "\n"
        (&b).Grow(len(title))
        _, _ = (&b).Write([]byte(title))
        daysLeft := strconv.Itoa(todo.TimeLeft)+"\n"
        (&b).Grow(len(daysLeft))
        _, _ = (&b).Write([]byte(daysLeft))
        mission := todo.Description + "\n"
        (&b).Grow(len(mission))
        _, _ = (&b).Write([]byte(mission))
	}
    return (&b).String()
}
/*



func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)

}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	content := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(content)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we extract the page title from the Request
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		// And call the provided handler "fn"
		fn(w, r, m[2])
	}
}
*/
func main() {
    fmt.Println("Program running")
    a := Todo{"task1",1,"this is task 1"}
    b := Todo{"task2",31,"this is task 2"}
    c := TodoList{a, b}
    c.save()
    loadTodoList()
}

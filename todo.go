package main

import (
	"fmt"
	"io/ioutil"
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
	// Implement it correctly
	file, _ := ioutil.ReadFile(filename)
	fileAsString := string(file)
	reader := strings.NewReader(fileAsString)
	fmt.Println(reader)


}

func (d Date) convertToString() string {
	return d.Day + "/" + d.Month + "/" + d.Year
}

func convertToDate(s string) Date {
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
		title := todo.Title + "\n"
		(&b).Grow(len(title))
		_, _ = (&b).Write([]byte(title))
		// Write dueDate
		dueDate := todo.Due.convertToString()
		(&b).Grow(len(dueDate))
		_, _ = (&b).Write([]byte(dueDate))
		// Write creationDate
		creationDate := todo.Creation.convertToString()
		(&b).Grow(len(creationDate))
		_, _ = (&b).Write([]byte(creationDate))
		// Write Description
		mission := todo.Description + "\n"
		(&b).Grow(len(mission))
		_, _ = (&b).Write([]byte(mission))
	}
	// Write the whole todo
	return (&b).String()
}

/*
Note : the entered date in the HTML form should satisfy the following Regexp :
"(0[1-9])|(1[012])"


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
	d1 := convertToDate("01/05/1997")
	d2 := convertToDate("07/12/2018")
	a := Todo{"task 1", "perform task 1", d1, d2}
	b := Todo{"task 2", "perform task 2", d2, d1}
	t := TodoList{a, b}
	t.save()
	c := loadTodoList()
	fmt.Println(c)
}

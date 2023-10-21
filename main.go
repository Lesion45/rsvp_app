package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Rsvp struct {
	Name 		string
	Email 		string
	Phone 		string
	WillAttend  bool
}

type formData struct {
	*Rsvp
	Errors 		[]string
	}
	

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() error {
	templateNames := [5]string{"welcome", "form", "list", "sorry", "thanks"}

	for _, templateName := range templateNames {
		t, err := template.ParseFiles("templates/layout.html", "templates/" + templateName + ".html")
		if err != nil {
			return err
		} else {
			templates[templateName] = t
			fmt.Println("Template loaded: ", templateName)
		}
	}

	return nil
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
	templates["form"].Execute(writer, formData {Rsvp: &Rsvp{}, Errors: []string {}})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Rsvp {
			Name: request.Form["name"][0],
			Email: request.Form["email"][0],
			Phone: request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}
		errors := []string {}
		if responseData.Name == "" {
			errors = append(errors, "Please enter your name")
			}
		if responseData.Email == "" {
			errors = append(errors, "Please enter your email address")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter your phone number")
		}
		if len(errors) > 0 {
			templates["form"].Execute(writer, formData {
			Rsvp: &responseData, Errors: errors,
		})} else {
			responses = append(responses, &responseData)
			if responseData.WillAttend {
			templates["thanks"].Execute(writer, responseData.Name)
			} else {
			templates["sorry"].Execute(writer, responseData.Name)
			}}
	}		
}	

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

func listHandler(writer http.ResponseWriter, request * http.Request) {
	templates["list"].Execute(writer, responses)
}

func main() {
	err := loadTemplates()
	if err != nil {
		panic(err)
	}
	fmt.Println("Ok!")

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	err = http.ListenAndServe(":5000", nil)
	if (err != nil) {
		fmt.Println(err)
	}

}
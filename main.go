package main

import (
	"encoding/csv"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
)

// Global, because we need to share this with the handler functions
var client *messagebird.Client

// Data is for passing data to our templates
type Data struct {
	Name    string
	Phone   string
	Message string
}

// Routes
func landing(w http.ResponseWriter, r *http.Request) {
	RenderDefaultTemplate(w, "views/landing.gohtml", nil)
}

func callme(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	msgBody := "You have a new lead: " + r.FormValue("name") + ". Call them at " + r.FormValue("phone")
	agentToCall := getRandomContact(getPhoneList("sample-contacts.csv"))
	msg, err := sms.Create(
		client,
		"MBCars",
		[]string{agentToCall},
		msgBody,
		nil,
	)
	if err != nil {
		log.Println(err)
		RenderDefaultTemplate(
			w,
			"views/landing.gohtml",
			&Data{
				Name:    r.FormValue("name"),
				Phone:   r.FormValue("phone"),
				Message: "We couldn't complete your request! Please try again.",
			},
		)
		return
	}
	log.Println(msg)
	RenderDefaultTemplate(w, "views/sent.gohtml", nil)
}

func main() {
	client = messagebird.New("<enter-your-api-key>")

	// Routes
	http.HandleFunc("/", landing)
	http.HandleFunc("/callme", callme)

	// Serve
	port := ":8080"
	log.Println("Serving application on", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println(err)
	}
}

// RenderDefaultTemplate takes:
// - a http.ResponseWriter
// - an array of strings to contain a list of template files to render
// - data to render to the template. If no data, should enter 'nil'
func RenderDefaultTemplate(w http.ResponseWriter, thisView string, data interface{}) {
	renderthis := []string{thisView, "views/layouts/default.gohtml"}
	t, err := template.ParseFiles(renderthis...)
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "default", data)
	if err != nil {
		log.Fatal(err)
	}
}

// getPhoneList reads a single column CSV file as a list of contacts, and
// returns its contents as a map of "int: string"
func getPhoneList(file string) map[int]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	phonelist := make(map[int]string)
	for i := range rows {
		phonelist[i] = rows[i][0]
	}
	return phonelist
}

// getRandomContact picks a random key from a map[int]string and returns its value.
func getRandomContact(phonelist map[int]string) string {
	rand.Seed(time.Now().Unix())
	return phonelist[rand.Intn(len(phonelist))]
}

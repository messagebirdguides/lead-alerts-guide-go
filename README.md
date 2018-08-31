# Instant Lead Alerts
### ⏱ 15 min build time

## Why build instant lead alerts for Sales? 

Even though a lot of business transactions happen on the web, from both a business and user perspective, it's still often preferred to switch the channel and talk on the phone. Especially when it comes to high-value transactions in industries such as real estate or mobility, personal contact is essential.

One way to streamline this workflow is by building callback forms onto your website. Through these forms, customers can enter their contact details and receive a call to their phone, thus skipping queues where prospective leads need to stay on hold. 

Callback requests reflect a high level of purchase intent and should be dealt with as soon as possible to increase the chance of converting a lead. Therefore it's paramount to get them pushed to a sales agent as quickly as possible. SMS messaging has proven to be one of the most instant and effective channels for this use case.

In this MessageBird Developer Guide, we'll show you how to implement a callback form on a Go web application with SMS integration powered by MessageBird for our fictitious car dealership, M.B. Cars.

## Getting Started

First things first, you can download or clone the complete source code from the [MessageBird Developer Guides GitHub repository](https://github.com/messagebirdguides/lead-alerts-guide-go) to run the application on your computer and follow along with the guide. Our sample application is built in Go, so in order to run the sample, you will need to have [Go installed](https://golang.org) on your machine.

We'll also need the [MessageBird REST API package for Go](https://github.com/messagebird/go-rest-api). After you've installed Go, run the following command to install the MessageBird REST API package for Go:

````bash
go get -u github.com/messagebird/go-rest-api
````

### Dependencies

To keep things straightforward and get you up and running as fast and easy as possible, we'll use only packages from the Go standard library and the MessageBird REST API package for Go.

### Create your API Key 🔑

To enable the MessageBird SDK, we need to provide an access key for the API. MessageBird provides keys in _live_ and _test_ modes. To get this application running, we will need to create and use a live API access key. Read more about the difference between test and live API keys [here](https://support.messagebird.com/hc/en-us/articles/360000670709-What-is-the-difference-between-a-live-key-and-a-test-key-).

Let's create your live API access key. First, go to the [MessageBird Dashboard](https://dashboard.messagebird.com/en/user/index); if you have already created an API key it will be shown right there. If you do not see any key on the dashboard or if you're unsure whether this key is in _live_ mode, go to the _Developers_ section and open the [API access (REST) tab](https://dashboard.messagebird.com/en/developers/access). Here, you can create new API keys and manage your existing ones.

If you are having any issues creating your API key, please reach out to our Customer Support team at support@messagebird.com.

**Pro-tip:** To keep our demonstration code simple, we will be saving our API key in `main.go`. However, hardcoding your credentials in the code is a risky practice that should never be used in production applications. A better method, also recommended by the [Twelve-Factor App Definition](https://12factor.net/), is to use environment variables. You can use open source packages such as [GoDotEnv](https://github.com/joho/godotenv) to read your API key from a `.env` file into your Go application. Your `.env` file should be written as follows:

`````env
MESSAGEBIRD_API_KEY=YOUR-API-KEY
`````

To use [GoDotEnv](https://github.com/joho/godotenv) in your application, you can install it by running:

````bash
go get -u github.com/joho/godotenv
````

Then, you'll need to import it in your application:

````go
import (
  // Other imported packages
  "os"

  "github.com/joho/godotenv"
)

func main(){
  // GoDotEnv loads any ".env" file located in the same directory as main.go
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  // Store the value for the key "MESSAGEBIRD_API_KEY" in the loaded '.env' file.
  apikey := os.Getenv("MESSAGEBIRD_API_KEY")

  // The rest of your application ...
}
````

### Initialize the MessageBird Client

Let's now install the [MessageBird's REST API package for Go](https://github.com/messagebird/go-rest-api) by running:

````go
go get -u github.com/messagebird/go-rest-api
````

In your project folder which we created earlier, let's create a `main.go` file, and write the following code:

````go
package main

import (
  "os"

  "github.com/messagebird/go-rest-api"
)

var client *messagebird.Client

func main(){
  client := messagebird.New("<enter-your-api-key>")
}
````

### Routes

To make the structure of our web application clear, we'll first stub the routes that we need. We need two routes:

- A route to ask for the user's name and phone number, so that our agent can contact them.
- A route to display a success state.

To write these routes, we need to modify `main.go` to look like the following:

````go
package main

import (
  "fmt"
  "net/http"

  "github.com/messagebird/go-rest-api"
)

func landing(w http.ResponseWriter, r *http.Request){
  // Do nothing
  fmt.Fprintln("We ask for the user's phone number here.")
}

func callme(w http.ResponseWriter, r *http.Request){
  // Do nothing
  fmt.Fprintln("We display a success state here.")
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
````

Here, we've defined a HTTP handler for each route, and assigned each handler to a corresponding URL path. At the bottom of our `main()` function, we start a HTTP server on port 8080.

**Pro-tip:** We use `http.ListenAndServe` to start our HTTP server here, but you may want to use `http.ListenAndServeTLS` to secure your application's communication over the web with TLS/SSL.

### Views

Now that we've got our routes set up, we can start writing templates to render for each route.
Let's create the template files for rendering views in your project folder:

- `views/layouts/default.gohtml`: The base HTML template we will use with all our views.
- `views/start.gohtml`: Our landing page that takes user input.
- `views/callme.gohtml`: Displays a success state.

The `.gohtml` files will contain our Go HTML template code for each of our views, and our "default" base template.

First, we'll write the base template. In `default.gohtml`, write the following code:

````html
{{ define "default" }}
<!DOCTYPE html>
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>MessageBird Verify Example</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="" type="text/css"/>
  </head>
  <body>
    <main>
    <h1>M.B. Cars 🚙</h1>
    <h2>Largest Car Dealership in Town!</h2>
    <h3>We buy &amp; sell all international brands.</h3>
    {{ template "yield" . }}
    <p><small>A sample application brought to you by <a href="https://developers.messagebird.com/">MessageBird</a> :)</small></p>
    </main>
  </body>
</html>
{{ end }}
````

Notice that in `default.gohtml`, we are calling `{{ template "yield" . }}`. This tells the Go templates package that the contents of any template named "yield" should be rendered in its place when the template is executed.

## Writing a RenderDefaultTemplate() function

Let's write a helper function to help us render our views. At this point, we know that the handlers for our views must:

- Parse the `default.gohtml` template file and one other template file.
- Execute the template.
- Handle errors from the `ParseFiles()` and `ExecuteTemplate()` methods.
- Display an error message if we have one.

We can write a `RenderDefaultTemplate()` helper function to make these calls, keeping our code DRY:

````go
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
````

Then, we can go ahead and rewrite all our handlers in `main.go`:

````go
func landing(w http.ResponseWriter, r *http.Request){
  RenderDefaultTemplate(w,"views/landing.gohtml", nil)
}

func callme(w http.ResponseWriter, r *http.Request){
  RenderDefaultTemplate(w,"views/callme.gohtml", nil)
}
````

Finally, remember to add the `template/html` package to your `import` statement:

````go
import (
  "fmt"
  "net/http"
  "html/template"

  "github.com/messagebird/go-rest-api"
)
````

## Showing a Landing Page

The landing page is a simple HTML page with information about our company, a call to action and a form with two input fields, name and number, and a submit button. All this page does is to take the name and phone number entered by the user and submits a "POST" request to the `/callme` route.

In `landing.gohtml`, let's write the following code:

````html
{{ define "yield" }}
<h4>Leave your number here all talk to a sales agent within 5 minutes, guaranteed!</h4>
<div>
  {{ if .Message }}
  <p><strong>{{ .Message }}</strong></p>
  {{ end }}
  <form action="/callme" method="post">
      <div>
          <label>Your name:</label>
          <br />
          <input type="text" name="name" 
          value="{{ if .Name }}{{ .Name }}{{ end }}" placeholder="Birdie" required/>
      </div>
      <div>
          <label>Your phone number:</label>
          <br />
          <input type="tel" name="phone" 
          value="{{ if .Phone }}{{ .Phone }}{{ end }}" placeholder="+319876543210" required/>
      </div>
      <div>
          <button type="submit">Call me</button>
      </div>        
  </form>
</div>
{{ end }}
````

## Handling Callback Requests

When the user submits the form, the `/callme` route receives their name and number. We then randomly assign an agent to call the user to follow up. We'll assume that we have a list of agent's phone numbers in a `contacts.csv` file, and that we'll pick a phone number from this list at random to send an SMS message to ask the agent who owns that number to call the user.

Our `contacts.csv` file should look like this:

````csv
+319876543210,
+319876543211,
+319876543212,
+319876543213,
+319876543214,
+319876543215,
+319876543216,
+319876543217,
+319876543218,
+319876543219,
+319876543220,
+319876543221,
+319876543222,
+319876543223,
+319876543224,
+319876543225,
+319876543226,
+319876543227,
+319876543228,
+319876543229,
````

**Note**: In order to send a message to a phone number, that phone number must have been added to your MessageBird contact list. See the [MessageBird API Reference](https://developers.messagebird.com/docs/contacts) for more information.

To load the `contacts.csv` file, add the following code under `RenderDefaultTemplate()` in `main.go`:

````go
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
````

And then add this line of code to the top of your `callme()` handler function:

````go
agentToCall := getRandomContact(getPhoneList("contacts.csv"))
````

Above, we've written two helper functions that work together to:

1. Load data from a given CSV file and returns it as a map. We'll use this to load our `contacts.csv` file.
2. Selects a random item from a map using its `key`, and return its value. We'll use this to select a random agent's contact number.

Now we're ready to write the rest of our `callme()` handler function. Modify `callme()` so that it looks like this:

````go
func callme(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    msgBody := "You have a new lead: " + r.FormValue("name") + ". Call them at " + r.FormValue("phone")
    agentToCall := getRandomContact(getPhoneList("contacts.csv"))
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
    // Logging msg for development. Can safely discard in production builds.
    log.Println(msg)
    RenderDefaultTemplate(w, "views/sent.gohtml", nil)
}
````

We'll also need to define a `Data` struct type, which we use to pass an error message and other data to our `landing.gohtml` template:

````go
// Data is for passing data to our templates
type Data struct {
    Name    string
    Phone   string
    Message string
}
````

In our `callme()` handler function, we're:

1. Parsing the form data submitted to the current route.
2. Writing the message to send to the assigned agent, and assigning it to `msgBody`.
3. Picking an agent to call at random, and assigning it to `agentToCall`.
4. Sending an SMS message to the selected agent by calling `sms.Create()`.
5. If `sms.Create()` returns an error, send the user back to the landing page and display and error.
6. If we successfully send the message, display a success state.

To send an SMS message to an agent, we make an `sms.Create()` call that takes four parameters:

- `originator`: The sender ID. This is the name that appears as the sender of the SMS message.
- `recipients`: The API supports an array of recipients; we're sending to only one but this parameter is still specified as a `[]string`. We're passing in our selected agent's contact number here by writing it into a `[]string` literal: `[]string{agentToCall}`.
- `body`: The text of the message that includes the input from the form.
- `params`: Other message parameters that you can specify. Here, we're leaving this as `nil`. See the [MessageBird API Reference](https://developers.messagebird.com/docs) for more information.

## Testing the Application

You're done! Make sure that you've replaced the API key in the sample code with your own, and you're ready to assign leads to agents.

To test your application, run the following command from your console:

````bash
go run main.go
````

Go to http://localhost:8080/ to see the form and request a lead!

## Nice work!

You've just built your own instant lead alerts application with MessageBird! 

You can now use the flow, code snippets and UI examples from this tutorial as an inspiration to build your own SMS-based lead alerts application. Don't forget to download the code from the [MessageBird Developer Guides GitHub repository](https://github.com/messagebirdguides/lead-alerts-guide-go).

## Next steps

Want to build something similar but not quite sure how to get started? Please feel free to let us know at support@messagebird.com, we'd love to help!

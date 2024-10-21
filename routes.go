package main

import (
	"fmt"
	"imperatorapp/models"
	"net/http"
	"strconv"

	"github.com/arc41t3ct/imperator/mailer"
	chi "github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes using aliases
	a.use(a.Middlware.Remember)
	// routes go here using the aloases
	a.get("/", a.Handlers.Home)
	a.get("/go-page", a.Handlers.GoPage)
	a.get("/jet-page", a.Handlers.JetPage)
	a.get("/sessions", a.Handlers.Sessions)

	a.get("/users/login", a.Handlers.Login)
	a.post("/users/login", a.Handlers.LoginPost)
	a.get("/users/logout", a.Handlers.Logout)
	a.get("/users/reset", a.Handlers.PasswordReset)
	a.get("/users/forgot-password", a.Handlers.PasswordForgot)
	a.post("/users/forgot-password", a.Handlers.PasswordForgotPost)
	a.get("/users/reset-password", a.Handlers.PasswordReset)
	a.post("/users/reset-password", a.Handlers.PasswordResetPost)

	a.get("/form", a.Handlers.Form)
	a.post("/form", a.Handlers.FormPost)

	a.get("/logo-download", a.Handlers.DownloadFile)
	a.get("/get-json.json", a.Handlers.GetJSON)
	a.get("/get-xml.xml", a.Handlers.GetXML)
	a.get("/test-encryption", a.Handlers.TestEncryption)

	a.get("/cache", a.Handlers.Cache)
	a.post("/api/cache/set", a.Handlers.CacheSet)
	a.post("/api/cache/get", a.Handlers.CacheGet)
	a.post("/api/cache/del", a.Handlers.CacheDelete)
	a.post("/api/cache/empty", a.Handlers.CacheEmpty)

	a.get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			From:        "someone@gutschein.promo",
			To:          "floridait@gmail.com",
			Subject:     "Test E-Mail from Imperator",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}
		// here we send email using channels with a size of 20 for now
		a.App.Mail.Jobs <- msg
		res := <-a.App.Mail.Results
		if res.Error != nil {
			a.App.ErrorLog.Println(res.Error)
		}
		// here we send a single email using the SendSMTPMessage
		// if err := a.App.Mail.SendSMTPMessage(msg); err != nil {
		// 	a.App.ErrorLog.Println("faileld to send email with error:", err)
		// 	return
		// }
		fmt.Fprint(w, "sent the mail")
	})
	a.App.Routes.Get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := models.User{
			FirstName: "Andre",
			LastName:  "Honsberg",
			Email:     "floridait@gmail.com",
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.Insert(u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d: %s", id, u.FirstName)
	})
	a.App.Routes.Get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		for _, u := range users {
			fmt.Fprintf(w, u.LastName)
		}
	})
	a.App.Routes.Get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		fmt.Fprintf(w, "%s %s %s", u.FirstName, u.LastName, u.Email)
	})
	a.App.Routes.Get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		u.LastName = a.App.CreateRadomString(10)

		u.LastName = ""
		validator := a.App.GetValidator()
		u.Validate(validator)

		if !validator.Valid() {
			for f, e := range validator.GetErrors() {
				fmt.Fprintf(w, "%s validation error: %s", f, e)
			}
			return
		}

		err = u.Update(*u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "updated %s %s %s", u.FirstName, u.LastName, u.Email)
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}

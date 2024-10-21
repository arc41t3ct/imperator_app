package handlers

import (
	"fmt"
	"net/http"

	jet "github.com/CloudyKit/jet/v6"
)

// These handlers below are for testing that is why they are not in the different folders

func (h *Handlers) GoPage(w http.ResponseWriter, r *http.Request) {
	if err := h.renderGo(w, r, "go-template", nil, nil); err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) JetPage(w http.ResponseWriter, r *http.Request) {
	if err := h.renderJet(w, r, "jet-template", nil, nil); err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) Sessions(w http.ResponseWriter, r *http.Request) {
	data := "bar"
	h.sessionPut(r.Context(), "foo", data)
	value := h.sessionGet(r.Context(), "foo")
	variables := make(jet.VarMap)
	variables.Set("foo", value)

	if err := h.renderJet(w, r, "sessions", variables, nil); err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) GetJSON(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		ID      int64    `json:"id"`
		Name    string   `json:"name"`
		Hobbies []string `json:"hobbies"`
	}
	data := payload{}
	data.ID = 10
	data.Name = "Jack"
	data.Hobbies = []string{"programming", "swimming", "smoking weed"}

	if err := h.renderJSON(w, data, http.StatusOK); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) GetXML(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		ID      int64    `xml:"id"`
		Name    string   `xml:"name"`
		Hobbies []string `xml:"hobbies>hoppy"`
	}
	data := payload{}
	data.ID = 10
	data.Name = "Jack"
	data.Hobbies = []string{"programming", "swimming", "smoking weed"}

	if err := h.renderXML(w, data, http.StatusOK); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	if err := h.download(w, r, "./public/images", "logo.jpg"); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) TestEncryption(w http.ResponseWriter, r *http.Request) {
	plainText := "Hello Imperitor"
	fmt.Fprint(w, "Unencrypted: "+plainText+"\n")
	encryptedText, err := h.encrypt(plainText)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Render.Error500(w, r)
		return
	}
	fmt.Fprint(w, "Encrypted: "+encryptedText+"\n")

	decryptedText, err := h.decrypt(encryptedText)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Render.Error500(w, r)
		return
	}
	fmt.Fprint(w, "Decrypted: "+decryptedText+"\n")
}

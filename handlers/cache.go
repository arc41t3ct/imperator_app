package handlers

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// Cache handles request for
func (h *Handlers) Cache(w http.ResponseWriter, r *http.Request) {
	h.App.InfoLog.Println("running handler: Cache")
	if err := h.render(w, r, "cache", nil, nil); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) CacheSet(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		CSRF  string `json:"csrf_token"`
	}
	if err := h.renderJSON(w, &userInput, http.StatusOK); err != nil {
		h.App.Render.Error500(w, r)
		return
	}

	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.App.Render.Error500(w, r)
		return
	}

	if err := h.App.Cache.Set(userInput.Name, userInput.Value); err != nil {
		h.App.Render.Error500(w, r)
		return
	}
	var response struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	response.Error = false
	response.Message = fmt.Sprintf("%s was set with %s", userInput.Name, userInput.Value)
	if err := h.App.Render.WriteJSON(w, response, http.StatusCreated); err != nil {
		h.App.Render.Error500(w, r)
		return
	}
}

func (h *Handlers) CacheGet(w http.ResponseWriter, r *http.Request) {
	var errMsg string
	var inCache = true
	var status = http.StatusOK
	var userInput struct {
		Name string `json:"name"`
		CSRF string `json:"csrf_token"`
	}
	if err := h.renderJSON(w, &userInput, http.StatusOK); err != nil {
		h.App.Render.Error500(w, r)
		return
	}

	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.App.Render.Error500(w, r)
		return
	}

	data, err := h.App.Cache.Get(userInput.Name)
	if err != nil {
		errMsg = "not found in cache"
		inCache = false
	}
	var response struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Value   string `json:"value"`
	}

	if inCache {
		response.Error = false
		response.Message = "success"
		response.Value = data.(string)
	} else {
		response.Error = true
		response.Message = errMsg
		status = http.StatusNotFound
	}

	if err := h.renderJSON(w, response, status); err != nil {
		h.App.Render.Error500(w, r)
		return
	}
}

func (h *Handlers) CacheDelete(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Name string `json:"name"`
		CSRF string `json:"csrf_token"`
	}
	if err := h.renderJSON(w, &userInput, http.StatusOK); err != nil {
		h.App.Render.Error500(w, r)
		return
	}

	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.App.Render.Error500(w, r)
		return
	}

	if err := h.App.Cache.Forget(userInput.Name); err != nil {
		h.App.Render.Error500(w, r)
		return
	}
	var response struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	response.Error = false
	response.Message = "successfully removed from cache if existed"

	if err := h.renderJSON(w, response, http.StatusOK); err != nil {
		h.App.Render.Error500(w, r)
		return
	}
}

func (h *Handlers) CacheEmpty(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		CSRF string `json:"csrf_token"`
	}
	if err := h.renderJSON(w, &userInput, http.StatusOK); err != nil {
		h.App.Render.Error500(w, r)
		return
	}

	if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
		h.App.Render.Error500(w, r)
		return
	}

	if err := h.App.Cache.Empty(); err != nil {
		h.App.Render.Error500(w, r)
		return
	}
	var response struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	response.Error = false
	response.Message = "successfully emptied cache"

	if err := h.renderJSON(w, response, http.StatusOK); err != nil {
		h.App.Render.Error500(w, r)
		return
	}
}

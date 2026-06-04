package handlers

import (
	"html/template"
	"net/http"

	"desafio4/services"
)

type Handlers struct {
	Service *services.AuthService
}

func NewHandlers(s *services.AuthService) *Handlers {
	return &Handlers{Service: s}
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("X-Session-Token")
	if token != "" {
		return token
	}
	return r.FormValue("token")
}

func (h *Handlers) RegisterGET(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func (h *Handlers) RegisterPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	_, err := h.Service.Register(username, password)
	if err != nil {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, map[string]string{"Error": err.Error()})
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handlers) LoginGET(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func (h *Handlers) LoginPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := h.Service.Login(username, password)
	if err != nil {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, map[string]string{"Error": "Credenciais inválidas"})
		return
	}

	http.Redirect(w, r, "/user/profile?token="+token, http.StatusSeeOther)
}


func (h *Handlers) Profile(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	user, err := h.Service.GetUserByToken(token)
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/profile.html"))
	
	data := map[string]interface{}{
		"User":  user,
		"Token": token,
	}
	tmpl.Execute(w, data)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token != "" {
		h.Service.Logout(token)
	}
	
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
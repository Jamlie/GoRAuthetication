package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"path"

	"golang.org/x/crypto/bcrypt"

	"server/internal/utils"
)

type SignupCredentials struct {
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

func (s SignupCredentials) IsEmpty() bool {
	if s.Email == "" || s.Password == "" || s.Name == "" {
		return true
	}

	return false
}

type LoginCredentials struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (s LoginCredentials) IsEmpty() bool {
	if s.Email == "" || s.Password == "" {
		return true
	}

	return false
}

func (s *Server) RegisterRoutes() http.Handler {
	s.app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("dist", "index.html"))
	})

	s.app.HandleFunc("/first", s.HandleFirstPage)

	s.app.HandleFunc("/second", s.HandleSecondPage)

	s.Use(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("dist/assets"))),
	)

	s.app.HandleFunc("/api/signup", s.HandleSignup)

	s.app.HandleFunc("/api/login", s.HandleLogin)

	s.app.HandleFunc("/does-session-exist", s.HandleDoesSessionExist)
	s.app.HandleFunc("/add-session", s.HandleAddSession())
	s.app.HandleFunc("/get-session", s.HandleGetSession)

	return s.app
}

func (s *Server) HandleFirstPage(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "Authenticated")
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := session.Values["isAdded"]; !ok {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	http.ServeFile(w, r, path.Join("dist", "index.html"))
}

func (s *Server) HandleSecondPage(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "Authenticated")
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := session.Values["isAdded"]; ok {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	http.ServeFile(w, r, path.Join("dist", "index.html"))
}

func (s *Server) HandleDoesSessionExist(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "Authenticated")
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := session.Values["isAdded"]; ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (s *Server) HandleAddSession(keyAndValue ...utils.Pair[any, any]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, "Authenticated")
		if err != nil {
			log.Fatal(err)
		}

		for _, kv := range keyAndValue {
			session.Values[kv.First] = kv.Second
		}

		if _, ok := session.Values["isAdded"]; !ok {
			session.Values["isAdded"] = true
		}

		session.Options.MaxAge = 7 * 24 * 60 * 60
		err = session.Save(r, w)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "Authenticated")
	if err != nil {
		log.Fatal(err)
	}

	sessionValues := map[string]any{
		"name": session.Values["name"],
	}

	jsonRes, err := json.Marshal(sessionValues)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(jsonRes)
}

func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "Authenticated")
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := session.Values["isAdded"]; ok {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	var user SignupCredentials

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	if user.IsEmpty() {
		m := map[string]string{
			"message": "",
			"error":   "Fields are empty",
		}

		jsonRes, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}

		_, _ = w.Write(jsonRes)
		return
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		m := map[string]string{
			"message": "",
			"error":   "Invalid email address",
		}

		jsonRes, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(jsonRes)
		return
	}

	if err := s.db.AddUser(user.Email, user.Name, user.Password); err != nil {
		m := map[string]string{
			"message": "",
			"error":   "Invalid Credentials",
		}

		jsonRes, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(jsonRes)
		return
	}

	m := map[string]string{
		"message": "Signed Up Successfully",
		"error":   "",
	}
	jsonRes, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonRes)
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "Authenticated")
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := session.Values["isAdded"]; ok {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	var user LoginCredentials

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	if user.IsEmpty() {
		m := map[string]string{
			"message": "",
			"error":   "Fields are empty",
		}

		jsonRes, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}

		_, _ = w.Write(jsonRes)
		return
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		m := map[string]string{
			"message": "",
			"error":   "Invalid email address",
		}

		jsonRes, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(jsonRes)
		return
	}

	userDB, err := s.db.GetUserByEmail(user.Email)
	if err != nil {
		m := map[string]string{
			"message": "",
			"error":   "Email does not exist",
		}

		jsonRes, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(jsonRes)
		return
	}

	if err := bcrypt.CompareHashAndPassword(userDB.Password, []byte(user.Password)); err != nil {
		m := map[string]string{
			"message": "",
			"error":   "Invalid email or password",
		}
		jsonRes, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(jsonRes)
		return
	}

	s.HandleAddSession(utils.NewPair[any, any]("name", userDB.Name))(w, r)

	m := map[string]string{
		"message": "Signed Up Successfully",
		"error":   "",
	}

	jsonRes, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = w.Write(jsonRes)
}

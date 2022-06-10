package site

import (
	"context"
	"log"
	"net/http"
)

func (h *site) indexFunc() http.HandlerFunc {
	type data struct {
		LatestPosts []Post
		CurrentUser *User
	}
	return func(w http.ResponseWriter, r *http.Request) {
		posts := []Post{
			{
				Title:   "Awesome Post!",
				Content: "lololol",
			},
		}
		d := data{LatestPosts: posts, CurrentUser: nil}
		err := h.t.ExecuteTemplate(w, "index.html.tmpl", d)
		if err != nil {
			// todo(javier): Write error message to response.
			// w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Error rendering tempalte. %s", err)
		}
	}
}

func (h *site) handleSignin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := h.t.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
			// todo(javier): log error instead of dying.
			if err != nil {
				log.Fatalf("Could not render sign-in template. %s", err)
			}
		case http.MethodPost:
			svc := NewSessionService(h.db)
			// todo(javier): check credentials
			err := r.ParseForm()
			if err != nil {
				log.Fatalf("Could not parse form. %s", err)
			}
			email := r.PostForm.Get("email")
			password := r.PostForm.Get("password")
			sid, err := svc.Authenticate(ctx, email, password)
			if err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				// TODO(javier): Add an error to the template
				err := h.t.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
				// todo(javier): log error instead of dying.
				if err != nil {
					log.Fatalf("Could not render sign-in template. %s", err)
				}
			}
			cookie := &http.Cookie{Name: "sid", Value: string(sid)}
			http.SetCookie(w, cookie)
			// todo(javier): redirect to / or redirect_to query param
			http.Redirect(w, r, "/", http.StatusSeeOther)
			w.Write([]byte("login successful."))

		}

	}
}

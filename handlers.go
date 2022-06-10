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
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error rendering tempalte. %s", err)
			return
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
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("Error rendering tempalte. %s", err)
				return
			}
		case http.MethodPost:
			// todo(javier): check credentials
			err := r.ParseForm()
			if err != nil {
				log.Fatalf("Could not parse form. %s", err)
			}
			email := r.PostForm.Get("email")
			password := r.PostForm.Get("password")
			sid, err := h.sessionsvc.New(ctx, email, password)
			if err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				// TODO(javier): Add an error to the template
				err := h.t.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("Could not render sign-in template. %s\n", err)
				}
			}
			cookie := &http.Cookie{Name: "sid", Value: string(sid)}
			http.SetCookie(w, cookie)
			url := r.Form.Get("redirect_to")
			if url == "" {
				url = "/"
			}
			http.Redirect(w, r, url, http.StatusSeeOther)
			w.Write([]byte("login successful."))

		}

	}
}

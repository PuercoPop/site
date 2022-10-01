package site

// func (h *site) NewPost() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.TODO()
// 		switch r.Method {
// 		case "GET":
// 			err := h.t.ExecuteTemplate(w, "newpost.html.tmpl", nil)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				log.Printf("Error rendering tempalte. %s", err)
// 				return
// 			}
// 		case "POST":
// 			err := r.ParseForm()
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				log.Printf("Error parsing form. %s", err)
// 				return
// 			}
// 			title := r.Form.Get("title")
// 			content := r.Form.Get("content")
// 			err = h.svc.CreatePost(ctx, title, content)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				log.Printf("Error creating post. %s", err)
// 				return
// 			}
// 			http.Redirect(w, r, "/", http.StatusFound)
// 		}
// 	}
// }

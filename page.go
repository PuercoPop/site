package swiki

// PageMux takes care of creating, updating, and deleting pages. It is mounted
// under /p/.
//
// - POST /p/ - Create a Page
// - PATCH /p/:slug
// - DELETE /p/:slug

// func NewPageMux() *http.ServeMux {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/")

// }

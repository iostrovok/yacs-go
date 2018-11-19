package httpserver

/*
	It a HTTP server for simple management of browsing, starting and canceling of jobs.
*/

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type contextKey string

// We use settings if we want to pass parameters into http handler.
type settings struct {
	Dir      string
	Ping     string
	Username string
	Password string
}

func getTimeHelper(key string, r *http.Request) (time.Time, error) {

	str := r.FormValue(key)
	if str == "" {
		return time.Time{}, fmt.Errorf("not found %s", key)
	}

	return time.Parse("2006-01-02", str)
}

func jsonWrite(w http.ResponseWriter, what interface{}) {

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(what)
	if err == nil {
		w.Write(b)
	} else {
		w.Write([]byte(`{"error":"internal error"}`))
	}
}

// It gives us a possibility to pass parameters into http handler.
func wrapHandler(h http.HandlerFunc, s *settings) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(s.Username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(s.Password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Duck out"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		// Get new context with key-value "settings"
		ctx := context.WithValue(r.Context(), contextKey("settings"), s)

		// Get new http.Request with the new context
		r = r.WithContext(ctx)

		// Call your original http.Handler
		h.ServeHTTP(w, r)
	})
}

func getContextHelper(r *http.Request) (*settings, error) {
	ctx := r.Context()
	s := ctx.Value(contextKey("settings"))

	var sets *settings

	switch s.(type) {
	case *settings:
		sets = s.(*settings)
	default:
		return nil, fmt.Errorf("Wrong context")
	}

	return sets, nil
}

func handlerState(w http.ResponseWriter, r *http.Request) {

	sets, err := getContextHelper(r)
	if err != nil {
		jsonWrite(w, map[string]string{"error": fmt.Sprint(err)})
		return
	}

	fmt.Printf("sets: %s\n", sets)

	res := map[string]string{
		"error": "",
		"id":    "jobID",
	}
	jsonWrite(w, res)
}

func handleFileServer(dir, prefix string) http.HandlerFunc {

	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix(prefix, fs).ServeHTTP

	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL)
		realHandler(w, req)
	}
}

// Run is main function. It starts the HTTP server
func Run(htmlDir string) {

	fmt.Println("SERVER.Run", "HTTP server starting...")

	s := &settings{
		Dir:      "string DIR",
		Ping:     "string Ping",
		Username: "Username",
		Password: "Password",
	}

	http.HandleFunc("/static/", wrapHandler(handleFileServer(htmlDir, "/static/"), s))
	http.HandleFunc("/state/", wrapHandler(handlerState, s))
	fmt.Println(http.ListenAndServe(":8080", nil))
}

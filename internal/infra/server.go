package infra

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]Route
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]Route),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc, method string) {
	var key = fmt.Sprintf("%s-%s", method, path)
	s.Handlers[key] = Route{
		Path:    path,
		Handler: handler,
		Method:  method,
	}
}

func (s *WebServer) Use(middleware func(next http.Handler) http.Handler) {
	s.Router.Use(middleware)
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	s.Use(middleware.Recoverer)
	for _, handler := range s.Handlers {
		s.Router.MethodFunc(handler.Method, handler.Path, handler.Handler)
		fmt.Println("Registering handler", handler.Path, handler.Method)
	}
	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}

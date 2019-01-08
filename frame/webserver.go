package frame

import (
	"net/http"
	"strconv"

	"nvmtech.nl/homemanager/middleware"

	"github.com/gorilla/mux"
	"nvmtech.nl/homemanager/tools/log"
)

// WebServer is an object responsible for
type WebServer struct {

	// Internal variables
	router    *mux.Router
	endpoints []*Endpoint
}

// NewWebServer creates a webserver struct with the provided
// attributes, providing a constructor
func NewWebServer() *WebServer {
	wb := WebServer{
		mux.NewRouter(),
		[]*Endpoint{},
	}

	return &wb
}

// RegisterEndpoint adds an endpoint to the router
func (ws *WebServer) RegisterEndpoint(name string, function func(http.ResponseWriter, *http.Request)) {
	// Add to the map
	endp := NewEndpoint(name, function)
	ws.AddEndpoint(endp)
}

// AddEndpoint adds the provided enpoint to the list of known endpoints
func (ws *WebServer) AddEndpoint(endp *Endpoint) {
	ws.endpoints = append(ws.endpoints, endp)
}

// AddEndpoints adds the provided enpoint to the list of known endpoints
func (ws *WebServer) AddEndpoints(endp []*Endpoint) {
	ws.endpoints = append(ws.endpoints, endp...)
}

// Run starts the webserver on the provided port
func (ws *WebServer) Run(port int) {
	// Parse endpoints and register them to the router
	if len(ws.endpoints) == 0 {
		log.Warn("WebServer", "No endpoints found for server, i'll be useless")
	}

	for _, endp := range ws.endpoints {
		log.Info("WebServer", "Registered Endpoint: "+endp.URL)
		ws.router.HandleFunc(
			endp.URL, endp.Function,
		)
	}

	// TODO Make this dynamic, the static path and middlware
	// Add static file handler
	ws.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	ws.router.Use(middleware.LogHTTP)

	// Run server
	log.Info("WebServer", "Running Webserver at port "+strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), ws.router)
	log.Fatal("WebServer", err.Error())
}

// Endpoint represents an endpoint for the webapp
type Endpoint struct {
	URL      string
	Function func(http.ResponseWriter, *http.Request)
}

// NewEndpoint is a constructor for the endoints
func NewEndpoint(URL string, function func(http.ResponseWriter, *http.Request)) *Endpoint {
	endp := Endpoint{
		URL,
		function,
	}

	return &endp
}

package service

import (
	"net/http"
	"strings"
)

// route stores information needed to make a request and build URLs
type route struct {
	method  string                                   // method allowed in http/https request
	prefix  string                                   // prefix of URL Path
	handler func(http.ResponseWriter, *http.Request) // request handler for route
	routes  []*route                                 // slice of all routes with the same root route
}

/*
Method : Sets the method that will be allowed in the given request.

This function is to be used with in tandem with the CreateEndpoint function.
Without this function the endpoints created will be able to use all methods by default.

IMPORTANT!! Order matters. Method will not take effect on any CreateEndpoint methods that
appear before it.
*/
func (route *route) Method(methodType string) *route {
	route.method = strings.ToUpper(methodType)
	return route
}

/*
Prefix : Sets the prefix for the Path.

This function is to be used when you want a set of endpoints to have the same base URL Path.

IMPORTANT!! Order matters. Prefix will not take effect on any CreateEndpoint methods that
appear before it.
*/
func (route *route) Prefix(prefix string) *route {
	route.prefix = prefix
	return route
}

// NewRoute : returns an empty route Object
func NewRoute() *route {
	return &route{}
}

/*
verifyMethod : Checked if the method in the http/https matches the allowed methods then called the handler

By default, all methods are allowed.
*/
func (route *route) verifyMethod(writer http.ResponseWriter, r *http.Request) {
	if route.method != r.Method && route.method != "" {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		SendStatus(r, http.StatusMethodNotAllowed)
		return
	}
	route.handler(writer, r)
	SendStatus(r, http.StatusOK)
}

/*
handleRoute : adds handler to the route and calls http.HandleFunc.

http.HandleFunc is called with `route.prefix + endpoint` as the pattern and verifyMethod as the handler
*/
func (route *route) handleRoute(endpoint string, handler func(http.ResponseWriter, *http.Request)) *route {
	route.handler = handler
	http.HandleFunc(
		route.prefix+endpoint,
		route.verifyMethod,
	)
	return route
}

/*
CreateEndpoint : Creates a new route and appends it to the root route's slice of routes
*/
func (route *route) CreateEndpoint(endpoint string, handler func(http.ResponseWriter, *http.Request)) *route {
	route.routes = append(route.routes, NewRoute().Method(route.method).Prefix(route.prefix).handleRoute(endpoint, handler))
	return route
}

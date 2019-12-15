package types

import (
	"net/http"
)

//Proxy works as a proxy server between
type Proxy interface {

	//ServerHTTP servers http request
	ServeHTTP(http.ResponseWriter, *http.Request)
}
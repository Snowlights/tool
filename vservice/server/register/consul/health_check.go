package consul

import (
	"net/http"
)

const DefaultCheckPath = "/health"

type Checker struct{}

func (Checker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("heath check success!"))
}

package linksservice

import (
	"fmt"
	"net/http"
)

func (s *Server) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

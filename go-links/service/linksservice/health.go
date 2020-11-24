package linksservice

import (
	"fmt"
	"net/http"
)

func (s *Server) CheckHealth(w http.ResponseWriter, r *http.Request) {
	// check if we can load links at all (our connection to the database is working)
	_, err := s.db.ListRecentLinks(r.Context(), 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "cannot list recent links: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

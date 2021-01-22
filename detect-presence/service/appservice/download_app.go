package appservice

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) DownloadApp(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "iPhone") {
		installURL := fmt.Sprintf("https://%s/app/install_manifest", r.Host)
		http.Redirect(w, r, "itms-services://?action=download-manifest&url="+installURL, http.StatusTemporaryRedirect)
		return
	}

	s.InstallApp(w, r)
}

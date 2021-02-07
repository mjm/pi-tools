package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/etherlabsio/healthcheck"
	vaultapi "github.com/hashicorp/vault/api"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/vault-proxy/service/authservice"
)

var (
	authPath     = flag.String("auth-path", "webauthn", "Path of webauthn auth method in Vault")
	cookieDomain = flag.String("cookie-domain", "homelab", "Domain to use for cookies, must be common base between callback and app hostnames")
	staticDir    = flag.String("static-dir", "/static", "Path to static assets")
)

func main() {
	rpc.SetDefaultHTTPPort(2220)
	flag.Parse()

	stopObs := observability.MustStart("vault-proxy")
	defer stopObs()

	vault, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		log.Panicf("Error creating Vault client: %v", err)
	}

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second))))

	authService, err := authservice.New(vault, authservice.Config{
		AuthPath:     *authPath,
		CookieDomain: *cookieDomain,
		CookieKey:    os.Getenv("COOKIE_KEY"),
	})
	if err != nil {
		log.Panicf("Error creating auth service: %v", err)
	}

	http.Handle("/auth",
		otelhttp.WithRouteTag("HandleAuthRequest", http.HandlerFunc(authService.HandleAuthRequest)))

	http.Handle("/webauthn/registration/start",
		otelhttp.WithRouteTag("StartRegistration", http.HandlerFunc(authService.StartRegistration)))
	http.Handle("/webauthn/registration/finish",
		otelhttp.WithRouteTag("FinishRegistration", http.HandlerFunc(authService.FinishRegistration)))
	http.Handle("/webauthn/login/start",
		otelhttp.WithRouteTag("StartLogin", http.HandlerFunc(authService.StartLogin)))
	http.Handle("/webauthn/login/finish",
		otelhttp.WithRouteTag("FinishLogin", http.HandlerFunc(authService.FinishLogin)))

	http.Handle("/webauthn/register",
		otelhttp.WithRouteTag("Register", serveStaticFile("register.html")))
	http.Handle("/webauthn/login",
		otelhttp.WithRouteTag("Login", serveStaticFile("login.html")))
	http.Handle("/webauthn/login_app",
		otelhttp.WithRouteTag("Login", serveStaticFile("login.html")))

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(*staticDir, r.URL.Path[8:])
		http.ServeFile(w, r, p)
	})

	go rpc.ListenAndServe()
	signal.Wait()
}

func serveStaticFile(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(*staticDir, name)
		http.ServeFile(w, r, p)
	})
}

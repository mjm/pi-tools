package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/etherlabsio/healthcheck"
	vaultapi "github.com/hashicorp/vault/api"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/vault-proxy/service/authservice"
)

var (
	cookieDomain = flag.String("cookie-domain", "homelab", "Domain to use for cookies, must be common base between callback and app hostnames")
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

	oauth := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_PROXY_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_PROXY_CLIENT_SECRET"),
		Scopes:       []string{"read:org"},
		Endpoint:     endpoints.GitHub,
	}

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second))))

	authService, err := authservice.New(vault, oauth, authservice.Config{
		CookieDomain: *cookieDomain,
	})
	if err != nil {
		log.Panicf("Error creating auth service: %v", err)
	}

	http.Handle("/auth",
		otelhttp.WithRouteTag("HandleAuthRequest", http.HandlerFunc(authService.HandleAuthRequest)))
	http.Handle("/oauth/start",
		otelhttp.WithRouteTag("StartOAuth", http.HandlerFunc(authService.StartOAuth)))
	http.Handle("/oauth/callback",
		otelhttp.WithRouteTag("HandleOAuthCallback", http.HandlerFunc(authService.HandleOAuthCallback)))

	go rpc.ListenAndServe()
	signal.Wait()
}

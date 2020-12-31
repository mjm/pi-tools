package main

import (
	"context"
	"flag"
	"log"
	"os"
	"text/template"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	nodeName = flag.String("node-name", "",
		"Name of the current Kubernetes node")
	localZoneFile = flag.String("local-zone-file", "/data/bind/etc/db.homelab.local",
		"Path to write the local zone file")
	tailscaleZoneFile = flag.String("tailscale-zone-file", "/data/bind/etc/db.homelab.tailscale",
		"Path to write the tailscale zone file")
)

type TemplateInput struct {
	Node  Node
	Nodes []Node
}

type Node struct {
	Name        string
	IP          string
	TailscaleIP string
}

const tailscaleIPKey = "tailscale.com/ip"

const (
	templates = `
{{ define "common" -}}
$TTL  1m
@   IN  SOA localhost. matt.mattmoriarity.com. (
                  1
	             1m     ; Refresh
			     1h		; Retry
			     1w		; Expire
			     1h )	; Negative Cache TTL
@   IN  NS  localhost.

*.homelab.  IN  CNAME {{ .Node.Name }}.homelab.
unifi	IN  A 10.0.0.1
nas		IN  A 10.0.0.10
{{ end }}

{{ define "local-zone" -}}
{{ template "common" . }}
{{ range .Nodes -}}
{{ .Name }} IN  A {{ .IP }}
{{ end -}}
{{ end }}

{{ define "tailscale-zone" -}}
{{ template "common" . }}
{{ range .Nodes -}}
{{ .Name }} IN  A {{ .TailscaleIP }}
{{ end -}}
{{ end }}
`
)

func main() {
	flag.Parse()

	tmpls := template.Must(template.New("zone-templates").Parse(templates))

	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset := kubernetes.NewForConfigOrDie(cfg)

	ctx := context.Background()
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var input TemplateInput
	for _, n := range nodes.Items {
		node := Node{
			Name:        n.GetName(),
			TailscaleIP: n.GetAnnotations()[tailscaleIPKey],
		}
		for _, addr := range n.Status.Addresses {
			if addr.Type == v1.NodeInternalIP {
				node.IP = addr.Address
			}
		}

		input.Nodes = append(input.Nodes, node)
		if n.GetName() == *nodeName {
			input.Node = node
		}
	}

	if input.Node.Name == "" {
		log.Fatal("Could not find node in list that matched name %q", *nodeName)
	}

	localFile, err := os.Create(*localZoneFile)
	if err != nil {
		log.Fatal(err)
	}
	defer localFile.Close()

	if err := tmpls.ExecuteTemplate(localFile, "local-zone", input); err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote local zone to %s", *localZoneFile)

	tailscaleFile, err := os.Create(*tailscaleZoneFile)
	if err != nil {
		log.Fatal(err)
	}
	defer tailscaleFile.Close()

	if err := tmpls.ExecuteTemplate(tailscaleFile, "tailscale-zone", input); err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote Tailscale zone to %s", *tailscaleZoneFile)
}

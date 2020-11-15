package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	nodeName = flag.String("node-name", "",
		"Name of the Kubernetes node the controller is running on")
	dnsTemplatePath = flag.String("dns-template-path", "/cfg/dns-template.tpl",
		"The template to use to generate the custom DNS records file")
)

var (
	dnsTemplate     *template.Template
	dnsTemplateName string
)

func main() {
	flag.Parse()

	dnsTemplate = template.Must(template.ParseFiles(*dnsTemplatePath))
	dnsTemplateName = filepath.Base(*dnsTemplatePath)

	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset := kubernetes.NewForConfigOrDie(cfg)
	updateDebouncer := &debouncer{
		after: 3 * time.Second,
	}

	watchlist := cache.NewListWatchFromClient(clientset.NetworkingV1().RESTClient(), "ingresses", "", fields.Everything())
	store, controller := cache.NewInformer(watchlist, &networkingv1.Ingress{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updateDebouncer.tick()
		},
		DeleteFunc: func(obj interface{}) {
			updateDebouncer.tick()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateDebouncer.tick()
		},
	})

	ctx := context.Background()
	updateDebouncer.f = func() {
		updateDNSRecords(ctx, clientset.CoreV1(), store)
	}

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

type TemplateInput struct {
	Node struct {
		Name string
		IP   string
	}
	Ingresses []Ingress
}

type Ingress struct {
	Hostname string
}

func updateDNSRecords(ctx context.Context, v1client v1.CoreV1Interface, store cache.Store) {
	node, err := v1client.Nodes().Get(ctx, *nodeName, metav1.GetOptions{})
	if err != nil {
		log.Panicf("Could not load node details: %v", err)
	}

	records := store.List()
	log.Printf("Updating custom DNS records, found %d ingresses", len(records))

	var input TemplateInput
	input.Node.Name = node.GetName()
	input.Node.IP = node.GetAnnotations()["tailscale.com/ip"]

	for _, record := range records {
		ingress := record.(*networkingv1.Ingress)
		for _, rule := range ingress.Spec.Rules {
			if rule.Host != "" && !strings.Contains(rule.Host, "*.") {
				input.Ingresses = append(input.Ingresses, Ingress{Hostname: rule.Host})
			}
		}
	}

	log.Printf("Produced %d DNS records from ingress rules", len(input.Ingresses))

	f, err := os.OpenFile("/etc/pihole/custom.list", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Error opening custom DNS entries file: %v", err)
		return
	}
	defer f.Close()

	if err := dnsTemplate.ExecuteTemplate(f, dnsTemplateName, input); err != nil {
		log.Printf("Error writing custom DNS entries: %v", err)
		return
	}

	// restart the DNS server
	cmd := exec.Command("/usr/bin/pkill", "pihole-FTL")
	if err := cmd.Run(); err != nil {
		log.Printf("Error restarting DNS server: %v", err)
		return
	}

	log.Printf("Successfully updated DNS records")
}

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

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
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
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updateDebouncer.tick()
		},
		DeleteFunc: func(obj interface{}) {
			updateDebouncer.tick()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateDebouncer.tick()
		},
	}

	watchIngresses := cache.NewListWatchFromClient(clientset.NetworkingV1().RESTClient(), "ingresses", "", fields.Everything())
	ingressStore, ingressController := cache.NewInformer(watchIngresses, &networkingv1.Ingress{}, 0, handlers)

	watchNodes := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "nodes", "", fields.Everything())
	nodeStore, nodeController := cache.NewInformer(watchNodes, &v1.Node{}, 0, handlers)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateDebouncer.f = func() {
		updateDNSRecords(ingressStore, nodeStore)
	}

	go ingressController.Run(ctx.Done())
	go nodeController.Run(ctx.Done())
	select {}
}

type TemplateInput struct {
	Node      Node
	Nodes     []Node
	Ingresses []Ingress
}

type Node struct {
	Name string
	IP   string
}

type Ingress struct {
	Hostname string
}

func updateDNSRecords(ingressStore cache.Store, nodeStore cache.Store) {
	nodeObj, exists, err := nodeStore.GetByKey(*nodeName)
	if err != nil {
		log.Panicf("Could not lookup node in node cache: %v", err)
		return
	}
	if !exists {
		log.Printf("Node not found with key %q. Node keys available: %v", *nodeName, nodeStore.ListKeys())
		return
	}
	node := nodeObj.(*v1.Node)

	ingresses := ingressStore.List()
	nodes := nodeStore.List()
	log.Printf("Updating custom DNS records, found %d ingresses, %d nodes", len(ingresses), len(nodes))

	var input TemplateInput
	input.Node.Name = node.GetName()
	input.Node.IP = node.GetAnnotations()["tailscale.com/ip"]

	for _, record := range nodes {
		n := record.(*v1.Node)
		input.Nodes = append(input.Nodes, Node{
			Name: n.GetName(),
			IP:   n.GetAnnotations()["tailscale.com/ip"],
		})
	}

	for _, record := range ingresses {
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

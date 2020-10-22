package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	ip             = flag.String("ip", "", "IP address to use for generated A records")
	extraHostnames = flag.String("extra-hostnames", "", "Comma-separated list of additional hostnames to include in the DNS records")
)

func main() {
	flag.Parse()

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

	updateDebouncer.f = func() {
		updateDNSRecords(store)
	}

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func updateDNSRecords(store cache.Store) {
	records := store.List()
	log.Printf("Updating custom DNS records, found %d ingresses", len(records))

	var dnsEntries []string

	if len(*extraHostnames) > 0 {
		for _, host := range strings.Split(*extraHostnames, ",") {
			dnsEntries = append(dnsEntries, fmt.Sprintf("%s %s", *ip, host))
		}
	}

	for _, record := range records {
		ingress := record.(*networkingv1.Ingress)
		for _, rule := range ingress.Spec.Rules {
			if rule.Host != "" && !strings.Contains(rule.Host, "*.") {
				dnsEntries = append(dnsEntries, fmt.Sprintf("%s %s", *ip, rule.Host))
			}
		}
	}

	log.Printf("Produced %d DNS records", len(dnsEntries))

	allEntriesStr := strings.Join(dnsEntries, "\n") + "\n"
	if err := ioutil.WriteFile("/etc/pihole/custom.list", []byte(allEntriesStr), 0666); err != nil {
		log.Printf("Error writing new DNS entries: %v", err)
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

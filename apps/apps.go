package apps

import (
	"github.com/mjm/pi-tools/apps/adminer"
	"github.com/mjm/pi-tools/apps/backup"
	"github.com/mjm/pi-tools/apps/blackboxexporter"
	"github.com/mjm/pi-tools/apps/blocky"
	"github.com/mjm/pi-tools/apps/consulexporter"
	"github.com/mjm/pi-tools/apps/deploy"
	"github.com/mjm/pi-tools/apps/grafana"
	"github.com/mjm/pi-tools/apps/ingress"
	"github.com/mjm/pi-tools/apps/loki"
	"github.com/mjm/pi-tools/apps/nodeexporter"
	"github.com/mjm/pi-tools/apps/nut"
	"github.com/mjm/pi-tools/apps/otel"
	"github.com/mjm/pi-tools/apps/presence"
	"github.com/mjm/pi-tools/apps/promtail"
	"github.com/mjm/pi-tools/apps/pushgateway"
	"github.com/mjm/pi-tools/apps/unifiexporter"
	"github.com/mjm/pi-tools/apps/vaultproxy"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

func Load() {
	nomadic.Register(adminer.New("adminer"))
	nomadic.Register(backup.New("backup"))
	nomadic.Register(blackboxexporter.New("blackbox-exporter"))
	nomadic.Register(blocky.New("blocky"))
	nomadic.Register(consulexporter.New("consul-exporter"))
	nomadic.Register(deploy.New("deploy"))
	nomadic.Register(grafana.New("grafana"))
	nomadic.Register(ingress.New("ingress"))
	nomadic.Register(loki.New("loki"))
	nomadic.Register(nodeexporter.New("node-exporter"))
	nomadic.Register(nut.New("nut"))
	nomadic.Register(otel.New("otel"))
	nomadic.Register(presence.New("presence", "beacon"))
	nomadic.Register(promtail.New("promtail"))
	nomadic.Register(pushgateway.New("pushgateway"))
	nomadic.Register(unifiexporter.New("unifi-exporter"))
	nomadic.Register(vaultproxy.New("vault-proxy"))
}

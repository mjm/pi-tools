package apps

import (
	"github.com/mjm/pi-tools/apps/blocky"
	"github.com/mjm/pi-tools/apps/ingress"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

func Load() {
	nomadic.Register(blocky.New("blocky"))
	nomadic.Register(ingress.New("ingress"))
}

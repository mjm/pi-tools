package apps

import (
	"github.com/mjm/pi-tools/apps/ingress"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

func Load() {
	nomadic.Register(ingress.New("ingress"))
}

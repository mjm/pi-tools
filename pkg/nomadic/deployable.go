package nomadic

import (
	"context"
)

type Deployable interface {
	Name() string
	Install(ctx context.Context, clients Clients) error
	Uninstall(ctx context.Context, clients Clients) error
}

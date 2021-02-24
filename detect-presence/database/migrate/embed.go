package migrate

import (
	"embed"
)

//go:embed *.sql
var FS embed.FS

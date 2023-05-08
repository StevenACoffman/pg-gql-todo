package assets

import (
	"embed"
)

//go:embed "queries" "migrations"
var EmbeddedFiles embed.FS

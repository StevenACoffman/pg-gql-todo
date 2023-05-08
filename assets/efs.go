package assets

import (
	"embed"
)

//go:embed "queries" "migrations" "static"
var EmbeddedFiles embed.FS

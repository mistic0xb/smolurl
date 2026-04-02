package static

import "embed"

//go:embed index.html styles.css
var StaticFiles embed.FS
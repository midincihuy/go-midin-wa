package main

import (
	"embed"
	"midincihuy/go-midin-wa/cmd"
)

//go:embed views/index.html
var embedIndex embed.FS

//go:embed views
var embedViews embed.FS

func main() {
	cmd.Execute(embedIndex, embedViews)
}

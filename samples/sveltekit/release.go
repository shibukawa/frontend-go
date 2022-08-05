//go:build release
// +build release

package main

import (
	"embed"

	"github.com/shibukawa/frontend-go"
)

//go:generate sh -c "cd frontend; npm run build"

//go:embed frontend/build/*
var asset embed.FS

func init() {
	frontend.SetFrontAsset(asset, frontend.Opt{
		FrameworkType: frontend.SvelteKit,
	})
}

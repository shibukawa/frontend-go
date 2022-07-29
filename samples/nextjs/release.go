//go:build release
// +build release

package main

import (
	"embed"

	"github.com/shibukawa/frontend-go"
)

//go:generate sh -c "cd frontend; npm run build"
//go:generate sh -c "cd frontend; npm exec next export"

//go:embed frontend/out/*
//go:embed frontend/out/_next/static/*/*
//go:embed frontend/out/_next/static/chunks/pages/*.js
var asset embed.FS

func init() {
	frontend.SetFrontAsset(asset)
}

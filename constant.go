package frontend

//go:generate enumer -type=Mode
type Mode int

const (
	Development Mode = iota
	Release
)

//go:generate enumer -type=FrameworkType
type FrameworkType int

const (
	AutoDetect FrameworkType = iota
	NextJS
	VueJS
	SvelteKit
	SolidJS
	skipDetect
	NotFound
)

type frameworkConfig struct {
	DistFolder       string
	DevServerCommand string
}

var frameworkConfigs = map[FrameworkType]frameworkConfig{
	NextJS: {
		DistFolder:       "out",
		DevServerCommand: "npm run dev",
	},
	VueJS: {
		DistFolder:       "dist",
		DevServerCommand: "npm run serve",
	},
	SvelteKit: {
		DistFolder:       "build",
		DevServerCommand: "npm run dev",
	},
	SolidJS: {
		DistFolder:       "dist",
		DevServerCommand: "npm run dev",
	},
}

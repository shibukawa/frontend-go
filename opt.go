package frontend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shibukawa/acquire-go"
)

// Opt specifies frontend configuration.
//
// frontend-go assumes frontend project is in "frontend" folder.
// And it tries to get other config as soon as possible.
// If you changed dist folder or build scripts and so on, use Opt and pass to
// [NewSPAHandler], [NewSPAHandlerFunc]
type Opt struct {
	FrontEndFolder       string        // Frontend application folder that contains package.json. default value is "frontend"
	SkipRunningDevServer bool          // Even if development mode, frontend-go doesn't run dev server
	FrameworkType        FrameworkType // NextJS, VueJS, SvelteKit, SolidJS is available instead of auto detect
	DistFolder           string        // Specify dist folder instead of auto detect
	Port                 uint16        // Specify port instead of auto detect
	DevServerCommand     string        // Specify dev server command instead of auto detect
	FallbackPath         string        // Specify fallback file path. Default is "index.html"
}

type packageJson struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (p packageJson) Has(f string) bool {
	if _, ok := p.Dependencies[f]; ok {
		return true
	}
	_, ok := p.DevDependencies[f]
	return ok
}

func normalizeOpt(currentFolder string, opts []Opt) (*Opt, error) {
	var opt Opt
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.FrontEndFolder == "" {
		opt.FrontEndFolder = "frontend"
	}
	if opt.FrameworkType == AutoDetect {
		matches, err := acquire.Acquire(acquire.File, filepath.Join(currentFolder, opt.FrontEndFolder, "package.json"))
		if err != nil {
			return nil, fmt.Errorf("package.json is not found under '%s' folder %w", opt.FrontEndFolder, err)
		}
		f, err := os.Open(matches[0])
		if err != nil {
			return nil, fmt.Errorf("file open error: '%s'", matches[0])
		}
		r := json.NewDecoder(f)
		var p packageJson
		err = r.Decode(&p)
		if err != nil {
			return nil, fmt.Errorf("json parse error: '%s'", matches[0])
		}
		if p.Has("@sveltejs/kit") {
			opt.FrameworkType = SvelteKit
		} else if p.Has("next") {
			opt.FrameworkType = NextJS
		} else if p.Has("vue") {
			opt.FrameworkType = VueJS
		} else if p.Has("solid-js") {
			opt.FrameworkType = SolidJS
		} else {
			opt.FrameworkType = NotFound
		}
	}
	if defaultConfig, ok := frameworkConfigs[opt.FrameworkType]; ok {
		if opt.DistFolder == "" {
			opt.DistFolder = defaultConfig.DistFolder
		}
		if opt.DevServerCommand == "" {
			opt.DevServerCommand = defaultConfig.DevServerCommand
		}
	}
	return &opt, nil
}

package frontend

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shibukawa/acquire-go"
)

var ErrPackageJsonNotFound = errors.New("package.json not found")

// Opt specifies frontend configuration.
//
// frontend-go assumes frontend project is in "frontend" folder.
// And it tries to get other config as soon as possible.
// If you changed dist folder or build scripts and so on, use Opt and pass to
// [NewSPAHandler], [NewSPAHandlerFunc]
type Opt struct {
	FrontEndFolderName   string        // Frontend application folder name that contains package.json. Default value is "frontend"
	FrontEndFolderPath   string        // Absolute frontend application folder that contains package.json.
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
	if opt.FrontEndFolderName == "" && opt.FrontEndFolderPath == "" {
		opt.FrontEndFolderName = "frontend"
	}
	if opt.FrontEndFolderPath == "" {
		matches, err := acquire.AcquireFromUnder(acquire.File, currentFolder, "", filepath.Join(opt.FrontEndFolderName, "package.json"))
		if errors.Is(err, acquire.ErrNotFound) {
			return nil, fmt.Errorf("package.json is not found under '%s' folder %w", opt.FrontEndFolderName, ErrPackageJsonNotFound)
		} else if err != nil {
			return nil, err
		}
		opt.FrontEndFolderPath = filepath.Dir(matches[0])
	} else {
		_, opt.FrontEndFolderName = filepath.Split(opt.FrontEndFolderPath)
		if _, err := os.Stat(filepath.Join(opt.FrontEndFolderPath, "package.json")); os.IsNotExist(err) {
			return nil, fmt.Errorf("package.json is not found under '%s' folder %w", opt.FrontEndFolderName, ErrPackageJsonNotFound)
		}
	}
	if opt.FallbackPath == "" {
		opt.FallbackPath = "index.html"
	}
	if opt.FrameworkType == AutoDetect {
		packageJsonPath := filepath.Join(opt.FrontEndFolderPath, "package.json")
		f, err := os.Open(packageJsonPath)
		if err != nil {
			return nil, fmt.Errorf("file open error: '%s'", packageJsonPath)
		}
		r := json.NewDecoder(f)
		var p packageJson
		err = r.Decode(&p)
		if err != nil {
			return nil, fmt.Errorf("json parse error: '%s'", packageJsonPath)
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

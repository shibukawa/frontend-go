package frontend

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/shibukawa/acquire-go"
	"github.com/stretchr/testify/assert"
)

func Test_normalizeOpt(t *testing.T) {
	type args struct {
		opts []Opt
	}
	testDataPaths := acquire.MustAcquire(acquire.Dir, "testdata")
	tests := []struct {
		name          string
		currentFolder string
		args          args
		wantErr       error
		want          *Opt
	}{
		{
			name:          "empty opt autofill",
			currentFolder: filepath.Join(testDataPaths[0], "emptyproject"),
			args: args{
				opts: []Opt{},
			},
			wantErr: nil,
			want: &Opt{
				FrontEndFolderName: "frontend",
				FrontEndFolderPath: filepath.Join(testDataPaths[0], "emptyproject", "frontend"),
				DevServerCommand:   "",
				FrameworkType:      NotFound,
				FallbackPath:       "index.html",
			},
		},
		{
			name:          "filled opt",
			currentFolder: filepath.Join(testDataPaths[0], "emptyproject"),
			args: args{
				opts: []Opt{
					{
						FrameworkType:        SvelteKit,
						DistFolder:           ".dist",
						FrontEndFolderName:   "web",
						Port:                 3000,
						SkipRunningDevServer: true,
						DevServerCommand:     "yarn dev",
						FallbackPath:         "200.html",
					},
				},
			},
			wantErr: nil,
			want: &Opt{
				FrameworkType:        SvelteKit,
				DistFolder:           ".dist",
				FrontEndFolderName:   "web",
				FrontEndFolderPath:   filepath.Join(testDataPaths[0], "emptyproject", "web"),
				Port:                 3000,
				SkipRunningDevServer: true,
				DevServerCommand:     "yarn dev",
				FallbackPath:         "200.html",
			},
		},
		{
			name:          "detect ancestor folder",
			currentFolder: filepath.Join(testDataPaths[0], "emptyproject/subfolder"),
			args: args{
				opts: []Opt{},
			},
			wantErr: nil,
			want: &Opt{
				FrontEndFolderName: "frontend",
				FrontEndFolderPath: filepath.Join(testDataPaths[0], "emptyproject", "frontend"),
				DevServerCommand:   "",
				FrameworkType:      NotFound,
				FallbackPath:       "index.html",
			},
		},
		{
			name:          "frontend abs path is specified",
			currentFolder: filepath.Join(testDataPaths[0], "emptyproject"),
			args: args{
				opts: []Opt{
					{
						FrontEndFolderPath: filepath.Join(testDataPaths[0], "emptyproject", "frontend"),
					},
				},
			},
			wantErr: nil,
			want: &Opt{
				FrontEndFolderName: "frontend",
				FrontEndFolderPath: filepath.Join(testDataPaths[0], "emptyproject", "frontend"),
				DevServerCommand:   "",
				FrameworkType:      NotFound,
				FallbackPath:       "index.html",
			},
		},
		{
			name:          "error: frontend abs path is specified (invalid)",
			currentFolder: filepath.Join(testDataPaths[0], "emptyproject"),
			args: args{
				opts: []Opt{
					{
						FrontEndFolderPath: filepath.Join(testDataPaths[0], "emptyproject", "not-exists"),
					},
				},
			},
			wantErr: ErrPackageJsonNotFound,
		},
		{
			name:          "error: frontend folder name is specified (invalid)",
			currentFolder: filepath.Join(testDataPaths[0], "emptyproject"),
			args: args{
				opts: []Opt{
					{
						FrontEndFolderName: "not-exists",
					},
				},
			},
			wantErr: ErrPackageJsonNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeOpt(tt.currentFolder, tt.args.opts)
			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				if err == nil {
					assert.Equal(t, *tt.want, *got)
				}
			}
		})
	}
}

func Test_normalizeOpt_Detect(t *testing.T) {
	samplesPaths := acquire.MustAcquire(acquire.Dir, "samples")
	tests := []struct {
		name          string
		projectFolder string
		want          *Opt
		wantErr       error
	}{
		{
			name:          "detect SvelteKit",
			projectFolder: "samples/sveltekit",
			want: &Opt{
				FrameworkType:      SvelteKit,
				DistFolder:         "build",
				FrontEndFolderName: "frontend",
				FrontEndFolderPath: filepath.Join(samplesPaths[0], "sveltekit", "frontend"),
				DevServerCommand:   "npm run dev",
				FallbackPath:       "index.html",
			},
		},
		{
			name:          "detect Next.js",
			projectFolder: "samples/nextjs",
			want: &Opt{
				FrameworkType:      NextJS,
				DistFolder:         "out",
				FrontEndFolderName: "frontend",
				FrontEndFolderPath: filepath.Join(samplesPaths[0], "nextjs", "frontend"),
				DevServerCommand:   "npm run dev",
				FallbackPath:       "index.html",
			},
		},
		{
			name:          "detect Vue.js",
			projectFolder: "samples/vuejs",
			want: &Opt{
				FrameworkType:      VueJS,
				DistFolder:         "dist",
				FrontEndFolderName: "frontend",
				FrontEndFolderPath: filepath.Join(samplesPaths[0], "vuejs", "frontend"),
				DevServerCommand:   "npm run serve",
				FallbackPath:       "index.html",
			},
		},
		{
			name:          "detect Solid.js",
			projectFolder: "samples/solidjs",
			want: &Opt{
				FrameworkType:      SolidJS,
				DistFolder:         "dist",
				FrontEndFolderName: "frontend",
				FrontEndFolderPath: filepath.Join(samplesPaths[0], "solidjs", "frontend"),
				DevServerCommand:   "npm run dev",
				FallbackPath:       "index.html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeOpt(tt.projectFolder, nil)
			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

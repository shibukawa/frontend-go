package frontend

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_normalizeOpt(t *testing.T) {
	type args struct {
		opts []Opt
	}
	tests := []struct {
		name          string
		projectFolder string
		args          args
		wantErr       error
		want          *Opt
	}{
		{
			name:          "empty opt autofill",
			projectFolder: "testdata/emptyproject",
			args: args{
				opts: []Opt{},
			},
			wantErr: nil,
			want: &Opt{
				FrontEndFolder:   "frontend",
				DevServerCommand: "",
				FrameworkType:    NotFound,
			},
		},
		{
			name:          "filled opt",
			projectFolder: "testdata/emptyproject",
			args: args{
				opts: []Opt{
					{
						FrameworkType:        SvelteKit,
						DistFolder:           ".dist",
						FrontEndFolder:       "web",
						Port:                 3000,
						SkipRunningDevServer: true,
						DevServerCommand:     "yarn dev",
					},
				},
			},
			wantErr: nil,
			want: &Opt{
				FrameworkType:        SvelteKit,
				DistFolder:           ".dist",
				FrontEndFolder:       "web",
				Port:                 3000,
				SkipRunningDevServer: true,
				DevServerCommand:     "yarn dev",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeOpt(tt.projectFolder, tt.args.opts)
			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, *tt.want, *got)
		})
	}
}

func Test_normalizeOpt_Detect(t *testing.T) {
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
				FrameworkType:    SvelteKit,
				DistFolder:       "build",
				FrontEndFolder:   "frontend",
				DevServerCommand: "npm run dev",
			},
		},
		{
			name:          "detect Next.js",
			projectFolder: "samples/nextjs",
			want: &Opt{
				FrameworkType:    NextJS,
				DistFolder:       "out",
				FrontEndFolder:   "frontend",
				DevServerCommand: "npm run dev",
			},
		},
		{
			name:          "detect Vue.js",
			projectFolder: "samples/vuejs",
			want: &Opt{
				FrameworkType:    VueJS,
				DistFolder:       "dist",
				FrontEndFolder:   "frontend",
				DevServerCommand: "npm run serve",
			},
		},
		{
			name:          "detect Solid.js",
			projectFolder: "samples/solidjs",
			want: &Opt{
				FrameworkType:    SolidJS,
				DistFolder:       "dist",
				FrontEndFolder:   "frontend",
				DevServerCommand: "npm run dev",
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

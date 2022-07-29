# frontend-go

 SPA (Single Page Application) style web application frontend helper for go server. This package works with Next.js, Vue.js, SvelteKit, Solid.js.

## Development Mode

It runs dev server behind go server. It works as reverse proxy. You can get all features that the web framework dev server provides (hot module replacement and so on).

![Development Mode](/docs/dev.png)

## Production Mode

It gets prebuilt assets (HTML, JS, CSS) from go's embed. Basically, all SPA needs to fallback to `index.html` when there is not have assets.

![Release Mode](/docs/rel.png)

## Integration

frontend-go assumes the following structure. Important points are the following:

* There is frontend project folder named `frontend`.
* There is `release.go` at the parent folder of frontend folder.
* Add frontend-go's handler to your web server.

```txt
awesome-your-web-app
├── LICENSE
├── README.md
├── cmd
│   └── server
│      └── main.go    # Add frontend-go's handler
├── go.mod
├── go.sum
├── api.go
├── release.go        # You should add this
└── frontend          # frontend project (Next.js, Vue.js, SvelteKit, Solid.js)
    └── package.json
```

`release.go` has the following code.

```go:release.go
//go:build release

package webapp

// ↓this folder is decided by your frontend framework

//go:embed frontend/build/*
var asset embed.FS

init() {
    frontend.SetFrontAsset(asset)
}
```

Use with `net/http`:

```go:cmd/server/main.go
package main

func handler() http.Handler {
    mux := http.NewServeMux()
    mux.Handle("/api", YourAPIHandler)
    mux.Handle("/", webfront.MustNewSPAHandler())
    return mux
}
```

Use with chi:

```go:cmd/server/main.go

import (
    "github.com/go-chi/chi/v5"
)

func handler() http.Handler  {
    r := chi.NewRouter()
    r.Post("/api", YourAPIHandler)
    r.NotFound(webfront.MustNewHandlerFunc())
    return r
}
```

```sh
# Development mode
$ go run main.go

# Production build
$ go build -tags release
```

## Web Frameworks Specific Instructions

### Next.js 12

frontend-go assumes Next.js's static generation result (that does'n need Node.js/bun to run).

To start project with Next.js, init project like this:

```sh
$ go mod init yourapp
$ mkdir -p cmd/yourapp
$ npx create-next-app@latest --ts frontend
```

To enable static generation, add `unoptimized` flag to `next.config.js`

```js:next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  experimental: { // Add here
    images: {
      unoptimized: true,
    }
  }
}
```

You should build frontend project by using the following commands.

```bash
$ npm run build
$ npm exec next export
```

The go:embed directive comment should be the following by default:

```go:release.go
//go:embed frontend/out/*
//go:embed frontend/out/_next/static/*/*
//go:embed frontend/out/_next/static/chunks/pages/*.js
var asset embed.FS
```

### Vue.js

To start project with Vue.js, init project like this:

```sh
$ go mod init yourapp
$ mkdir -p cmd/yourapp
$ vue create frontend
```

You should build frontend project by using the following commands.

```bash
$ npm run build
```

The go:embed directive comment should be the following by default:

```go:release.go
//go:embed frontend/dist/*
var asset embed.FS
```

### SvelteKit

Default SvelteKit provides frontend program that requires Node.js to run. To work with this module, you should configure the front end project with [static site mode](https://kit.svelte.dev/docs/adapters#supported-environments-static-sites).

```
$ go mod init yourapp
$ mkdir -p cmd/server
$ npm create yourapp@latest frontend
```

To modify front end site, add `static-adaptor`.

```sh
$ npm install @sveltejs/adapter-static
```

```js:frontend/svelte.config.js
import adapter from '@sveltejs/adapter-static'; // Modify here

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: preprocess(),

  kit: {
    adapter: adapter({       // Modify here for static HTML/JS generation
      fallback: 'index.html'
    }),
    prerender: {             // Modify here for SPA mode
      default: false
    },
    trailingSlash: 'always'  // Modify here
  }
};

export default config;
```

You should build frontend project by using the following commands.

```bash
$ npm run build
```

The go:embed directive comment should be the following by default:

```go:release.go
//go:embed frontend/build/*
var asset embed.FS
```


### Solid.js

```
$ go mod init yourapp
$ mkdir -p cmd/server
$ npx degit solidjs/templates/ts frontend
```

You should build frontend project by using the following commands.

```bash
$ npm run build
```

The go:embed directive comment should be the following by default:

```go:release.go
//go:embed frontend/dist/*
var asset embed.FS
```

## Configuration

`webfront.NewSPAHandler()` has option that modifies package's behavior. 

If you put frontend project at `frontend` folder and doesn't change npm scripts, you don't have to modify configuration.

```go
handler := webfront.NewSPAHandler(ctx, webfront.Opt{
    FrontEndFolder: "frontend",              // Frontend application folder that contains package.json. default value is "frontend"
    ProjectType:    webfront.AutoDetect,     // NextJS, SvelteKit, VueJS, SolidJS is available
    SkipRunningDevServer:     false,         // Skip running dev server even if development mode
    DistFolder:     "",                      // Specify dist folder instead of auto detect
    Port:           0,                       // Specify port instead of auto detect
    DevelopmentCommand: "npm run dev",       // Specify dev server command instead of auto detect
    FallbackPath:       string               // Specify fallback file path. Default is "index.html"
})
```

## Credits

Yoshiki Shibukawa

## License

Apache 2
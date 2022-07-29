//go:build dev
// +build dev

package frontend

func init() {
	mode = Development
	/*u, err := url.Parse("http://localhost:3000/")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", httputil.NewSingleHostReverseProxy(u))*/

}

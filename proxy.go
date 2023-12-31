// poc: Forge short links
// take simple endpoints and redirect them to helpful destinations

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type ProxyResponse struct {
}

// ProxyHandler is a struct that implements the http.Handler interface
type ProxyHandler struct {
	proxyList  map[string]*httputil.ReverseProxy
	shortLinks map[string]string
}

// ServeHTTP is a method of the ProxyHandler struct. It implements the http.Handler interface. It logs the url and the linked endpoint.
func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("serving: %s at %s", r.URL.Path, ph.shortLinks[r.URL.Path])
	// http.Redirect(w, r, ph.shortLinks[r.URL.Path], http.StatusMovedPermanently)
	// r.URL.Path = ph.shortLinks[r.URL.Path]
	ph.proxyList[r.URL.Path].ServeHTTP(w, r)
}

// return map of short links and their destinations
func loadShortLinks(filename string) (map[string]string, error) {
	// load short links from file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode json into map
	var shortLinks map[string]string
	err = json.NewDecoder(file).Decode(&shortLinks)
	if err != nil {
		return nil, err
	}
	return shortLinks, nil
}

func (ph *ProxyHandler) hostReverseProxy() error {
	for k, v := range ph.shortLinks {
		remote, err := url.Parse(v)
		if err != nil {
			return err
		}
		// set up reverse proxy
		proxy := &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetXForwarded()
				r.SetURL(remote)
			},
		}
		ph.proxyList[k] = proxy
		http.Handle(k, ph)
	}
	return nil
}

func main() {
	filename := flag.String("filename", "shortlinks.json", "file with endpoints and their destinations")
	flag.Parse()

	// load short links from file
	shortLinks, err := loadShortLinks(*filename)
	if err != nil {
		log.Fatal(err)
	}
	proxyHandler := &ProxyHandler{
		shortLinks: shortLinks,
		proxyList:  make(map[string]*httputil.ReverseProxy),
	}

	// create reverse proxy for each short link
	err = proxyHandler.hostReverseProxy()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "welcome to the forge proxy server.\n")
		fmt.Fprintf(w, "available endpoints:\n")
		for k := range proxyHandler.shortLinks {
			fmt.Fprintf(w, "%s\n", k)
		}
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

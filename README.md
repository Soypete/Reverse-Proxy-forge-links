# Reverse-Proxy-forge-links

This is a reverse proxy server used to forward forge custom links to public urls. It would almost like browser bookmarks.

--

## Questions

* How could someone get started with your codebase?
  * this software can be run in 3 ways (way 3 prefered as it has links file) 

```bash
 curl -LJO https://github.com/Soypete/Reverse-Proxy-forge-links/releases/download/v0.5.0/Reverse-Proxy-forge-links_Darwin_arm64.tar.gz
 tar -xfz Reverse-Proxy-forge-links_Darwin_arm64.tar.gz
 ./Reverse-Proxy-forge-links
```

```bash
go install github.com/Soypete/Reverse-Proxy-forge-links@v0.5.0
Reverse-Proxy-forge-links 
```

```bash
git checkout
go build .
./Reverse-Proxy-forge-links -filename="forgelinks.json"
```

* What resources did you use to build your implementation?
  * I have heard the term reverse proxy but didn't know full definition https://en.wikipedia.org/wiki/Reverse_proxy
  * Look at OSS implementation https://github.com/caddyserver/caddy/blob/master/modules/caddyhttp/reverseproxy/reverseproxy.go
  * Example https://gist.github.com/thurt/2ae1be5fd12a3501e7f049d96dc68bb9
  * Docs https://pkg.go.dev/net/http/httputil@go1.20.5#ReverseProxy
* Explain any design decisions you made, including limitations of the system.
  * I wanted this example to have a purpose, so I replicated a shortlink [server](https://www.golinks.io/) that I have used at previous jobs. I know that this will not send to https since I am only serving an http without certificates, but I thought it a decent POC as all the routes are being called. I chose to use a json file instead of a DB since it is easier for I/O and testing. I wanted a test endpoint like `/` that would allow me to check that it was live since the `https:/` will return a 404.
* How would you scale this?
  * first thing would be adding a sqlite db to manage links. This will allow users to add and delete links easier. Also add auth so that members of forge can have access to links and can share links as they needed. A simple sharing token can be generated as an parameter or a ehader arguement to allow one time access sharing. most of the traffic I expect would be handled by go routing and connection pooling so no need to add additional proxying for requests.
* How would you make it more secure?
  * as mentioned above, adding auth to requests. Since this is publish basic auth would suffice with a list of users/passwords in the db for verification. The urls are all public so nothing needs to be protected. If we start using this for internal docs then something like tailscale and device sharing would scale really well.

## todo:

* make a way to save the short links and their destinations.
* enable forwarding to https (http -> https returns 404)
  * something like caddy would probably solve this
* enable access to _forge members_

// # TCP Proxy
// We want to listen on a TCP socket. And forward the traffic to a Unix domain
// socket.
package main

var services proxy.ServiceDesc{{Subdomain: "blog"}; {Subdomain: "www"}}

func main(){}

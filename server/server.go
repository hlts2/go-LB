package server

import (
	"crypto/tls"
	"net/http"

	"github.com/hlts2/go-LB/config"
	iphash "github.com/hlts2/ip-hash"
	"github.com/hlts2/least-connections"
	"github.com/hlts2/round-robin"
	"github.com/kpango/glg"
)

// Balancing is custom type of balancing algorithm
type Balancing interface{}

// LBServer represents load balancing server object
type LBServer struct {
	*http.Server
	balancing Balancing
}

// NewLBServer returns LBServer object
func NewLBServer(addr string) *LBServer {
	lbs := new(LBServer)
	lbs.Addr = addr
	return lbs
}

func (lbs *LBServer) getLeastConnections() leastconnections.LeastConnections {
	return lbs.balancing.(leastconnections.LeastConnections)
}

func (lbs *LBServer) getRoundRobin() roundrobin.RoundRobin {
	return lbs.balancing.(roundrobin.RoundRobin)
}

func (lbs *LBServer) getIPHash() iphash.IPHash {
	return lbs.balancing.(iphash.IPHash)
}

// Build builds LB config
func (lbs *LBServer) Build(conf config.Config) *LBServer {
	switch conf.Balancing {
	case "ip-hash":
		ih, err := iphash.New(conf.Servers.ToStringSlice())
		if err == nil {
			lbs.balancing = ih
		}
		lbs.Handler = http.HandlerFunc(lbs.leastConnectionsBalancing)
	case "round-robin":
		rr, err := roundrobin.New(conf.Servers.ToStringSlice())
		if err == nil {
			lbs.balancing = rr
		}
		lbs.Handler = http.HandlerFunc(lbs.roundRobinBalancing)
	case "least-connections":
		lc, err := leastconnections.New(conf.Servers.ToStringSlice())
		if err == nil {
			lbs.balancing = lc
		}
		lbs.Handler = http.HandlerFunc(lbs.ipHashBalancing)
	default:
		glg.Fatal("balancing algorithm dose not found")
	}

	return lbs
}

// ListenAndServeTLS runs load balancing server with TLS
func (lbs *LBServer) ListenAndServeTLS(tlsConfig *tls.Config, certFile, keyFile string) error {
	lbs.TLSConfig = tlsConfig

	err := lbs.Server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		return err
	}

	return nil
}

// ListenAndServe runs load balancing server
func (lbs *LBServer) ListenAndServe() error {
	err := lbs.Server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

package meta

import (
	"crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/influxdata/influxdb/monitor/diagnostics"
	"github.com/influxdata/influxdb/toml"
)

const (
	// DefaultHostname is the default hostname if one is not provided.
	DefaultHostname = "localhost"

	// DefaultRaftBindAddress is the default address to bind to.
	DefaultRaftBindAddress = ":8089"

	// DefaultHTTPBindAddress is the default address to bind the API to.
	DefaultHTTPBindAddress = ":8091"

	// DefaultGossipFrequency is the default frequency with which the node will gossip its known announcements.
	DefaultGossipFrequency = 5 * time.Second

	// DefaultAnnouncementExpiration is the default length of time an announcement is kept before it is considered too old.
	DefaultAnnouncementExpiration = 30 * time.Second

	// DefaultElectionTimeout is the default election timeout for the store.
	DefaultElectionTimeout = 1000 * time.Millisecond

	// DefaultHeartbeatTimeout is the default heartbeat timeout for the store.
	DefaultHeartbeatTimeout = 1000 * time.Millisecond

	// DefaultLeaderLeaseTimeout is the default leader lease for the store.
	DefaultLeaderLeaseTimeout = 500 * time.Millisecond

	// DefaultConsensusTimeout is the default consensus timeout for the store.
	DefaultConsensusTimeout = 30 * time.Second

	// DefaultCommitTimeout is the default commit timeout for the store.
	DefaultCommitTimeout = 50 * time.Millisecond

	// DefaultLeaseDuration is the default duration for leases.
	DefaultLeaseDuration = 60 * time.Second

	// DefaultLoggingEnabled determines if log messages are printed for the meta service.
	DefaultLoggingEnabled = true
)

// Config represents the meta configuration.
type Config struct {
	MetaTLSEnabled           bool   `toml:"meta-tls-enabled"`
	MetaInsecureTLS          bool   `toml:"meta-insecure-tls"`
	MetaAuthEnabled          bool   `toml:"meta-auth-enabled"`
	MetaInternalSharedSecret string `toml:"meta-internal-shared-secret"`

	Dir string `toml:"dir"`

	RetentionAutoCreate bool `toml:"retention-autocreate"`
	LoggingEnabled      bool `toml:"logging-enabled"`

	// RemoteHostname is the hostname portion to use when registering meta node
	// addresses.  This hostname must be resolvable from other nodes.
	RemoteHostname string `toml:"-"`

	// SingleServer is used to start the meta server in single server mode.
	SingleServer bool `toml:"-"`

	// TLS is a base tls config to use for https clients.
	TLS *tls.Config `toml:"-"`

	// BindAddress is the bind address(port) for meta node communication
	BindAddress string `toml:"bind-address"`

	AuthEnabled bool `toml:"auth-enabled"`
	LDAPAllowed bool `toml:"ldap-allowed"`

	// HTTPBindAddress is the bind address for the metaservice HTTP API
	HTTPBindAddress  string `toml:"http-bind-address"`
	HTTPSEnabled     bool   `toml:"https-enabled"`
	HTTPSCertificate string `toml:"https-certificate"`
	HTTPSPrivateKey  string `toml:"https-private-key"`
	HTTPSInsecureTLS bool   `toml:"https-insecure-tls"`

	DataUseTLS      bool `toml:"data-use-tls"`
	DataInsecureTLS bool `toml:"data-insecure-tls"`

	GossipFrequency        toml.Duration `toml:"gossip-frequency"`
	AnnouncementExpiration toml.Duration `toml:"announcement-expiration"`

	ElectionTimeout    toml.Duration `toml:"election-timeout"`
	HeartbeatTimeout   toml.Duration `toml:"heartbeat-timeout"`
	LeaderLeaseTimeout toml.Duration `toml:"leader-lease-timeout"`
	ConsensusTimeout   toml.Duration `toml:"consensus-timeout"`
	CommitTimeout      toml.Duration `toml:"commit-timeout"`
	ClusterTracing     bool          `toml:"cluster-tracing"`
	PprofEnabled       bool          `toml:"pprof-enabled"`
	LeaseDuration      toml.Duration `toml:"lease-duration"`

	SharedSecret         string `toml:"shared-secret"`
	InternalSharedSecret string `toml:"internal-shared-secret"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		RetentionAutoCreate:    true,
		LoggingEnabled:         DefaultLoggingEnabled,
		BindAddress:            DefaultRaftBindAddress,
		HTTPBindAddress:        DefaultHTTPBindAddress,
		GossipFrequency:        toml.Duration(DefaultGossipFrequency),
		AnnouncementExpiration: toml.Duration(DefaultAnnouncementExpiration),
		ElectionTimeout:        toml.Duration(DefaultElectionTimeout),
		HeartbeatTimeout:       toml.Duration(DefaultHeartbeatTimeout),
		LeaderLeaseTimeout:     toml.Duration(DefaultLeaderLeaseTimeout),
		ConsensusTimeout:       toml.Duration(DefaultConsensusTimeout),
		CommitTimeout:          toml.Duration(DefaultCommitTimeout),
		PprofEnabled:           true,
		LeaseDuration:          toml.Duration(DefaultLeaseDuration),
	}

}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if c.Dir == "" {
		return errors.New("Meta.Dir must be specified")
	}
	return nil
}

// Diagnostics returns a diagnostics representation of a subset of the Config.
func (c *Config) Diagnostics() (*diagnostics.Diagnostics, error) {
	return diagnostics.RowFromMap(map[string]interface{}{
		"dir": c.Dir,
	}), nil
}

func RemoteAddr(hostname, addr string) string {
	if hostname == "" {
		hostname = DefaultHostname
	}
	remote, err := DefaultHost(hostname, addr)
	if err != nil {
		return addr
	}
	return remote
}

func DefaultHost(hostname, addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	if host == "" || host == "0.0.0.0" || host == "::" {
		return net.JoinHostPort(hostname, port), nil
	}
	return addr, nil
}

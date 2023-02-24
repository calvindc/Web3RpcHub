package mainimpl

// SvrCfg_PrintVersion
var SvrCfg_PrintVersion = true
var SvrCfg_PrintVersion_I = "print version number, build date, git commit, and compiler version"

// SvrCfg_SecretHandsharkeKey shs
var SvrCfg_SecretHandsharkeKey = "1KHLiKZvAvjbY1ziZEHMXawbCEIM6qwjCDm3VYRan/s="
var SvrCfg_SecretHandsharkeKey_I = "secret handshake key, if change makes you part of a different network"

// SvrCfg_SecretHandsharkeKey
var SvrCofg_ListenAddrShsMux = ":8008"
var SvrCofg_ListenAddrShsMux_I = "address to listen on for secret handshake + muxrpc"

// SvrCfg_SecretHandsharkeKey
var SvrCfg_ListenAddrHttp = "8001"
var SvrCfg_ListenAddrHttp_I = "address to listen on for HTTP requests"

// SvrCfg_EnableUnixSock
var SvrCfg_EnableUnixSock = false
var SvrCfg_EnableUnixSock_I = "disable or enable the UNIX socket RPC interface"

// SvrCfg_EnableUnixSock
var SvrCfg_RepoDir = ".web3rpchub"
var SvrCfg_RepoDir_I = "where to put the log and indexes"

// SvrCfg_LogDir
var SvrCfg_LogDir = "logs"
var SvrCfg_LogDir_I = "where to write debug output to"

// SvrCfg_ListenAddrMetricsPprof as prometheus
var SvrCfg_ListenAddrMetricsPprof = "localhost:8002"
var SvrCfg_ListenAddrMetricsPprof_I = "prometheus, listen addr for metrics and pprof HTTP server"

// SvrCfg_HttpsDomain
var SvrCfg_HttpsDomain = ""
var SvrCfg_HttpsDomain_I = "which domain to use for TLS and AllowedHosts checks"

// SvrCfg_HubMode
var SvrCfg_HubMode = func(val string) error {
	return nil
}
var SvrCfg_HubMode_I = "the privacy mode(values: open, community, restricted) determining hub access controls"

//SvrCfg_AliasesAsSubdomains
var SvrCfg_AliasesAsSubdomains = true
var SvrCfg_AliasesAsSubdomains_I = "needs to be disabled if a wildcard certificate for the hub is not available"

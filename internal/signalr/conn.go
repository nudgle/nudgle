package signalr

import (
	"net/url"
)

// ParsedConnString is the structure extracted from a SignalR connection string
type ParsedConnString struct {
	Endpoint *url.URL
	Key      string
	Version  string
}

// ParseConnectionString will parse the SignalR connection string from the Azure Portal
func ParseConnectionString(connStr string) (*ParsedConnString, error) {
	parsed := new(ParsedConnString)
	u, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}
	parsed.Endpoint = u
	parsed.Key = ""
	parsed.Version = ""
	return parsed, nil
}

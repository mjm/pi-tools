package nomadic

import (
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
)

func ConsulUpstream(dest string, port int) *nomad.ConsulUpstream {
	return &nomad.ConsulUpstream{
		DestinationName: dest,
		LocalBindPort:   port,
	}
}

func NewServiceDefaults(name string, protocol string) *consul.ServiceConfigEntry {
	return &consul.ServiceConfigEntry{
		Kind:     consul.ServiceDefaults,
		Name:     name,
		Protocol: protocol,
	}
}

func NewServiceIntentions(name string, sources ...*consul.SourceIntention) *consul.ServiceIntentionsConfigEntry {
	return &consul.ServiceIntentionsConfigEntry{
		Kind:    consul.ServiceIntentions,
		Name:    name,
		Sources: sources,
	}
}

func AppAwareIntention(name string, permissions ...*consul.IntentionPermission) *consul.SourceIntention {
	return &consul.SourceIntention{
		Name:        name,
		Precedence:  9,
		Type:        consul.IntentionSourceConsul,
		Permissions: permissions,
	}
}

func AllowHTTP(http *consul.IntentionHTTPPermission) *consul.IntentionPermission {
	return &consul.IntentionPermission{
		Action: consul.IntentionActionAllow,
		HTTP:   http,
	}
}

func DenyHTTP(http *consul.IntentionHTTPPermission) *consul.IntentionPermission {
	return &consul.IntentionPermission{
		Action: consul.IntentionActionDeny,
		HTTP:   http,
	}
}

func DenyAllHTTP() *consul.IntentionPermission {
	return DenyHTTP(HTTPPathPrefix("/"))
}

func HTTPPathPrefix(prefix string) *consul.IntentionHTTPPermission {
	return &consul.IntentionHTTPPermission{
		PathPrefix: prefix,
	}
}

func HTTPPathExact(path string) *consul.IntentionHTTPPermission {
	return &consul.IntentionHTTPPermission{
		PathExact: path,
	}
}

func AllowIntention(name string) *consul.SourceIntention {
	precedence := 9
	if name == "*" {
		precedence = 8
	}

	return &consul.SourceIntention{
		Name:       name,
		Action:     consul.IntentionActionAllow,
		Precedence: precedence,
		Type:       consul.IntentionSourceConsul,
	}
}

func DenyIntention(name string) *consul.SourceIntention {
	precedence := 9
	if name == "*" {
		precedence = 8
	}

	return &consul.SourceIntention{
		Name:       name,
		Action:     consul.IntentionActionDeny,
		Precedence: precedence,
		Type:       consul.IntentionSourceConsul,
	}
}

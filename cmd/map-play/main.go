package main

import "log"

type Provider struct {
	name         string
	capabilities []string
}

func main() {
	providers := []Provider{
		{
			"nats-core",
			[]string{"publish", "subscription", "request/reply"},
		},
		{
			"nats-jetstream",
			[]string{"jetstream-publish", "consumer"},
		},
		{
			"nats-kv",
			[]string{"keyvalue"},
		},
	}
	flattenCapability := func(providers []Provider) []string {
		capabilities := []string{}
		for _, provider := range providers {
			capabilities = append(capabilities, provider.capabilities...)
		}
		return capabilities
	}
	log.Println(flattenCapability(providers))
}

package component

type ComponentConfig struct {
	RootPath, ComponentName                  string
	NatsCore, NatsJetstream, NatsKv, License bool
}

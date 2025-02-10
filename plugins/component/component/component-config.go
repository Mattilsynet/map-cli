package component

type Config struct {
	Path, ComponentName, Repository string
	Capabilities                    []string
	ComponentNatsConn,              // these booleans should be deducted from Capabilities list above
	ComponentNatsJetstream,
	ComponentNatsKeyValue bool

	WitPackage, // can be deducted from Repository
	ImportNatsCoreWit,
	ExportNatsCoreWit,
	ImportNatsJetstreamWit,
	ExportNatsJetstreamWit,
	ImportNatsKvWit bool
}

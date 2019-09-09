type DoNotSend struct {
	What       string   `yaml:"what"`
	Peers string   `yaml:"peers"`
	Community       int64   `yaml:"community"`
}

type SetLocPref struct {
	Value       int   `yaml:"value"`
	Destination string   `yaml:"dest"`
	Community       int64   `yaml:"community"`
}

type NoExport struct {
	DoNotSends []DoNotSend `yaml:"donotsend"`
	SetLocPrefs []SetLocPref `yaml:"setlocpref"`
}

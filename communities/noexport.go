package communities
type DoNotSend struct {
	What       string   `yaml:"what"`
	Peers string   `yaml:"peers"`
	Community       int   `yaml:"community"`
}

type SetLocPref struct {
	Value       int   `yaml:"value"`
	Destination string   `yaml:"dest"`
	Community       int   `yaml:"community"`
}

type NoExport struct {
	DoNotSends []DoNotSend `yaml:"donotsend"`
	SetLocPrefs []SetLocPref `yaml:"setlocpref"`
}

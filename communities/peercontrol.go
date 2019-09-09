package communities
type DoNotAnnounce struct {
	Peer string `yaml:"peer"`
	Asn int64	`yaml:"asn"`
	Community int64 `yaml:"community"`
}

type Prepend struct {
	What string	`yaml:"what`
	Times int	`yaml:"times"`
	Community int64 `yaml:"community"`
}

type PeerControl struct {
	DoNotAnnounces []DoNotAnnounce `yaml:"donotannounce"`
	Prepends []Prepend `yaml:"prepend"`
}

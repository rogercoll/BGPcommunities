package communities
type DoNotAnnounce struct {
	Peer string `yaml:"peer"`
	Asn int	`yaml:"asn"`
	Community int `yaml:"community"`
}

type Prepend struct {
	What string	`yaml:"what`
	Times int	`yaml:"times"`
	Community int `yaml:"community"`
}

type PeerControl struct {
	DoNotAnnounces []DoNotAnnounce `yaml:"donotannounce"`
	Prepends []Prepend `yaml:"prepend"`
}

package communities

type SetCustomerRoute struct {
	Value       int   `yaml:"value"`
	Community       int64   `yaml:"community"`
}

type LocalPreference struct {
	SetCustomersRoute []SetCustomerRoute `yaml:"setcustomerroute"`
}


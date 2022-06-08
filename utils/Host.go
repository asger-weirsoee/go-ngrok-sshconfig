package utils

type Host struct {
	Name string
	Usr  string
	Host string
	Port string
}

func NewHost(name string, usr string, res *ApiRes) *Host {
	if name == "" || usr == "" || res.Host == "" || res.Port == "" {
		panic("Some or all of the fields are empty or not set correctly")
	}
	return &Host{Name: name, Usr: usr, Host: res.Host, Port: res.Port}
}

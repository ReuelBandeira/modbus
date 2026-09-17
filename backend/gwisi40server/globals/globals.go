package globals

type Configurations struct {
	JWT struct {
		Expiration int `json:"expiration"`
	}
}

var Conf = &Configurations{}

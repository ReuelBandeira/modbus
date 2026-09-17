package status

type Dev struct {
	Id     uint
	Status bool
}

var Devs []Dev

func GetDeviceStatus(id uint) (int, Dev, bool) {

	for i, v := range Devs {
		if v.Id == id {
			return i, v, true
		}
	}
	return -1, Dev{}, false
}

func SetDeviceStatus(id uint, status bool) {

	i, _, ok := GetDeviceStatus(id)

	if ok {
		Devs[i] = Dev{Id: id, Status: status}
	} else {
		Devs = append(Devs, Dev{Id: id, Status: status})
	}

}

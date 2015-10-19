package manager

type MiddlemanInterface interface {
	GenInnerKey(string, interface{}) string
	SerializeValue(string, interface{}) string
	UnserializeValue(string, interface{}) string
}

var Middlemen = make(map[string]MiddlemanInterface)

func MiddlemanRegister(id string, thisMiddleman MiddlemanInterface) bool {
	if _, ok := Middlemen[id]; ok == true {
		return false
	}
	Middlemen[id] = thisMiddleman
	return true
}

func init() {

}

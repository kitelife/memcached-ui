package manager

type MiddlemanInterface interface {
	Config(map[string]string) bool
	GenInnerKey(string) string
	SerializeValue(string) string
	UnserializeValue(string) interface{}
}

var Middlemen = make(map[string]MiddlemanInterface)

func MiddlemanRegister(id string, thisMiddleman MiddlemanInterface) bool {
	if _, ok := Middlemen[id]; ok == true {
		return false
	}
	Middlemen[id] = thisMiddleman
	return true
}

func Get(id string, config map[string]string) MiddlemanInterface {
	targetMiddleman, ok := Middlemen[id]
	if ok == false {
		return nil
	}
	if targetMiddleman.Config(config) == false {
		return nil
	}
	return targetMiddleman
}

func init() {

}

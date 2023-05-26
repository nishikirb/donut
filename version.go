package donut

var version = "source"

func SetVersion(v string) {
	version = v
}

func GetVersion() string {
	return version
}

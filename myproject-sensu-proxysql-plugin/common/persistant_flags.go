package common

type PersistantOptions struct {
	Json           bool
	JsonPretty     bool
	DefaultsFile   string
	User           string
	Pass           string
	Host           string
	Port           uint64
	MydbToken      string
	EncryptedCreds bool
}

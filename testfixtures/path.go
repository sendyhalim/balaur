package testfixtures

const TomlConfigPath string = "testfixtures/testconfig.toml"
const TestAppPath string = "testfixtures/postapp"

var TestAppConfigPath = map[string]string{
	"app":        "config.toml",
	"route":      "routes.toml",
	"middleware": "middlewares.toml",
}

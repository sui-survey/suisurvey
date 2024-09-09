package config

type config struct {
	System system `yaml:"system"`
	Logger logger `yaml:"logger"`
}

type system struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var Config *config

func Init() {
	Config = &config{}
}

type logger struct {
	Level          string `yaml:"level"`
	Prefix         string `yaml:"prefix"`
	Director       string `yaml:"director"`
	ShowLineNumber bool   `yaml:"show_line_number"`
	LogInConsole   bool   `yaml:"log_in_console"`
}

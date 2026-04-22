package section

type Monitor struct {
	LogLevel    string `default:"debug" split_word:"true"`
	Environment string `default:"development"`
}

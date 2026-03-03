package section

type Monitor struct {
	LogLevel    string `default:"debag" split_word:"true"`
	Environment string `default:"development"`
}

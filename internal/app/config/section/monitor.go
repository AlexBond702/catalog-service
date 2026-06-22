package section

type Monitor struct {
	Prometheus  MonitorPrometheus
	LogLevel    string `default:"debug" split_word:"true"`
	Environment string `default:"development"`
}
type MonitorPrometheus struct {
	Enabled bool `default:"false"`
}

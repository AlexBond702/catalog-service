package section

type (
	Processor struct {
		WebServer ProcessorWebServer
	}
	ProcessorWebServer struct {
		ListenPort uint32 `required:"true" default:"8080" split_words:"true"`
	}
)

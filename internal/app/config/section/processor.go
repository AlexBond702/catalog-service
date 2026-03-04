package section

type (
	Processor struct {
		WebServer ProcessorWebServer `required:"true" split_words:"true"`
	}
	ProcessorWebServer struct {
		ListenPort uint32 `required:"true" default:"8080" split_words:"true"`
	}
)

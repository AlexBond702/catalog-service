package section

type (
	Processor struct {
		Gateway   ProcessorGateway
		Grpc      ProcessorGrpc
		WebServer ProcessorWebServer `required:"true" split_words:"true"`
	}
	ProcessorWebServer struct {
		ListenPort uint32 `required:"true" default:"8080" split_words:"true"`
	}
	ProcessorGrpc struct {
		Host       string `required:"true"`
		ListenPort uint32 `required:"true" split_words:"true"`
	}
	ProcessorGateway struct {
		Host         string `required:"true"`
		ListenPort   uint32 `required:"true" split_words:"true"`
		GrpcEndpoint string `required:"true" split_words:"true"`
	}
)

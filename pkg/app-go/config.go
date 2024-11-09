package appgo

type (
	Config struct {
		HttpPort   int
		GrpcPort   int
		EnableGrpc bool
	}
)

func (conf *Config) ApplyDefaults() {

}

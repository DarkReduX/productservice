package config

type Server struct {
	ListenAddress string `env:"LISTEN_ADDRESS,notEmpty" envDefault:"0.0.0.0:8081"`
}

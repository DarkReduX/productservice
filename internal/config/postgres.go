package config

type Postgres struct {
	URL string `env:"POSTGRES_URL" envDefault:"postgres://postgres:postgres@localhost:5432/products"`
}

package envs

import "crg.eti.br/go/config"

type Envs struct {
	Environment       string `cfg:"ENVIRONMENT" cfgRequired:"true" cfgDefault:"local"`
	Service           string `cfg:"SERVICE" cfgRequired:"true" cfgDefault:"r2-challenge"`
	Version           string `cfg:"VERSION" cfgRequired:"true" cfgDefault:"1"`
	HTTPHost          string `cfg:"HTTP_HOST" cfgRequired:"true" cfgDefault:"localhost"`
	HTTPPort          string `cfg:"HTTP_PORT" cfgRequired:"true" cfgDefault:"8080"`
	ReadHeaderTimeout string `cfg:"READ_HEADER_TIMEOUT" cfgDefault:"15s"`
	HTTPTimeout       string `cfg:"HTTP_TIMEOUT" cfgDefault:"10s"`

	DBHost            string `cfg:"DB_HOST" cfgDefault:"localhost"`
	DBPort            string `cfg:"DB_PORT" cfgDefault:"5432"`
	DBUser            string `cfg:"DB_USER" cfgDefault:"postgres"`
	DBPassword        string `cfg:"DB_PASSWORD" cfgDefault:"postgres"`
	DBName            string `cfg:"DB_NAME" cfgDefault:"r2_db"`
	DBSSLMode         string `cfg:"DB_SSLMODE" cfgDefault:"disable"`
	DBMaxOpenConns    int    `cfg:"DB_MAX_OPEN_CONNS" cfgDefault:"10"`
	DBMaxIdleConns    int    `cfg:"DB_MAX_IDLE_CONNS" cfgDefault:"5"`
	DBConnMaxLifetime string `cfg:"DB_CONN_MAX_LIFETIME" cfgDefault:"1h"`
}

func NewEnvs() (Envs, error) {
	var e Envs
	return e, config.Parse(&e)
}


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
}

func NewEnvs() (Envs, error) {
	var e Envs
	return e, config.Parse(&e)
}


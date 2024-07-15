package config

import (
	"os"

	"main/internal/logger"

	"github.com/caarlos0/env/v10"
)

type cors struct {
	AllowedOrigins []string `env:"ALLOWED_ORIGINS,required,notEmpty" envSeparator:","`
}

var (
	Cors = &cors{}
)

func init() {
	toParse := []any{Cors}
	errors := []error{}

	for _, v := range toParse {
		if err := env.Parse(v); err != nil {
			if er, ok := err.(env.AggregateError); ok {
				errors = append(errors, er.Errors...)

				continue
			}

			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		logger.Stderr.Error("errors found while parsing environment variables", logger.ErrorsAttr(errors...))

		os.Exit(1)
	}
}

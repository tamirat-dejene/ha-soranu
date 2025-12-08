package main

import (
	"github.com/tamirat-dejene/ha-soranu/shared/env"
)

type Service struct {
	Environment *env.Env
}

func NewService(env *env.Env) *Service {
	return &Service{
		Environment: env,
	}
}

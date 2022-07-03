package mixture

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Envs int

const (
	ForAnyEnv Envs = 1 << iota
	ForLocal
	ForDocker
	ForIntegration
	ForProduction
	ForTest
	ForSandbox
)

var DefaultConfig = Config{
	AllowedEnvironments:      ForAnyEnv,
	FailOnUnknownEnvironment: true,
}

type migration struct {
	migration *gormigrate.Migration
	envs      Envs
}

type Config struct {
	AllowedEnvironments      Envs
	FailOnUnknownEnvironment bool
}

type mixture struct {
	migrations []migration
	config     *Config
	db         *gorm.DB
}

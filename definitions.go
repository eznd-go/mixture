package mixture

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type envs int

const (
	ForAnyEnv envs = 1 << iota
	ForLocal
	ForDocker
	ForIntegration
	ForProduction
	ForTest
	ForSandbox
)

var DefaultOptions = config{
	AllowedEnvironments:      ForAnyEnv,
	FailOnUnknownEnvironment: true,
}

type migration struct {
	migration *gormigrate.Migration
	envs      envs
}

type config struct {
	AllowedEnvironments      envs
	FailOnUnknownEnvironment bool
}

type mixture struct {
	migrations []migration
	config     config
	db         *gorm.DB
}

package mixture

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var mx = NewMixture(DefaultOptions, nil)

func Add(e envs, m *gormigrate.Migration) {
	mx.Add(e, m)
}

func Apply(db *gorm.DB, env string) error {
	e, err := parseEnv(env)
	if err != nil {
		return err
	}

	return mx.
		SetDB(db).
		Apply(e)
}

func parseEnv(env string) (envs, error) {
	switch env {
	case "local":
		return ForLocal, nil
	case "docker", "ci":
		return ForDocker, nil
	case "int", "integration":
		return ForIntegration, nil
	case "prod", "production":
		return ForProduction, nil
	case "test":
		return ForTest, nil
	case "sandbox":
		return ForSandbox, nil

	}

	return 0, fmt.Errorf("failed to parse env %s", env)
}

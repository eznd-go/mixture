package static

import (
	"fmt"
	"github.com/ezn-go/mixture"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var mx = mixture.New(nil)

func Add(e mixture.Envs, m *gormigrate.Migration) {
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

func parseEnv(env string) (mixture.Envs, error) {
	switch env {
	case "local":
		return mixture.ForLocal, nil
	case "docker", "ci":
		return mixture.ForDocker, nil
	case "int", "integration":
		return mixture.ForIntegration, nil
	case "prod", "production":
		return mixture.ForProduction, nil
	case "test":
		return mixture.ForTest, nil
	case "sandbox":
		return mixture.ForSandbox, nil

	}

	return 0, fmt.Errorf("failed to parse env: %s", env)
}

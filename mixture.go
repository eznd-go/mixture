package mixture

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func NewMixture(c config, db *gorm.DB) *mixture {
	return &mixture{
		migrations: []migration{},
		config:     c,
		db:         db,
	}
}

func (m *mixture) Add(e envs, mig *gormigrate.Migration) *mixture {
	m.migrations = append(m.migrations, migration{
		envs:      e,
		migration: mig,
	})
	return m
}

func (m *mixture) SetDB(db *gorm.DB) *mixture {
	m.db = db
	return m
}

func (m *mixture) Apply(env envs) error {
	migrator := gormigrate.New(m.db, gormigrate.DefaultOptions, m.filter(env))
	return migrator.Migrate()
}

func (m *mixture) filter(env envs) []*gormigrate.Migration {
	var t []*gormigrate.Migration

	for _, r := range m.migrations {
		if r.envs == ForAnyEnv || env&r.envs == env {
			t = append(t, r.migration)
		}
	}

	return t
}

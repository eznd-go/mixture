package mixture

import (
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type M gormigrate.Migration

func (m *M) ToGormMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID:       m.ID,
		Migrate:  m.Migrate,
		Rollback: m.Rollback,
	}
}

func New(db *gorm.DB) *mixture {
	return &mixture{
		migrations: []migration{},
		config:     &DefaultConfig,
		db:         db,
	}
}

func NewWithConfig(db *gorm.DB, c *Config) *mixture {
	return &mixture{
		migrations: []migration{},
		config:     c,
		db:         db,
	}
}

func (m *mixture) Add(e Envs, mig *M) *mixture {
	m.migrations = append(m.migrations, migration{
		envs:      e,
		migration: mig.ToGormMigration(),
	})
	return m
}

func (m *mixture) SetDB(db *gorm.DB) *mixture {
	m.db = db
	return m
}

func (m *mixture) Apply(destEnv Envs) error {
	migrator := gormigrate.New(m.db, gormigrate.DefaultOptions, m.filter(destEnv))
	return errors.WithStack(migrator.Migrate())
}

func (m *mixture) ApplyVerbose(destEnv Envs) error {
	migrations := m.filter(destEnv)
	fmt.Printf("applying %s%d%s migrations...\n", Cyan, len(migrations), Reset)
	for _, migration := range migrations {
		migrationset := make([]*gormigrate.Migration, 0)
		migrationset = append(migrationset, migration)

		fmt.Printf("  %s‣%s %s... ", Yellow, Reset, migration.ID)

		err := gormigrate.New(m.db, gormigrate.DefaultOptions, migrationset).Migrate()
		if err != nil {
			fmt.Printf("%sFAIL:%s %v\n", Red, Reset, err)
			return errors.WithStack(err)
		}

		fmt.Printf("%sok%s\n", Green, Reset)
	}
	fmt.Printf("%s✔%s migrations applied\n", Green, Reset)
	return nil
}

func (m *mixture) filter(env Envs) []*gormigrate.Migration {
	var t []*gormigrate.Migration

	for _, r := range m.migrations {
		if r.envs == ForAnyEnv || env&r.envs == env {
			t = append(t, r.migration)
		}
	}

	return t
}

package mixture

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNewMixture(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *mixture
	}{
		{
			name: "Happy path",
			args: args{
				db: nil,
			},
			want: &mixture{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.db)
			if reflect.DeepEqual(got.migrations, tt.want.migrations) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mixture_Add(t *testing.T) {
	type fields struct {
		migrations []migration
		config     *Config
		db         *gorm.DB
	}
	type args struct {
		e   Envs
		mig []*M
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mixture
	}{
		{
			name: "Happy path",
			args: args{
				e: 0,
				mig: []*M{
					{ID: "1"},
					{ID: "2"},
					{ID: "3"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mixture{
				migrations: tt.fields.migrations,
				config:     tt.fields.config,
				db:         tt.fields.db,
			}
			for _, r := range tt.args.mig {
				m.Add(tt.args.e, r)
			}
			if got := len(m.migrations); got != len(tt.args.mig) {
				t.Errorf("Add(): count mismatch, got %v, want %v", got, len(tt.args.mig))
			}
		})
	}
}

func Test_mixture_filter(t *testing.T) {
	type fields struct {
		migrations []migration
		config     *Config
		db         *gorm.DB
	}
	type migs struct {
		e   Envs
		mig []*M
	}
	type args struct {
		env  Envs
		migs []migs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*gormigrate.Migration
	}{
		{
			name: "Adding simple migrations",
			fields: fields{
				migrations: nil,
				config:     &Config{},
				db:         nil,
			},
			args: args{
				env: ForAnyEnv,
				migs: []migs{
					{
						e: ForAnyEnv,
						mig: []*M{
							{ID: "0001"},
							{ID: "0002"},
							{ID: "0003"},
						},
					},
				},
			},
			want: []*gormigrate.Migration{
				{ID: "0001"},
				{ID: "0002"},
				{ID: "0003"},
			},
		},
		{
			name: "Mixing different Envs",
			fields: fields{
				migrations: nil,
				config:     &Config{},
				db:         nil,
			},
			args: args{
				env: ForProduction,
				migs: []migs{
					{
						e: ForProduction,
						mig: []*M{
							{ID: "0001"},
							{ID: "0002"},
							{ID: "0003"},
						},
					},
					{
						e: ForIntegration,
						mig: []*M{
							{ID: "0004"},
							{ID: "0005"},
							{ID: "0006"},
						},
					},
					{
						e: ForAnyEnv,
						mig: []*M{
							{ID: "0007"},
							{ID: "0008"},
							{ID: "0009"},
						},
					},
				},
			},
			want: []*gormigrate.Migration{
				{ID: "0001"},
				{ID: "0002"},
				{ID: "0003"},
				{ID: "0007"},
				{ID: "0008"},
				{ID: "0009"},
			},
		},
		{
			name: "Empty result set",
			fields: fields{
				migrations: nil,
				config:     &Config{},
				db:         nil,
			},
			args: args{
				env: ForLocal,
				migs: []migs{
					{
						e: ForProduction,
						mig: []*M{
							{ID: "0001"},
							{ID: "0002"},
							{ID: "0003"},
						},
					},
					{
						e: ForIntegration,
						mig: []*M{
							{ID: "0004"},
							{ID: "0005"},
							{ID: "0006"},
						},
					},
					{
						e: ForDocker,
						mig: []*M{
							{ID: "0007"},
							{ID: "0008"},
							{ID: "0009"},
						},
					},
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mixture{
				migrations: tt.fields.migrations,
				config:     tt.fields.config,
				db:         tt.fields.db,
			}
			for _, r := range tt.args.migs {
				for _, rr := range r.mig {
					m.Add(r.e, rr)
				}
			}
			if got := m.filter(tt.args.env); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

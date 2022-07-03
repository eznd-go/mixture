package mixture

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"reflect"
)

func CreateTable(t interface{}) gormigrate.MigrateFunc {
	return func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(t)
	}
}

func DropTable(t interface{}) gormigrate.RollbackFunc {
	return func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(t)
	}
}

// CreateBatch imports slice of structs into DB, may have lots of limitations
// I realize, that playing with reflect is not the best performance trick
// But, since migrations are not expected to be used in high-loaded code,
// I believe that syntax sugar and readability >> speed here
func CreateBatch(b interface{}) gormigrate.MigrateFunc {
	t := reflect.TypeOf(b)
	if t.Kind() != reflect.Slice {
		panic("CreateBatch: input is not a slice")
	}
	return func(tx *gorm.DB) error {
		s := reflect.ValueOf(b)
		for i := 0; i < s.Len(); i++ {
			tt := s.Index(i).Type()
			v := s.Index(i).Convert(tt)
			tm := v.MethodByName("TableName")
			if !tm.IsValid() {
				panic("CreateBatch: input does not implement TableName() func")
			}
			tns := tm.Call(nil)
			if len(tns) == 0 {
				panic("CreateBatch: input does not implement TableName() func")
			}
			tn := tns[0].String()

			vv := make(map[string]interface{})
			for j := 0; j < s.Index(i).NumField(); j++ {
				vv[strcase.ToSnake(tt.Field(j).Name)] = s.Index(i).Field(j).Interface()
			}

			err := tx.Table(tn).Create(vv).Error
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func DeleteBatch(b interface{}) gormigrate.RollbackFunc {
	t := reflect.TypeOf(b)
	if t.Kind() != reflect.Slice {
		panic("CreateBatch: input is not a slice")
	}
	return func(tx *gorm.DB) error {
		s := reflect.ValueOf(b)
		for i := 0; i < s.Len(); i++ {
			tt := s.Index(i).Type()
			v := s.Index(i).Convert(tt)
			tm := v.MethodByName("TableName")
			if !tm.IsValid() {
				panic("DeleteBatch: input does not implement TableName() func")
			}
			tns := tm.Call(nil)
			if len(tns) == 0 {
				panic("DeleteBatch: input does not implement TableName() func")
			}
			tn := tns[0].String()

			vv := make(map[string]interface{})
			for j := 0; j < s.Index(i).NumField(); j++ {
				vv[strcase.ToSnake(tt.Field(j).Name)] = s.Index(i).Field(j).Interface()
			}

			err := tx.Table(tn).Delete(vv).Error
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Update(table, where, column, value string) gormigrate.MigrateFunc {
	return func(tx *gorm.DB) error {
		return tx.Table(table).Where(where).Update(column, value).Error
	}
}

func RollbackUpdate(table, where, column, value string) gormigrate.RollbackFunc {
	return func(tx *gorm.DB) error {
		return tx.Table(table).Where(where).Update(column, value).Error
	}
}

package mixture

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"reflect"
)

type Func = func(tx *gorm.DB) error

func createTable(t interface{}) Func {
	return func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(t)
	}
}

func CreateTableM(t interface{}) gormigrate.MigrateFunc {
	return createTable(t)
}

func CreateTableR(t interface{}) gormigrate.RollbackFunc {
	return createTable(t)
}

func dropTable(t interface{}) Func {
	return func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(t)
	}
}

func DropTableM(t interface{}) gormigrate.MigrateFunc {
	return dropTable(t)
}

func DropTableR(t interface{}) gormigrate.RollbackFunc {
	return dropTable(t)
}

// CreateBatchM imports slice of structs into DB, may have lots of limitations
// I realize, that playing with reflect is not the best performance trick
// But, since migrations are not expected to be used in high-loaded code,
// I believe that syntax sugar and readability >> speed here
func createBatch(b interface{}) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		t := reflect.TypeOf(b)
		if t.Kind() != reflect.Slice {
			panic("CreateBatch: input is not a slice")
		}

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

func CreateBatchM(b interface{}) gormigrate.MigrateFunc {
	return createBatch(b)
}

func CreateBatchR(b interface{}) gormigrate.RollbackFunc {
	return createBatch(b)
}

func deleteBatch(b interface{}) func(tx *gorm.DB) error {
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

			found := false
			for j := 0; j < s.Index(i).NumField(); j++ {
				if strcase.ToSnake(tt.Field(j).Name) == "id" {
					err := tx.Table(tn).Where("id = ?", s.Index(i).Field(j).Interface()).Delete(s.Index(i).Field(j).Type()).Error
					if err != nil {
						return err
					}
					found = true
				}
			}

			if !found {
				panic("DeleteBatch: input does not have id column")
			}
		}

		return nil
	}
}

func DeleteBatchM(b interface{}) gormigrate.MigrateFunc {
	return deleteBatch(b)
}

func DeleteBatchR(b interface{}) gormigrate.RollbackFunc {
	return deleteBatch(b)
}

func update(table, where, column, value string) Func {
	return func(tx *gorm.DB) error {
		return tx.Table(table).Where(where).Update(column, value).Error
	}
}

func UpdateM(table, where, column, value string) gormigrate.MigrateFunc {
	return update(table, where, column, value)
}

func UpdateR(table, where, column, value string) gormigrate.RollbackFunc {
	return update(table, where, column, value)
}

func delete(table, where string) Func {
	return func(tx *gorm.DB) error {
		return tx.Table(table).Delete(where).Error
	}
}

func DeleteM(table, where string) gormigrate.MigrateFunc {
	return delete(table, where)
}

func DeleteR(table, where string) gormigrate.RollbackFunc {
	return delete(table, where)
}

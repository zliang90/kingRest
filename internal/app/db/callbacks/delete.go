package callbacks

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// DeleteCallback overwrite default delete callback
func DeleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedField, hasDeletedField := scope.FieldByName("Deleted")
		deletedAtField, hasDeletedAtField := scope.FieldByName("DeletedAt")

		if !hasDeletedField {
			return
		}

		if !scope.Search.Unscoped {
			sql := fmt.Sprintf("UPDATE %v SET ", scope.QuotedTableName())
			// deleted
			sql += fmt.Sprintf(
				"%v=%v",
				scope.Quote(deletedField.DBName),
				scope.AddToVars(1),
			)
			// deleted at
			if hasDeletedAtField {
				sql += fmt.Sprintf(
					",%v=%v",
					scope.Quote(deletedAtField.DBName),
					scope.AddToVars(time.Now()),
				)
			}
			// extra option
			sql += addExtraSpaceIfExist(scope.CombinedConditionSql())
			if extraOption != "" {
				sql += addExtraSpaceIfExist(extraOption)
			}
			scope.Raw(sql).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

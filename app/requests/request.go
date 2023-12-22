package requests

import (
	"GoBlog/pkg/model"
	"errors"
	"fmt"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

// 此方法会在初始化时执行
func init(){
	//not_exits:users,email
	govalidator.AddCustomRule("not_exits", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exits:"), ",")

		tableName := rng[0]
		dbFiled := rng[1]
		val := value.(string)

		var count int64
		model.DB.Table(tableName).Where(dbFiled+" = ?", val).Count(&count)

		if count != 0 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已被占用", val)
		}

		return nil
	})
}
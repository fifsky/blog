package validate

//
// import (
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestValidate(t *testing.T) {
//
// 	t.Run("验证通过", func(t *testing.T) {
// 		data := struct {
// 			Name string `validate:"omitempty,oneof=Y N" label:"名称"`
// 		}{
// 			Name: "Y",
// 		}
//
// 		err := Validate(data)
// 		assert.Nil(t, err)
// 	})
//
// 	t.Run("验证不通过", func(t *testing.T) {
// 		data := struct {
// 			Name string `validate:"omitempty,oneof=Y N" label:"名称"`
// 		}{
// 			Name: "ABC",
// 		}
//
// 		err := Validate(data)
// 		assert.Error(t, err)
// 		t.Logf("error: %s", err)
// 	})
//
// 	t.Run("验证时间格式", func(t *testing.T) {
// 		data := struct {
// 			Date string `validate:"required,datetime=2006-01-02" label:"开始时间"`
// 		}{
// 			Date: "2022-01-01 1",
// 		}
//
// 		t.Logf("errors: %s", Validate(data))
// 	})
//
// 	t.Run("验证测试", func(t *testing.T) {
// 		type d struct {
// 			StartTime string `validate:"omitempty,datetime=2006-01-02 15:04:05"`
// 		}
//
// 		d1 := &d{}
// 		err := Validate(d1)
// 		t.Logf("error: %s", err)
// 	})
//
// 	t.Run("var", func(t *testing.T) {
// 		err := validate.Var("http://test.cn?pass=1", "http_url")
// 		assert.NoError(t, err)
//
// 		err = validate.Var("", "http_url")
// 		assert.Error(t, err)
// 	})
// }

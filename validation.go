package margo

import (
	"github.com/gin-gonic/gin/binding"
	"reflect"
	"gopkg.in/go-playground/validator.v8"
)

type SpecificBinding interface {
	getBinding() binding.Binding
}

func bodyParamsValidator(params interface{}) HandlerFunc {
	return paramsValidator(params, bodyParams, binding.JSON)
}

func queryParamsValidator(params interface{}) HandlerFunc {
	return paramsValidator(params, queryParams, binding.Query)
}

func paramsValidator(params interface{}, key string, deflt binding.Binding) HandlerFunc {
	return func(c *Context) Response {
		if params == nil {
			return nil
		}

		var b binding.Binding
		if p, ok := params.(SpecificBinding); ok {
			// if SpecifiedBinding interface is implemented,
			// use the specified binding
			b = p.getBinding()
		} else {
			// otherwise, use default binding
			b = deflt
		}

		instance := reflect.New(reflect.TypeOf(params)).Interface()

		if err := c.ShouldBindWith(instance, b); err != nil {
			var errs []*MargoError

			// ValidationErrors is a map[string]*FieldError
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, val := range ve {
					errs = append(errs, InvalidParamsError(&val.Name, &val.ActualTag))
				}
			} else {
				errs = append(errs, InvalidParamsError(nil, nil))
			}

			return BadRequest(errs...)
		}

		c.Set(key, instance)
		return nil
	}
}

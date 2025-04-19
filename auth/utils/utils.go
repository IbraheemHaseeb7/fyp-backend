package utils

import (
	"reflect"

	"github.com/IbraheemHaseeb7/pubsub"
)

var Requests map[string](chan pubsub.PubsubMessage)

func StructToMap(data interface{}) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    v := reflect.ValueOf(data)
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        // Get the JSON tag
        jsonTag := fieldType.Tag.Get("json")
        if jsonTag != "" && jsonTag != "-" { // Exclude fields with "-" JSON tag
            result[jsonTag] = field.Interface()
        }
    }

    return result, nil
}

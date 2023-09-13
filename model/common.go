package model

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

func LoadJsonSchema(name string, i interface{}) (err error) {
	var b []byte
	if b, err = os.ReadFile(fmt.Sprintf("./model/%s.json", name)); err != nil {
		return
	}
	return bson.UnmarshalExtJSON(b, true, &i)
}

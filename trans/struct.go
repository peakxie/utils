// Package trans ...
package trans

import (
	"bytes"
	"encoding/json"
)

func StructIntf(st any) (any, error) {
	bts, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(bts))
	decoder.UseNumber()
	// 反序列化JSON到map[string]interface{}
	var mt map[string]any
	err = decoder.Decode(&mt)
	if err != nil {
		return nil, err
	}

	return mt, nil
}

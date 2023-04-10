package entity

import (
	"encoding/json"
	"testing"
)

func TestTime(t *testing.T) {
	now := Now()
	byt, err := json.Marshal(&now)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(byt, &now); err != nil {
		panic(err)
	}

}

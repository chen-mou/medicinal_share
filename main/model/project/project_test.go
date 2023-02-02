package project

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetHospitalByNear(t *testing.T) {
	byt, _ := json.Marshal(GetHospitalByNear(10.5, 200, 1, 500))
	fmt.Println(string(byt))
}

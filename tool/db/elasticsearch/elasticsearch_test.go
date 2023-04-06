package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestIndex(t *testing.T) {
	type doctorTag struct {
		DoctorId  int64   `json:"doctor_id"`
		Ids       []int64 `json:"ids"`
		Longitude float64
		Latitude  float64
	}
	tag := doctorTag{
		123456,
		[]int64{1, 2, 3, 4, 5},
		12.477,
		24.899,
	}
	byt, _ := json.Marshal(tag)
	GetClient().Indices.Create("doctor_tag")
	res, _ := GetClient().Index("doctor_tag",
		bytes.NewBuffer(byt),
		GetClient().Index.WithDocumentID(GetRandomId("doctor_tag")))
	res, _ = GetClient().Search(GetClient().Search.WithIndex("doctor_tag"), GetClient().Search.WithSource())
	byt, _ = ioutil.ReadAll(res.Body)
	fmt.Println(string(byt))
}

package elasticsearch

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"math/rand"
	"medicinal_share/tool/encrypt/md5"
	"net/http"
	"time"
)

var client *elasticsearch.Client

var config = elasticsearch.Config{
	Addresses: []string{
		"localhost:9300",
		"localhost:9301",
		"localhost:9302",
	},
	Transport: &http.Transport{MaxIdleConns: 10},
}

func init() {
	var err error
	client, err = elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
}

func GetClient() *elasticsearch.Client {
	return client
}

func Save(index string, data any) error {
	var err error
	var b []byte
	if b, err = json.Marshal(data); err != nil {
		if _, err = client.Indices.Create(index, client.Indices.Create.WithBody(bytes.NewReader(b))); err != nil {
			return nil
		}
	}
	return err
}

func Get() error {
	return nil
}

func GetRandomId(name string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	random := make([]rune, 10)
	for i := 0; i < 10; i++ {
		random[i] = 'a' + rune(rand.Intn(26))
	}
	return md5.Hash(now + string(random) + name)
}

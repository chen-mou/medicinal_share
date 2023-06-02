package elasticsearch

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io/ioutil"
	"math/rand"
	"medicinal_share/tool/encrypt/md5"
	"net/http"
	"sync"
	"time"
)

var client *elasticsearch.Client

var config = elasticsearch.Config{
	Addresses: []string{
		"http://localhost:9200",
		"http://localhost:9201",
		"http://localhost:9202",
	},
	Transport: &http.Transport{MaxIdleConns: 10},
}

func GetClient() *elasticsearch.Client {
	once := &sync.Once{}
	once.Do(func() {
		var err error
		client, err = elasticsearch.NewClient(config)
		if err != nil {
			panic(err)
		}
	})
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

func Get(v any, o ...func(*esapi.SearchRequest)) error {
	var err error
	var res *esapi.Response
	if res, err = client.Search(o...); err == nil {
		var byt []byte
		if byt, err = ioutil.ReadAll(res.Body); err == nil {
			err = json.Unmarshal(byt, v)
		}
	}
	return err
}

func GetRandomId(name string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	random := make([]rune, 10)
	for i := 0; i < 10; i++ {
		random[i] = 'a' + rune(rand.Intn(26))
	}
	return md5.Hash(now + string(random) + name)
}

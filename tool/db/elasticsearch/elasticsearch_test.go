package elasticsearch

import "testing"

func TestIndex(t *testing.T) {
	GetClient().Indices.Create("api_test")
}

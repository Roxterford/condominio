package filter_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

const json_query = `
{
  "status": "active",
  "volume": {
    "gt": 50,
    "lt": 100
  },
  "or": [
    {
      "status": {
        "like": "active"
      },
      "and": [
        {
          "volume": {
            "gt": 50
          }
        }
      ]
    },
    {
      "not": {
        "volume": 100
      }
    }
  ]
}
`

func TestQueryJSON(t *testing.T) {

	_, ok := filter.NewQueryFrom(json_query)
	if !ok {
		t.Fatal("failed to create query")
	}

}

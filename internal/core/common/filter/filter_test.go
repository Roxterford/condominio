package filter_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
)

const complicated_json_query = `{"and":[{"account_status":"verified"},{"login_attempts":{"le":3}},{"balance":{"bt":[1000,50000]}},{"or":[{"username":{"rx":"^[a-zA-Z0-9._%+-]+@company\\.com$"}},{"and":[{"tags":{"in":["premium","vip","early_adopter"]}},{"bio":{"li":"%expert%"}},{"not":{"display_name":{"leq":5}}}]}]},{"and":[{"password_strength":{"lge":12}},{"comments_content":{"not":{"llt":10}}}]},{"or":[{"age":{"ge":18}},{"internal_rank":{"gt":100}},{"metadata.description":{"lgt":50}}]},{"last_activity":{"lt":"2024-01-01T00:00:00Z"}},{"not":{"or":[{"region":{"eq":"restricted_zone"}},{"legacy_id":{"rx":"old_.*"}}]}}]}`

type generic_filterable struct {
}

func (generic_filterable) FilterSpec() filter.Spec {
	return filter.Spec{
		"oop":                  filter.String,
		"foo":                  filter.String,
		"bar":                  filter.String,
		"account_status":       filter.String,
		"login_attempts":       filter.Int,
		"balance":              filter.Float,
		"username":             filter.String,
		"tags":                 filter.String,
		"bio":                  filter.String,
		"display_name":         filter.String,
		"password_strength":    filter.Int,
		"comments_content":     filter.String,
		"age":                  filter.Int,
		"internal_rank":        filter.Int,
		"metadata.description": filter.String,
		"last_activity":        filter.String,
		"region":               filter.String,
		"legacy_id":            filter.String,
	}
}

func TestNewFilter(t *testing.T) {

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "simple query",
			query: `{"oop": "bar"}`,
		},
		{
			name:  "complicated query",
			query: complicated_json_query,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			query, ok := filter.NewQueryFrom(test.query)
			if !ok {
				t.Fatal("failed to create query")
			}

			logger.Debug("===Test creating filter===\n")
			_, err := filter.New[generic_filterable](query)
			if err != nil {
				t.Fatal("failed to create filter")
			}
			logger.Debug("===Test created filter===\n")

		})
	}

}

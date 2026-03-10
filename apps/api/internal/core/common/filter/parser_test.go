package filter_test

import (
	"encoding/json"
	"testing"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string // Usamos JSON string para simular entrada real
		wantError bool
		check     func(t *testing.T, node filter.Clause)
	}{
		{
			name:  "Simple Equality",
			input: `{"status": "active"}`,
			check: func(t *testing.T, node filter.Clause) {
				p, ok := node.(*filter.PredicateClause)
				if !ok {
					t.Fatalf("Expected PredicateNode, got %T", node)
				}
				if p.Field != "status" || p.Condition != filter.CONDITON_EQ || p.Value != "active" {
					t.Errorf("Mismatch parsing simple equality: %+v", p)
				}
			},
		},
		{
			name:  "Compound Condition",
			input: `{"age": {"gt": 18}}`,
			check: func(t *testing.T, node filter.Clause) {
				p, ok := node.(*filter.PredicateClause)
				if !ok {
					t.Fatalf("Expected PredicateNode")
				}
				if p.Field != "age" || p.Condition != filter.CONDITON_GT {
					t.Errorf("Mismatch parsing compound: %+v", p)
				}
			},
		},
		{
			name:  "Logical AND Implicit",
			input: `{"status": "active", "age": {"gte": 18}}`,
			check: func(t *testing.T, node filter.Clause) {
				l, ok := node.(*filter.LogicalClause)
				if !ok || l.Operator != filter.OPERATOR_AND {
					t.Fatalf("Expected LogicalNode AND due to multiple keys")
				}
				if len(l.Children) != 2 {
					t.Errorf("Expected 2 children, got %d", len(l.Children))
				}
			},
		},
		{
			name:  "Logical OR Explicit",
			input: `{"or": [{"status": "pending"}, {"status": "failed"}]}`,
			check: func(t *testing.T, node filter.Clause) {
				l, ok := node.(*filter.LogicalClause)
				if !ok || l.Operator != filter.OPERATOR_OR {
					t.Fatalf("Expected LogicalNode OR")
				}
				if len(l.Children) != 2 {
					t.Errorf("Expected 2 children for OR")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simular decodificación inicial como lo haría un controlador HTTP
			var raw map[string]any
			_ = json.Unmarshal([]byte(tt.input), &raw)

			got, err := filter.Parse(raw)
			if (err != nil) != tt.wantError {
				t.Errorf("Parse() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

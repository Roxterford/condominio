package filter_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery_Validate_Simple(t *testing.T) {
	tests := []struct {
		name        string
		query       filter.Query
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid simple string field",
			query:       filter.Query{"status": "active"},
			expectError: false,
		},
		{
			name:        "valid simple numeric field",
			query:       filter.Query{"price": 100},
			expectError: false,
		},
		{
			name:        "valid simple float field",
			query:       filter.Query{"price": 100.50},
			expectError: false,
		},
		{
			name:        "valid simple boolean field",
			query:       filter.Query{"enabled": true},
			expectError: false,
		},
		{
			name:        "empty query should be valid",
			query:       filter.Query{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQuery_Validate_Compound(t *testing.T) {
	tests := []struct {
		name        string
		query       filter.Query
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid compound string field with equals",
			query:       filter.Query{"status": map[string]any{"eq": "active"}},
			expectError: false,
		},
		{
			name:        "valid compound string field with like",
			query:       filter.Query{"name": map[string]any{"li": "john%"}},
			expectError: false,
		},
		{
			name:        "valid compound numeric field with greater than",
			query:       filter.Query{"price": map[string]any{"gt": 100}},
			expectError: false,
		},
		{
			name:        "valid compound numeric field with less than",
			query:       filter.Query{"price": map[string]any{"lt": 200}},
			expectError: false,
		},
		{
			name:        "valid compound numeric field with greater or equal",
			query:       filter.Query{"price": map[string]any{"ge": 100}},
			expectError: false,
		},
		{
			name:        "valid compound numeric field with less or equal",
			query:       filter.Query{"price": map[string]any{"le": 200}},
			expectError: false,
		},
		{
			name:        "valid compound boolean field with equals",
			query:       filter.Query{"enabled": map[string]any{"eq": true}},
			expectError: false,
		},
		{
			name:        "invalid condition for string field",
			query:       filter.Query{"status": map[string]any{"gt": "active"}},
			expectError: true,
			errorMsg:    "Condicion 'gt' no valida para el campo string 'status'",
		},
		{
			name:        "invalid condition for numeric field",
			query:       filter.Query{"price": map[string]any{"li": 100}},
			expectError: true,
			errorMsg:    "Condicion 'li' no valida para el campo numerico 'price'",
		},
		{
			name:        "invalid condition for boolean field",
			query:       filter.Query{"enabled": map[string]any{"gt": true}},
			expectError: true,
			errorMsg:    "Condicion 'gt' no valida para el campo bool 'enabled'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQuery_Validate_Operators(t *testing.T) {
	tests := []struct {
		name        string
		query       filter.Query
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid AND operator with simple queries",
			query: filter.Query{
				"and": []any{
					filter.Query{"status": "active"},
					map[string]any{"price": 100},
				},
			},
			expectError: false,
		},
		{
			name: "valid AND operator with compound queries",
			query: filter.Query{
				"and": []any{
					filter.Query{"status": map[string]any{"eq": "active"}},
					filter.Query{"price": map[string]any{"gt": 100}},
				},
			},
			expectError: false,
		},
		{
			name: "valid OR operator with simple queries",
			query: filter.Query{
				"or": []any{
					filter.Query{"status": "active"},
					filter.Query{"enabled": true},
				},
			},
			expectError: false,
		},
		{
			name: "valid NOT operator with simple query",
			query: filter.Query{
				"not": filter.Query{"status": "inactive"},
			},
			expectError: false,
		},
		{
			name: "valid NOT operator with compound query",
			query: filter.Query{
				"not": filter.Query{"price": map[string]any{"lt": 50}},
			},
			expectError: false,
		},
		{
			name: "nested operators",
			query: filter.Query{
				"and": []any{
					filter.Query{"status": "active"},
					filter.Query{
						"or": []any{
							filter.Query{"price": map[string]any{"gt": 100}},
							filter.Query{"discount": true},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "AND operator with non-array value",
			query:       filter.Query{"and": filter.Query{"status": "active"}},
			expectError: true,
			errorMsg:    "Operador 'and' debe ser un array",
		},
		{
			name:        "OR operator with non-array value",
			query:       filter.Query{"or": "invalid"},
			expectError: true,
			errorMsg:    "Operador 'or' debe ser un array",
		},
		{
			name:        "NOT operator with non-object value",
			query:       filter.Query{"not": "invalid"},
			expectError: true,
			errorMsg:    "Operador 'not' debe ser un objeto",
		},
		{
			name: "AND operator with invalid child",
			query: filter.Query{
				"and": []any{
					"invalid",
				},
			},
			expectError: true,
			errorMsg:    "Condicion 'and[0]' debe ser un objeto",
		},
		{
			name: "AND operator with nested error",
			query: filter.Query{
				"and": []any{
					filter.Query{"status": "active"},
					filter.Query{"price": map[string]any{"li": 100}}, // invalid: like on numeric field
				},
			},
			expectError: true,
			errorMsg:    "Condicion 'li' no valida para el campo numerico 'price'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQuery_ValidateWithSpec(t *testing.T) {
	tests := []struct {
		name        string
		query       filter.Query
		spec        filter.Spec
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid field in whitelist",
			query:       filter.Query{"status": "active"},
			spec:        filter.Spec{"status": filter.String, "price": filter.Int},
			expectError: false,
		},
		{
			name:        "valid compound field in whitelist",
			query:       filter.Query{"price": map[string]any{"gt": 100}},
			spec:        filter.Spec{"status": filter.String, "price": filter.Int},
			expectError: false,
		},
		{
			name:        "field not in whitelist",
			query:       filter.Query{"invalid": "value"},
			spec:        filter.Spec{"status": filter.String, "price": filter.Int},
			expectError: true,
			errorMsg:    "Campo 'invalid' no permitido",
		},
		{
			name:        "compound field not in whitelist",
			query:       filter.Query{"invalid": map[string]any{"eq": "value"}},
			spec:        filter.Spec{"status": filter.String, "price": filter.Int},
			expectError: true,
			errorMsg:    "Campo 'invalid' no permitido",
		},
		{
			name: "operator with whitelisted fields",
			query: filter.Query{
				"and": []any{
					filter.Query{"status": "active"},
					filter.Query{"price": map[string]any{"gt": 100}},
				},
			},
			spec:        filter.Spec{"status": filter.String, "price": filter.Int},
			expectError: false,
		},
		{
			name: "operator with non-whitelisted field",
			query: filter.Query{
				"and": []any{
					filter.Query{"status": "active"},
					filter.Query{"invalid": "value"},
				},
			},
			spec:        filter.Spec{"status": filter.String, "price": filter.Int},
			expectError: true,
			errorMsg:    "Campo 'invalid' no permitido",
		},
		{
			name:        "empty whitelist allows any field",
			query:       filter.Query{"anyfield": "value"},
			spec:        filter.Spec{},
			expectError: false,
		},
		{
			name:        "nil whitelist allows any field",
			query:       filter.Query{"anyfield": "value"},
			spec:        nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.ValidateWithSpec(tt.spec)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func BenchmarkQuery_Validate_Operators(b *testing.B) {
	query := filter.Query{
		"and": []any{
			filter.Query{"status": "active"},
			filter.Query{
				"or": []any{
					filter.Query{"price": map[string]any{"gt": 100}},
					filter.Query{"discount": true},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = query.Validate()
	}
}

func BenchmarkQuery_ValidateWithSpec(b *testing.B) {
	query := filter.Query{"status": "active", "price": 100}
	spec := filter.Spec{"status": filter.String, "price": filter.Int}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = query.ValidateWithSpec(spec)
	}
}

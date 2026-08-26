package filter_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

func TestValidator_Validate(t *testing.T) {
	// Definimos un Spec de prueba
	testSpec := filter.Spec{
		"name":   filter.TypeString,
		"age":    filter.TypeInt,
		"active": filter.TypeBool,
		"score":  filter.TypeFloat,
	}

	validator := &filter.Validator{FilterSpec: testSpec}

	tests := []struct {
		name      string
		node      filter.Clause
		wantError bool
	}{
		// Casos Felices
		{
			name: "Valid String Eq",
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_EQ,
				Value:     "Juan",
			},
			wantError: false,
		},
		{
			name: "Valid Int GreaterThan",
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_GT,
				Value:     18,
			},
			wantError: false,
		},
		{
			name: "Valid Int from Float (JSON behavior)",
			// Simula que JSON envió 20.0, pero esperamos Int. Debería pasar.
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_EQ,
				Value:     20.0,
			},
			wantError: false,
		},

		// Errores de Estructura / Spec
		{
			name: "Field Not in Spec",
			node: &filter.PredicateClause{
				Field:     "hack_attempt",
				Condition: filter.CONDITON_EQ,
				Value:     "yes",
			},
			wantError: true,
		},

		// Errores de Tipo
		{
			name: "Type Mismatch: String for Int",
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_EQ,
				Value:     "twenty",
			},
			wantError: true,
		},
		{
			name: "Type Mismatch: Float with decimals for Int",
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_EQ,
				Value:     20.5,
			},
			wantError: true,
		},

		// Errores de Lógica
		{
			name: "Invalid Condition for Type (Gt on String)",
			// Según nuestra lógica simple, Gt no aplica a String (aunque en DBs a veces sí, aquí lo restringimos para el test)
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_GT,
				Value:     "Juan",
			},
			wantError: true,
		},

		// Lógica Compuesta
		{
			name: "Valid Logic AND",
			node: &filter.LogicalClause{
				Operator: filter.OPERATOR_AND,
				Children: []filter.Clause{
					&filter.PredicateClause{
						Field:     "active",
						Condition: filter.CONDITON_EQ,
						Value:     true,
					},
					&filter.PredicateClause{
						Field:     "age",
						Condition: filter.CONDITON_GTE,
						Value:     18,
					},
				},
			},
			wantError: false,
		},

		// Null Filtering
		{
			name: "Valid Eq Null on String",
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_EQ,
				Value:     nil,
			},
			wantError: false,
		},
		{
			name: "Valid Neq Null on String",
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_NEQ,
				Value:     nil,
			},
			wantError: false,
		},
		{
			name: "Valid Eq Null on Int",
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_EQ,
				Value:     nil,
			},
			wantError: false,
		},
		{
			name: "Valid Neq Null on Bool",
			node: &filter.PredicateClause{
				Field:     "active",
				Condition: filter.CONDITON_NEQ,
				Value:     nil,
			},
			wantError: false,
		},
		{
			name: "Invalid Null on Gt",
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_GT,
				Value:     nil,
			},
			wantError: true,
		},
		{
			name: "Invalid Null on Like",
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_LIKE,
				Value:     nil,
			},
			wantError: true,
		},
		{
			name: "Invalid Null on In",
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_IN,
				Value:     nil,
			},
			wantError: true,
		},

		// Condición IN (lista de valores)
		{
			name: "Valid In String list",
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_IN,
				Value:     []any{"Juan", "Ana"},
			},
			wantError: false,
		},
		{
			name: "Valid In Int list",
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_IN,
				Value:     []any{10, 20, 30},
			},
			wantError: false,
		},
		{
			name: "Invalid Eq with list",
			node: &filter.PredicateClause{
				Field:     "name",
				Condition: filter.CONDITON_EQ,
				Value:     []any{"Juan", "Ana"},
			},
			wantError: true,
		},
		{
			name: "Invalid In with wrong element type",
			node: &filter.PredicateClause{
				Field:     "age",
				Condition: filter.CONDITON_IN,
				Value:     []any{"not_a_number"},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.node)
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestValidator_NilRoot verifica que un filtro nulo (sin cláusulas) se valide
// como "sin filtro" en lugar de panic o error.
func TestValidator_NilRoot(t *testing.T) {
	validator := &filter.Validator{FilterSpec: filter.Spec{"name": filter.TypeString}}
	if err := validator.Validate(nil); err != nil {
		t.Errorf("Validate(nil) error = %v, want nil", err)
	}
}

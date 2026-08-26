package model

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
)

func TestStructToMap_Omittable(t *testing.T) {
	type Condition struct {
		Eq  graphql.Omittable[*string] `json:"eq,omitempty"`
		Neq graphql.Omittable[*string] `json:"neq,omitempty"`
	}

	t.Run("campo no enviado", func(t *testing.T) {
		c := Condition{}
		m, err := structToMap(c)
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 0 {
			t.Errorf("esperaba mapa vacío, obtuve: %v", m)
		}
	})

	t.Run("eq null explícito", func(t *testing.T) {
		c := Condition{Eq: graphql.OmittableOf[*string](nil)}
		m, err := structToMap(c)
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 1 {
			t.Errorf("esperaba 1 clave, obtuve %d: %v", len(m), m)
		}
		val, ok := m["eq"]
		if !ok {
			t.Fatal("clave 'eq' no encontrada")
		}
		if val != nil {
			t.Errorf("esperaba nil, obtuve: %v", val)
		}
	})

	t.Run("eq con valor", func(t *testing.T) {
		v := "test"
		c := Condition{Eq: graphql.OmittableOf(&v)}
		m, err := structToMap(c)
		if err != nil {
			t.Fatal(err)
		}
		val, ok := m["eq"]
		if !ok {
			t.Fatal("clave 'eq' no encontrada")
		}
		if val == nil {
			t.Fatal("esperaba valor no nulo")
		}
		if s, ok := val.(string); !ok || s != "test" {
			t.Errorf("esperaba string 'test', obtuve: %T %v", val, val)
		}
	})

	t.Run("neq null explícito", func(t *testing.T) {
		c := Condition{Neq: graphql.OmittableOf[*string](nil)}
		m, err := structToMap(c)
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 1 {
			t.Errorf("esperaba 1 clave, obtuve %d: %v", len(m), m)
		}
		val, ok := m["neq"]
		if !ok {
			t.Fatal("clave 'neq' no encontrada")
		}
		if val != nil {
			t.Errorf("esperaba nil, obtuve: %v", val)
		}
	})
}

func TestStructToMap_NestedOmittable(t *testing.T) {
	type Condition struct {
		Eq  graphql.Omittable[*string] `json:"eq,omitempty"`
		Neq graphql.Omittable[*string] `json:"neq,omitempty"`
	}

	type Filter struct {
		Cuota    *Condition `json:"cuota,omitempty"`
		Concepto *Condition `json:"concepto,omitempty"`
	}

	t.Run("solo cuota null", func(t *testing.T) {
		f := Filter{
			Cuota: &Condition{Eq: graphql.OmittableOf[*string](nil)},
		}
		m, err := structToMap(f)
		if err != nil {
			t.Fatal(err)
		}
		cuotaMap, ok := m["cuota"].(map[string]any)
		if !ok {
			t.Fatalf("esperaba map para 'cuota', obtuve: %T %v", m["cuota"], m["cuota"])
		}
		if cuotaMap["eq"] != nil {
			t.Errorf("esperaba nil en cuota.eq, obtuve: %v", cuotaMap["eq"])
		}
		if _, exists := cuotaMap["neq"]; exists {
			t.Error("neq no debería existir en cuota")
		}
		if _, exists := m["concepto"]; exists {
			t.Error("concepto no debería existir en el mapa")
		}
	})
}

func TestStructToMap_NilInput(t *testing.T) {
	m, err := structToMap(nil)
	if err != nil {
		t.Fatal(err)
	}
	if m != nil {
		t.Errorf("esperaba nil, obtuve: %v", m)
	}
}

func TestStructToMap_SliceOmittable(t *testing.T) {
	type Filter struct {
		Cuota *struct {
			Eq graphql.Omittable[*string] `json:"eq,omitempty"`
		} `json:"cuota,omitempty"`
		And []*Filter `json:"and,omitempty"`
	}

	t.Run("slice nil no aparece en mapa", func(t *testing.T) {
		f := Filter{}
		m, err := structToMap(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, exists := m["and"]; exists {
			t.Error("'and' no debería existir cuando el slice es nil")
		}
	})

	t.Run("slice vacío no aparece en mapa", func(t *testing.T) {
		f := Filter{And: []*Filter{}}
		m, err := structToMap(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, exists := m["and"]; exists {
			t.Error("'and' no debería existir cuando el slice está vacío")
		}
	})

	t.Run("slice con elementos se convierte a []any", func(t *testing.T) {
		eqVal := "test"
		child := Filter{
			Cuota: &struct {
				Eq graphql.Omittable[*string] `json:"eq,omitempty"`
			}{Eq: graphql.OmittableOf(&eqVal)},
		}
		f := Filter{And: []*Filter{&child}}
		m, err := structToMap(f)
		if err != nil {
			t.Fatal(err)
		}
		andSlice, ok := m["and"].([]any)
		if !ok {
			t.Fatalf("esperaba []any, obtuve: %T", m["and"])
		}
		if len(andSlice) != 1 {
			t.Fatalf("esperaba 1 elemento, obtuve %d", len(andSlice))
		}
		nestedMap, ok := andSlice[0].(map[string]any)
		if !ok {
			t.Fatalf("esperaba map[string]any, obtuve: %T", andSlice[0])
		}
		cuotaMap, ok := nestedMap["cuota"].(map[string]any)
		if !ok {
			t.Fatalf("esperaba map en cuota, obtuve: %T", nestedMap["cuota"])
		}
		if cuotaMap["eq"] != "test" {
			t.Errorf("esperaba 'test', obtuve: %v", cuotaMap["eq"])
		}
	})
}

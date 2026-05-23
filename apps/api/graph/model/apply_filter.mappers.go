package model

import (
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

func applyFilter[T filter.Filterable](input any) filter.Filter[T] {
	nill := *filter.NewFilter[T](nil)
	if input == nil {
		return nill
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nill
	}

	var inputMap map[string]any
	if err := json.Unmarshal(jsonBytes, &inputMap); err != nil {
		return nill
	}

	return *filter.NewFilter[T](inputMap)
}

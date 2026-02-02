package utils

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/stretchr/testify/assert"
)

type ValidableTest struct {
	Name    string
	Input   core.Validable
	WantErr bool
	ErrMsg  string
}

func TestValidables(t *testing.T, tests []ValidableTest, parallel bool) {
	for _, tt := range tests {
		test := tt
		t.Run(test.Name, func(t *testing.T) {
			if parallel {
				t.Parallel()
			}
			test.Run(t)
		})
	}
}

func (v ValidableTest) Run(t *testing.T) {

	if v.WantErr && v.ErrMsg == "" {
		t.Fatal("Debe proporcionar un mensaje de error si WantErr es true")
	}

	err := v.Input.Validate()

	if err != nil {
		t.Log(err.Error())
	}

	if v.WantErr {
		assert.Error(t, err, "El test '%s' debería haber fallado", v.Name)
		assert.Contains(
			t,
			err.Error(),
			v.ErrMsg,
			"Se esperaba que el error contenga el mensaje: %s\nPero en su lugar se obtuvo: %s",
			v.ErrMsg,
			err.Error(),
		)
	} else {
		assert.NoError(t, err)
	}
}

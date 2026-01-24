package query_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/testing/utils"
	"github.com/stretchr/testify/assert"
)

func TestObtenerGastosSegunCuota(t *testing.T) {

	db, mock := utils.SetupMockTestDB(t)

	ctx := context.New(t.Context(), db)
	uc := query.NewObtenerGastosSegunCuota()
	input := query.ObtenerGastosSegunCuotaDTO{
		CuotaID: "1",
	}

	// Verificar que la cuota existe
	mock.ExpectQuery(
		fmt.Sprintf(
			"^SELECT %s FROM %s WHERE id = \\? LIMIT 1$",
			utils.GetMockTableRegex("id"),
			utils.GetMockTableRegex("cuotas"),
		),
	).WithArgs(input.CuotaID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(input.CuotaID))
	// SELECT COUNT
	mock.ExpectQuery(
		fmt.Sprintf(
			"^SELECT (COUNT|count)\\(\\*\\) FROM %s WHERE cuota = \\?$",
			utils.GetMockTableRegex("gastos"),
		),
	).WithArgs(input.CuotaID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// SELECT gastos
	mock.ExpectQuery(
		fmt.Sprintf(
			"^SELECT \\* FROM %s WHERE cuota = \\? LIMIT %d$",
			utils.GetMockTableRegex("gastos"),
			common.DEFAULT_PAGE_SIZE,
		),
	).
		WithArgs(input.CuotaID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	gastos, err := uc.Exec(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, gastos)

}

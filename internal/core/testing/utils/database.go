package utils

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/pagos"
	"github.com/Sanaruca/condominio/internal/villas"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupMockTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creando mock db: %v", err)
	}

	mock.ExpectQuery("select sqlite_version()").
		WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.35.5"))

	newLogger := logger.New(
		log.New(os.Stdout, "\x1b[1;45m[GORM]\t\x1b[0m", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info, // nivel: Info muestra SQL
			Colorful:      true,
			// ... other config
		},
	)

	gdb, err := gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		t.Fatalf("error abriendo gorm con mock: %v", err)
	}
	return gdb, mock
}

func SetupInMemoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la BD en memoria: %v", err)
	}

	// Migramos las tablas necesarias
	err = db.AutoMigrate(
		&villas.Villa{},
		&administracion.Cuota{},
		&administracion.Proveedor{},
		&administracion.IGasto{},
		&pagos.Destino{},
	)

	if err != nil {
		t.Fatalf("no se pudo migrar: %v", err)
	}

	runSQLFile(t, db, "../../../../sql/views/pagos.sql")
	runSQLFile(t, db, "../../../../sql/views/deudas.sql")
	runSQLFile(t, db, "../../../../sql/views/gastos.sql")

	return db
}

func runSQLFile(
	t *testing.T,
	db *gorm.DB,
	file string,
) {
	sql_bytes, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("no se pudo leer el archivo: %v", err)
	}

	if err := db.Exec(string(sql_bytes)).Error; err != nil {
		t.Fatalf("no se pudo ejecutar el archivo: %v", err)
	}
}

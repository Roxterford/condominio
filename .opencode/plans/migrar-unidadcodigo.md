# Plan: Migrar string → UnidadCodigo

## Archivos a modificar (11 archivos)

### 1. `apps/api/internal/administracion/models/deuda/deuda.go`
- Agregar import: `"github.com/Sanaruca/condominio/internal/unidades/models/unidad"`
- L19: `unidad string` → `unidad unidad.UnidadCodigo`
- L27: `func (d *Deuda) Unidad() string` → `func (d *Deuda) Unidad() unidad.UnidadCodigo`

### 2. `apps/api/internal/administracion/models/deuda/deuda.factory.go`
- Agregar import: `"github.com/Sanaruca/condominio/internal/unidades/models/unidad"`
- L28: `unidad string` → `unidad unidad.UnidadCodigo`
- L31: `if unidad == ""` → `if unidad == unidad.UnidadCodigo("")` (o se deja igual porque `UnidadCodigo` es `string` subyacente)
- L41: `unidad: unidad` (se deja igual)
- L51: `unidad string` → `unidad unidad.UnidadCodigo`
- L59: `unidad: unidad` (se deja igual)

### 3. `apps/api/internal/administracion/models/deuda/deuda.repository.go`
- L15: `unidadCodigo string` → `unidadCodigo unidad.UnidadCodigo`

### 4. `apps/api/internal/unidades/models/unidad/unidad.repository.go`
- L21: `codigo string` → `codigo UnidadCodigo`
- L28: `[]string` → `[]UnidadCodigo`
- L29: `unidad string` → `unidad UnidadCodigo`

### 5. `apps/api/internal/unidades/models/unidad/unidad.factory.go`
- L27: `codigo string` → `codigo UnidadCodigo`

### 6. `apps/api/internal/unidades/adapters/gorm/unidad.repository.impl.go`
- L44: `codigo string` → `codigo unidad.UnidadCodigo`
- L114: `[]string` → `[]unidad.UnidadCodigo`
- L120: `make([]string, ...)` → `make([]unidad.UnidadCodigo, ...)`
- L122: `codigos[i] = u.Codigo` → `codigos[i] = unidad.UnidadCodigo(u.Codigo)`
- L130: `unidadCodigo string` → `unidadCodigo unidad.UnidadCodigo`
- L133: `"codigo = ?", unidadCodigo` → `"codigo = ?", string(unidadCodigo)`

### 7. `apps/api/internal/administracion/adapters/gorm/deuda.repository.impl.go`
- L49: `unidadCodigo string` → `unidadCodigo unidad.UnidadCodigo`
- L52: `unidadCodigo` → `string(unidadCodigo)`
- L126: `unidad string` → `unidad unidad.UnidadCodigo`
- L130: `unidad,` → `string(unidad),`
- L156: `_deuda.UnidadID,` → `unidad.UnidadCodigo(_deuda.UnidadID),`
- L170: `deudaEntity.Unidad(),` → `string(deudaEntity.Unidad()),`

### 8. `apps/api/internal/administracion/adapters/gorm/types.go`
- L208: `t.UnidadID,` → `unidad.UnidadCodigo(t.UnidadID),`

### 9. `apps/api/internal/pagos/adapters/gorm/pago.table.go`
- L45: `p.Unidad,` → `unidad.UnidadCodigo(p.Unidad),`
- L89: `p.Unidad(),` → `string(p.Unidad()),`
- Agregar import: `"github.com/Sanaruca/condominio/internal/unidades/models/unidad"`

### 10. `apps/api/graph/registrar_pago.resolvers.go`
- L25: `unidad.UnidadID(input.Unidad)` → `unidad.UnidadCodigo(input.Unidad)`

### 11. `apps/api/graph/model/models_mappers.go`
- L26: `p.Unidad(),` → `string(p.Unidad()),`
- L59: `d.Unidad(),` → `string(d.Unidad()),`
- L93: `unidad.Codigo(),` → `string(unidad.Codigo()),`

### 12. `apps/api/internal/pagos/app/command/registrar_pago.usecase_test.go`
- L27: `unidad.UnidadID("1")` → `unidad.UnidadCodigo("1")`
- L40: `unidad.UnidadID("0")` → `unidad.UnidadCodigo("0")`
- L48: `unidad.UnidadID("1")` → `unidad.UnidadCodigo("1")`
- L56: `unidad.UnidadID("1")` → `unidad.UnidadCodigo("1")`
- L64: `unidad.UnidadID("1")` → `unidad.UnidadCodigo("1")`
- L76: `unidad.UnidadID("1")` → `unidad.UnidadCodigo("1")`
- L84: `unidad.UnidadID("1")` → `unidad.UnidadCodigo("1")`

## Post-ejecución
```bash
cd apps/api && go build ./...
```

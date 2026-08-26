# 📊 Backlog de Producto — Condominio

> Documento de propiedad del producto. Prioriza por **valor de negocio y velocidad para
> conseguir clientes y facturar**, no por deuda técnica. Los TODOs de código viven en
> `readme.md` y son un artefacto distinto (técnico), no de producto.

## 🎯 Tesis del producto

El producto le hace *obvia* al administrador la gestión de cobro y gastos de su
condominio, y le da los reportes para **cobrar y rendir cuentas**.

No se construye para un solo cliente piloto: se diseña escalable (procesos del
condominio donde residimos + investigación de otros), con una "ventana" amplia pero
enfocada — sin caer en un producto imposible.

## ✅ Definición de éxito del primer milestone (la demo que vende)

> Un admin carga su edificio (unidades + propietarios), genera la cuota del mes, ve
> deudas y recaudación, registra un pago y ve la deuda bajar y la recaudación subir —
> **en una sesión, con montos correctos**.

Si eso funciona en la demo, se cierran clientes. Todo lo demás es ruido hasta ahí.

## 🚀 Prioridades

### P0 — La rebanada vertical que VENDE (demo de cobro de punta a punta)

| # | Item | Por qué es prioritario (negocio) | Notas técnicas |
|---|------|----------------------------------|----------------|
| A | **Onboarding del condominio** (crear Unidades + Sujetos/propietarios) | Sin datos del edificio no hay nada que cobrar → **bloqueador #1 de la demo** | ⚠️ Hoy NO existe mutación ni UI para crear unidades/propietarios (solo listar/detalle). Hay que construirla. |
| B | **Bucle de cobro**: Cuota → deudas auto por unidad → Pago → abono → Recaudación | Es el "wow" que justifica el sistema | Validar de punta a punta en la demo, no solo que compila |
| C | **Reportes que el admin enseña**: deudas por unidad, recaudación del período, estado de cuenta | Es lo que el admin lleva a la asamblea para cobrar | Dashboard existe; verificar que trae datos reales |
| D | **Corregir bug del filtro `in`** (`apps/api/internal/core/common/filter/validator.go:140`) | Bloquea ver listas/deudas en pantalla → **bloqueador de venta** | No es "técnica menor": sin listas no hay demo |
| E | **Exactitud multi-moneda/tasa** (VED/VES/USD) | Un monto mal en la demo destruye la credibilidad y la venta | *Gate de calidad P0*, no P2 |

### P1 — Transparencia (retiene y justifica adoptar)

| # | Item | Por qué |
|---|------|---------|
| F | **Gastos + Proveedores** en la demo (registrar gasto, asociar proveedor, ver egresos) | El admin debe mostrar a propietarios dónde fue el dinero |
| G | **Reporte gastos vs recaudación** | Lo que el admin lleva a la asamblea para justificar el sistema |
| H | **Login de operador/admin mínimo** (1 admin seed) | Suficiente para la demo; roles completos después |

### P2 — Escalar / higiene (después de los primeros clientes)

- Límite de gastos por cuota, N+1, tests stub, CI / branch protection.
- **Portal de propietario self-service** — *solo si* se valida demanda.
- **Modelo de cobro real** (freemium condicionado a hosting, etc.) — *tras* tener
  clientes y datos de uso, no antes.

## 🔭 Fuera de alcance (por ahora)

- Roles finos de usuario (admin/propietario/etc.) más allá de un operador seed.
- Internacionalización fuera de español.
- Cualquier feature que no alimente directamente el bucle de cobro + gastos de la demo.

## 📌 Principio rector

> Antes de pulir, agregar o "mejorar", preguntarse: *¿esto acelera que un administrador
> vea su edificio cobrado y recaudado en la demo?* Si no, es P2.

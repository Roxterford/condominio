# Changelog

## [Unreleased]

### feat
- **transacciones**: unificar pagos/gastos en modelo `TransaccionFinanciera/Movimiento` con factory, repositorio GORM, casos de uso y resolvers GraphQL
- **schema**: migrar Prisma de IPago/IGasto a ITransaccion/IMovimiento con enums TipoDeMovimiento/RolDelMovimiento
- **views**: actualizar vistas SQL (pagos, gastos, unidades_info, tasas_de_cambio, recaudacion) para usar internal_transacciones/internal_movimientos
- **graphql**: agregar mutation `registrarTransaccion`, query `obtenerMovimientos`, tipos Transaccion/Movimiento, inputs y filtros
- **factory**: implementar `NuevoPago`, `NuevoGasto`, `NuevoReembolso` con auditoría de `sum(movimientos) == monto_total`
- **wiring**: integrar TransaccionFactory, GORMTransaccionRepository y TransaccionService en server.go

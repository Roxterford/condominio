DROP VIEW IF EXISTS gastos;
CREATE VIEW gastos AS
SELECT
  o.id AS operacion,
  COALESCE(bo.transaccion_id, '') AS transaccion,
  o.tipo,
  o.rol,
  o.monto,
  o.cuota,
  o.unidad_codigo,
  o.proveedor_id,
  o.concepto,
  o.fecha,
  o.monto AS monto_total,
  o.moneda,
  o.metodo,
  o.tasa
FROM operaciones o
LEFT JOIN transaccion_operaciones bo ON bo.operacion_id = o.id
WHERE o.tipo = 'DEBITO';

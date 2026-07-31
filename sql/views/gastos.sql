DROP VIEW IF EXISTS gastos;
CREATE VIEW gastos AS
SELECT
  m.id,
  t.concepto,
  m.proveedor_id AS proveedor,
  t.cuota_id AS cuota,
  m.monto,
  t.moneda,
  t.tasa,
  CAST(
    CASE
      WHEN t.moneda <> 'USD' THEN Cast(m.monto * 1.0 / t.tasa * 100 AS INTEGER)
      ELSE m.monto
    END AS INTEGER
  ) AS total,
  t.fecha,
  t.concepto AS descripcion,
  t.registro,
  t.registrado_por
FROM internal_movimientos m
JOIN internal_transacciones t ON t.id = m.transaccion_id
WHERE m.tipo = 'DEBITO' AND (m.rol = 'PROVEEDOR' OR m.rol = 'CONDOMINIO');

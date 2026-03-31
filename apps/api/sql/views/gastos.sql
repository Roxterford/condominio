DROP VIEW IF EXISTS gastos;
CREATE VIEW gastos AS
SELECT
  g.id,
  g.proveedor,
  g.cuota,
  g.monto,
  g.moneda,
  g.tasa,
  CAST(
    CASE
      WHEN g.moneda <> 'USD' THEN Cast(g.monto * 1.0 / g.tasa * 100 AS INTEGER)
      ELSE g.monto
    END AS INTEGER
  ) AS total,
  g.fecha,
  g.descripcion,
  g.registro,
  g.registrado_por
FROM
  internal_gastos g
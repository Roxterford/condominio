DROP VIEW IF EXISTS operaciones;
CREATE VIEW operaciones AS
SELECT
  o.*,
  u.id AS unidad_id,
  CASE
    WHEN o.moneda = "USD" THEN o.monto 
    ELSE o.monto / o.tasa
  END AS total
FROM
  internal_operaciones o
  LEFT JOIN unidades u ON u.codigo = o.unidad_codigo;

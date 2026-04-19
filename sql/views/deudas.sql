DROP VIEW IF EXISTS deudas;
CREATE VIEW deudas AS
SELECT
  d.id,
  d.unidad,
  d.cuota,
  c.monto AS monto,
  c.monto - COALESCE(SUM(dp.destinado), 0) AS deuda,
  CASE
    WHEN COALESCE(SUM(dp.destinado), 0) = c.monto THEN 'PAGADA'
    WHEN COALESCE(SUM(dp.destinado), 0) = 0 THEN 'PENDIENTE'
    WHEN COALESCE(SUM(dp.destinado), 0) < c.monto THEN 'ABONADA'
  END AS estado,
  d.registro,
  d.actualizacion
FROM
  internal_deudas d
  JOIN cuotas c ON d.cuota = c.id
  LEFT JOIN destino_de_pagos dp ON dp.deuda = d.id
GROUP BY
  d.id,
  c.monto;
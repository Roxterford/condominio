DROP VIEW IF EXISTS deudas;
CREATE VIEW deudas AS
SELECT
  d.id,
  u.id AS unidad_id,
  u.codigo AS unidad_codigo,
  d.cuota,
  c.monto AS monto,
  c.monto - COALESCE(SUM(dp.destinado), 0) AS deuda,
  CASE
    WHEN COALESCE(SUM(dp.destinado), 0) = c.monto THEN 'SALDADA'
    WHEN COALESCE(SUM(dp.destinado), 0) = 0 THEN 'PENDIENTE'
    WHEN COALESCE(SUM(dp.destinado), 0) < c.monto THEN 'ABONADA'
  END AS estado,
  d.registro,
  d.actualizacion
FROM
  internal_deudas d
  JOIN cuotas c ON d.cuota = c.id
  LEFT JOIN destino_de_pagos dp ON dp.deuda = d.id
  LEFT JOIN unidades u ON u.codigo = d.unidad
GROUP BY
  d.id,
  c.monto;
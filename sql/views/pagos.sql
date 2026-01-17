-- DROP VIEW IF EXISTS pagos;
CREATE VIEW pagos AS
WITH
  pagos_con_total AS (
    SELECT
      p.*,
      CASE
        WHEN p.moneda <> 'USD' THEN Cast(p.monto * 1.0 / p.tasa * 100 AS INTEGER)
        ELSE p.monto
      END AS total
    FROM
      internal_pagos p
  )
SELECT
  pt.id,
  pt.villa,
  pt.fecha,
  pt.metodo,
  pt.monto,
  pt.referencia,
  pt.moneda,
  pt.tasa,
  pt.total,
  COALESCE(Sum(d.destinado), 0) AS destinado,
  pt.total - COALESCE(Sum(d.destinado), 0) AS cuenta,
  pt.registro,
  pt.registrado_por,
  pt.actualizacion,
  pt.actualizado_por
FROM
  pagos_con_total pt
  LEFT JOIN destino_de_pagos d ON pt.id = d.pago
GROUP BY
  pt.id,
  pt.total;
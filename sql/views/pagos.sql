DROP VIEW IF EXISTS pagos;
CREATE VIEW pagos AS
WITH
  movimientos_pago AS (
    SELECT
      m.id,
      m.transaccion_id,
      m.monto,
      m.unidad_codigo AS unidad,
      t.fecha,
      t.concepto,
      t.monto_total,
      t.moneda,
      t.tasa,
      t.registro,
      t.registrado_por,
      CASE
        WHEN t.moneda <> 'USD' THEN Cast(t.monto_total * 1.0 / t.tasa * 100 AS INTEGER)
        ELSE t.monto_total
      END AS total
    FROM internal_movimientos m
    JOIN internal_transacciones t ON t.id = m.transaccion_id
    WHERE m.tipo = 'CREDITO' AND m.rol = 'UNIDAD'
  )
SELECT
  mp.id,
  mp.unidad,
  mp.fecha,
  mp.concepto,
  mp.monto,
  mp.monto_total,
  mp.moneda,
  mp.tasa,
  mp.total,
  COALESCE(Sum(d.destinado), 0) AS destinado,
  mp.total - COALESCE(Sum(d.destinado), 0) AS cuenta,
  mp.registro,
  mp.registrado_por
FROM
  movimientos_pago mp
  LEFT JOIN destino_de_pagos d ON d.movimiento = mp.id
GROUP BY
  mp.id, mp.unidad, mp.fecha, mp.concepto, mp.monto, mp.monto_total,
  mp.moneda, mp.tasa, mp.total, mp.registro, mp.registrado_por;

DROP VIEW IF EXISTS recaudacion;
CREATE VIEW recaudacion AS
SELECT 
  c.id AS cuota,
  c.tipo,
  c.monto,
  c.mes,
  c.anio,
  (SELECT COUNT(*) FROM villas) AS villas,
  COUNT(DISTINCT d.id) AS villas_aplicadas,
  (SELECT COUNT(DISTINCT d2.villa) 
   FROM internal_deudas d2 
   JOIN destino_de_pagos dp2 ON dp2.deuda = d2.id 
   WHERE d2.cuota = c.id) AS villas_solventes,
  (SELECT COUNT(DISTINCT d2.villa) 
   FROM internal_deudas d2 
   LEFT JOIN destino_de_pagos dp2 ON dp2.deuda = d2.id 
   WHERE d2.cuota = c.id AND (dp2.destinado IS NULL OR dp2.destinado = 0)) AS villas_pendientes,
  (CAST(c.monto AS INTEGER) * COUNT(DISTINCT d.id)) AS total_estimado,
  COALESCE(SUM(dp.destinado), 0) AS recaudado,
  (CAST(c.monto AS INTEGER) * COUNT(DISTINCT d.id)) - COALESCE(SUM(dp.destinado), 0) AS pendiente,
  COUNT(DISTINCT dp.pago) AS pagos_asociados
FROM cuotas c
LEFT JOIN internal_deudas d ON d.cuota = c.id
LEFT JOIN destino_de_pagos dp ON dp.deuda = d.id
GROUP BY c.id, c.monto;
DROP VIEW IF EXISTS tasas_de_cambio;
CREATE VIEW tasas_de_cambio AS
SELECT
  moneda,
  tasa,
  fecha,
  'pagos' AS origen,
  id AS origen_id
FROM
  internal_pagos
UNION ALL
SELECT
  moneda,
  tasa,
  'gastos' AS origen,
  id AS origen_id,
  fecha
FROM
  internal_gastos;

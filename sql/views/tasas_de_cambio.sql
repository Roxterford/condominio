DROP VIEW IF EXISTS tasas_de_cambio;
CREATE VIEW tasas_de_cambio AS
SELECT
  DATE(fecha) AS fecha,
  tasa,
  'VED' AS moneda,
  COUNT(*) AS cantidad_operaciones
FROM
  operaciones
WHERE tasa > 0
GROUP BY DATE(fecha), tasa
ORDER BY fecha DESC;
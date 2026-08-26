DROP VIEW IF EXISTS tasas_de_cambio;
CREATE VIEW tasas_de_cambio AS
SELECT
  moneda,
  tasa,
  fecha,
  'operaciones' AS origen,
  id AS origen_id
FROM
  operaciones;

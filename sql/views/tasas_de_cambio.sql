DROP VIEW IF EXISTS tasas_de_cambio;
CREATE VIEW tasas_de_cambio AS
SELECT
  moneda,
  tasa,
  fecha,
  'transacciones' AS origen,
  id AS origen_id
FROM
  internal_transacciones;

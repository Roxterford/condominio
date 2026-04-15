DROP VIEW IF EXISTS villas_stats;
CREATE VIEW villas_stats AS
SELECT 
  (SELECT COUNT(*) FROM villas) AS total_villas,
  (SELECT COUNT(DISTINCT villa) FROM deudas WHERE deuda > 0) AS villas_con_pendientes,
  (SELECT SUM(deuda) FROM deudas) AS total_pendiente;
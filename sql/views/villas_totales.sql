DROP VIEW IF EXISTS villas_totales;
CREATE VIEW villas_totales AS
SELECT 
  (SELECT COUNT(*) FROM villas) AS total_villas,
  (SELECT COUNT(*) FROM villas WHERE estado = 'ACTIVA') AS villas_activas,
  (SELECT COUNT(*) FROM villas WHERE estado = 'INHABITADA') AS villas_inhabitadas,
  (SELECT COUNT(*) FROM villas WHERE estado = 'EXENTA') AS villas_exentas,
  (SELECT COUNT(*) FROM villas WHERE estado = 'EN_LITIGIO') AS villas_en_litigio,
  (SELECT COUNT(*) FROM villas WHERE estado = 'SUSPENDIDA') AS villas_suspendidas,
  (SELECT COUNT(*) FROM villas WHERE estado = 'PREVENTA') AS villas_preventa,
  (SELECT COUNT(DISTINCT villa) FROM deudas WHERE deuda > 0) AS villas_con_pendientes,
  (SELECT COUNT(DISTINCT villa) FROM deudas WHERE deuda = 0 OR deuda IS NULL) AS villas_solventes,
  (SELECT SUM(deuda) FROM deudas) AS total_pendiente,
  (SELECT SUM(c.monto) FROM internal_deudas id JOIN cuotas c ON id.cuota = c.id) AS total_asignado;
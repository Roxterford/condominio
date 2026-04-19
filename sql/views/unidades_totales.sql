DROP VIEW IF EXISTS unidades_totales;
CREATE VIEW unidades_totales AS
SELECT 
  (SELECT COUNT(*) FROM unidades) AS total_unidades,
  (SELECT COUNT(*) FROM unidades WHERE estado = 'ACTIVA') AS unidades_activas,
  (SELECT COUNT(*) FROM unidades WHERE estado = 'INHABITADA') AS unidades_inhabitadas,
  (SELECT COUNT(*) FROM unidades WHERE estado = 'EXENTA') AS unidades_exentas,
  (SELECT COUNT(*) FROM unidades WHERE estado = 'EN_LITIGIO') AS unidades_en_litigio,
  (SELECT COUNT(*) FROM unidades WHERE estado = 'SUSPENDIDA') AS unidades_suspendidas,
  (SELECT COUNT(*) FROM unidades WHERE estado = 'PREVENTA') AS unidades_preventa,
  (SELECT COUNT(DISTINCT unidad) FROM deudas WHERE deuda > 0) AS unidades_con_pendientes,
  (SELECT COUNT(DISTINCT unidad) FROM deudas WHERE deuda = 0 OR deuda IS NULL) AS unidades_solventes,
  (SELECT SUM(deuda) FROM deudas) AS total_pendiente,
  (SELECT SUM(c.monto) FROM internal_deudas id JOIN cuotas c ON id.cuota = c.id) AS total_asignado;
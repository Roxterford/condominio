DROP VIEW IF EXISTS unidades_info;
CREATE VIEW unidades_info AS
SELECT 
  u.id,
  u.codigo,
  u.estado,
  u.contacto,
  u.descripcion,
  COALESCE(SUM(d.deuda), 0) AS deuda_total,
  CASE 
    WHEN COALESCE(SUM(d.deuda), 0) = 0 THEN 'SOLVENTE'
    WHEN COALESCE(SUM(d.deuda), 0) < (SELECT SUM(c.monto) FROM internal_deudas id JOIN cuotas c ON id.cuota = c.id WHERE id.unidad = u.codigo) THEN 'ABONADA'
    ELSE 'PENDIENTE'
  END AS estado_cuenta,
  COUNT(d.id) AS cuotas_pendientes
FROM unidades u
LEFT JOIN deudas d ON u.codigo = d.unidad
GROUP BY u.codigo;
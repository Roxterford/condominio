DROP VIEW IF EXISTS unidades_resumen;
CREATE VIEW unidades_resumen AS
SELECT 
  v.codigo AS unidad_codigo,
  SUM(d.deuda) AS deuda_total,
  CASE 
    WHEN SUM(d.deuda) = 0 THEN 'SOLVENTE'
    WHEN SUM(d.deuda) < (SELECT SUM(c.monto) FROM internal_deudas id JOIN cuotas c ON id.cuota = c.id WHERE id.unidad = v.codigo) THEN 'ABONADA'
    ELSE 'PENDIENTE'
  END AS estado_cuenta,
  COUNT(d.id) AS cuotas_pendientes
FROM unidades v 
LEFT JOIN deudas d ON v.codigo = d.unidad 
GROUP BY v.codigo;
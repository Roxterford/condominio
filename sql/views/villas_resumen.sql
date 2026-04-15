DROP VIEW IF EXISTS villas_resumen;
CREATE VIEW villas_resumen AS
SELECT 
  v.numero AS villa_numero,
  SUM(d.deuda) AS deuda_total,
  CASE 
    WHEN SUM(d.deuda) = 0 THEN 'SOLVENTE'
    WHEN SUM(d.deuda) < (SELECT SUM(c.monto) FROM internal_deudas id JOIN cuotas c ON id.cuota = c.id WHERE id.villa = v.numero) THEN 'ABONADA'
    ELSE 'PENDIENTE'
  END AS estado_cuenta,
  COUNT(d.id) AS cuotas_pendientes
FROM villas v 
LEFT JOIN deudas d ON v.numero = d.villa 
GROUP BY v.numero;
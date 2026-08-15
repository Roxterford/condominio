DROP VIEW IF EXISTS unidades_info;

CREATE VIEW unidades_info AS
SELECT
    u.id,
    u.codigo,
    u.estado,
    u.contacto,
    u.titular_primario,
    u.descripcion,
    COALESCE(SUM(d.deuda), 0) AS deuda_total,
    CASE
        WHEN COALESCE(SUM(d.deuda), 0) = 0 THEN 'SOLVENTE'
        WHEN COALESCE(SUM(d.deuda), 0) < (
            SELECT SUM(c.monto)
            FROM internal_deudas id
                JOIN cuotas c ON id.cuota = c.id
            WHERE
                id.unidad = u.codigo
        ) THEN 'ABONADA'
        ELSE 'PENDIENTE'
    END AS estado_cuenta,
    COUNT(d.id) AS cuotas_pendientes,

-- Cálculo de la cuenta (Total Pagado - Total Destinado)
(
    SELECT COALESCE(SUM(op.monto), 0)
    FROM operaciones op
    WHERE
        op.unidad_codigo = u.codigo
        AND op.tipo = 'CREDITO'
        AND op.rol = 'UNIDAD'
) - (
    SELECT COALESCE(SUM(dp.destinado), 0)
    FROM
        destino_de_pagos dp
        JOIN operaciones op ON dp.operacion = op.id
    WHERE
        op.unidad_codigo = u.codigo
        AND op.tipo = 'CREDITO'
        AND op.rol = 'UNIDAD'
) AS cuenta
FROM unidades u
    LEFT JOIN deudas d ON u.codigo = d.unidad
GROUP BY
    u.codigo;

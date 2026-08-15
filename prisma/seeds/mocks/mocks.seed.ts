import {
  EstadoDeProyecto,
  MetodoDeOperacion,
  Moneda,
  RolDelMovimiento,
  TipoDeCuota,
  TipoDeMovimiento,
} from "../../../generated/prisma/client";
import { prisma } from "../../client";

export async function main() {
  console.log("🌱 Seeding with mock data database...");

  await prisma.sujeto.createMany({
      data: [
        {
          id: "juan",
          tipo: "PERSONA_NATURAL",
          documento_identidad: "V-12345678",
          nombres: "Juan Andres",
          apellidos: "Ramirez Rodriguez",
          email: "juan@email.com",
          telefono: "+584123456789",
        },
        {
          id: "roxterford",
          tipo: "ENTE_JURIDICO",
          documento_identidad: "J-12345678",
          razon_social: "Roxterford",
          email: "roxter@roxterford.com",
          telefono: "+584123456789",
        },
        {
          id: "santiago",
          tipo: "PERSONA_NATURAL",
          documento_identidad: "V-12345679",
          nombres: "Santiago Mariño",
          apellidos: "Mariño Rodriguez",
          email: "santiago@email.com",
          telefono: "+584128888888",
        },
      ],
    });

    await prisma.sujeto.update({
      where: { id: "roxterford" }, data: {
        representante: "santiago"
      }
    })

    await prisma.unidad.createMany({
      data: Array(500)
        .fill(null)
        .map((_, i) => ({
          id: "u" + (i + 1).toString(),
          codigo: "villa-" + (i + 1).toString(),
          estado: "ACTIVA",
        })),
    });

    await prisma.titularidad.createMany({
      data: [
        {
          titular: "juan",
          unidad: "u1",
        },
        {
          titular: "roxterford",
          unidad: "u2",
        },
        {
          titular: "roxterford",
          unidad: "u3",
        },
        {
          titular: "roxterford",
          unidad: "u4",
        },
      ],
    });

    await prisma.unidad.update({ where: { id: "u1" }, data: { contacto: "juan", titular_primario: "juan" } })
    await prisma.unidad.updateMany({ where: { id: { in: ["u2", "u3", "u4"] } }, data: { contacto: "santiago", titular_primario: "roxterford" } })

    await prisma.usuario.create({
      data: {
        id: "tester",
        email: "tester@example.com",
        password:
          "$2a$12$TWaUL3tJEuMNUfC7uiiAjelPshhEWyBePLxVzs34LWdi1TnpJN0ZO", // password
      },
    });

    const proveedor = await prisma.proveedor.create({
      data: {
        id: "pvdr0",
        nombre: "Proveedor 0",
        rif: "J-123456789",
        email: "proveedor0@example.com",
        telefono: "04121234567",
      },
    });

    const gastoOperacion = await prisma.iOperacion.create({
      data: {
        id: "mg0",
        concepto: "Some",
        monto: 25_00,
        moneda: Moneda.VED,
        metodo: MetodoDeOperacion.EFECTIVO,
        tasa: 12_50,
        tipo: TipoDeMovimiento.DEBITO,
        rol: RolDelMovimiento.PROVEEDOR,
        proveedor_id: proveedor.id,
        registrado_por: "tester",
      },
    });

    const cuotas = await prisma.cuota.createManyAndReturn({
      data: [
        {
          id: "c0",
          monto: 8_75,
          mes: 1,
          anio: 2026,
          tipo: TipoDeCuota.REGULAR,
          registrado_por: "tester",
          actualizado_por: "tester",
        },
        {
          id: "c1",
          monto: 9_22,
          mes: 2,
          anio: 2026,
          tipo: TipoDeCuota.REGULAR,
          registrado_por: "tester",
          actualizado_por: "tester",
        },
        {
          id: "c3",
          monto: 8_94,
          mes: 3,
          anio: 2026,
          tipo: TipoDeCuota.ESPECIAL,
          registrado_por: "tester",
          actualizado_por: "tester",
        },

        {
          id: "c4",
          monto: 8_75,
          mes: 4,
          anio: 2026,
          tipo: TipoDeCuota.ESPECIAL,
          registrado_por: "tester",
          actualizado_por: "tester",
        },
      ],
    });

    await prisma.proyecto.createMany({
      data: cuotas
        .filter((cuota) => cuota.tipo === TipoDeCuota.ESPECIAL)
        .map((cuota) => ({
          titulo: `Proyecto para cuota ${cuota.id}`,
          cuota: cuota.id,
          descripcion: `Proyecto para cuota ${cuota.id}`,
          registrado_por: "tester",
          estado: EstadoDeProyecto.ACTIVO,
          actualizado_por: "tester",
          justificacion: `Proyecto para cuota ${cuota.id}`,
          fecha_limite: new Date(cuota.anio, cuota.mes - 1, 1),
        })),
    });

    for (const [i, c] of cuotas.entries()) {
      await prisma.iDeuda.createMany({
        data: Array(500)
          .fill(null)
          .map((_, v) => {
            // console.log(`d${i + v + 1}[${c.id}]v[${v + 1}]`);
            return {
              id: `d${i + v + 1}[${c.id}]v[${v + 1}]`,
              unidad: "villa-" + (v + 1).toString(),
              cuota: c.id,
            };
          }),
      });
    }

    const pagoOperacion = await prisma.iOperacion.create({
      data: {
        id: "mp0",
        concepto: "Pago de prueba",
        monto: 8134_50,
        moneda: Moneda.VED,
        metodo: MetodoDeOperacion.EFECTIVO,
        tasa: 325_38,
        tipo: TipoDeMovimiento.CREDITO,
        rol: RolDelMovimiento.UNIDAD,
        unidad_codigo: "villa-500",
        registrado_por: "tester",
      },
    });

    await prisma.destinoDePago.create({
      data: { deuda: "d500[c0]v[500]", operacion: pagoOperacion.id, destinado: 8_75 },
    });

    console.log("✅ Seeding completed.");
}


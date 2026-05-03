import {
  EstadoDeProyecto,
  MetodoDePago,
  Moneda,
  TipoDeCuota,
} from "../../../generated/prisma/client";
import { prisma } from "../../client";

export async function main() {
  console.log("🌱 Seeding with mock data database...");

  await prisma.$transaction(async (tx) => {
    await tx.sujeto.createMany({
      data: [
        {
          id: "sp0",
          tipo: "PERSONA_NATURAL",
          documento_identidad: "V-12345678",
          nombres: "Pepe Andres",
          apellidos: "Ramirez Rodriguez",
          email: "persona@email.com",
          telefono: "+584123456789",
        },
        {
          id: "se0",
          tipo: "ENTE_JURIDICO",
          documento_identidad: "J-12345678",
          razon_social: "Empresa SA",
          email: "empresa@empresa.com",
          telefono: "+584123456789",
        },
        {
          id: "sp1",
          tipo: "PERSONA_NATURAL",
          documento_identidad: "V-12345679",
          nombres: "Santiago Mariño",
          apellidos: "Mariño Rodriguez",
          email: "santiago@email.com",
          telefono: "+584128888888",
        },
      ],
    });

    await tx.sujeto.update({
      where: { id: "se0" }, data: {
        representante: "sp1"
      }
    })

    await tx.unidad.createMany({
      data: Array(500)
        .fill(null)
        .map((_, i) => ({
          id: "u" + (i + 1).toString(),
          codigo: "villa-" + (i + 1).toString(),
          estado: "ACTIVA",
        })),
    });

    await tx.titularidad.createMany({
      data: [
        {
          propietario: "sp0",
          unidad: "u1",
        },
        {
          propietario: "se0",
          unidad: "u2",
        },
        {
          propietario: "se0",
          unidad: "u3",
        },
        {
          propietario: "se0",
          unidad: "u4",
        },
      ],
    });

    await tx.unidad.update({ where: { id: "u1" }, data: { contacto: "sp0" } })
    await tx.unidad.updateMany({ where: { id: { in: ["u2", "u3", "u4"] } }, data: { contacto: "sp1" } })

    await tx.usuario.create({
      data: {
        id: "tester",
        email: "tester@example.com",
        password:
          "$2a$12$TWaUL3tJEuMNUfC7uiiAjelPshhEWyBePLxVzs34LWdi1TnpJN0ZO", // password
      },
    });

    const proveedor = await tx.proveedor.create({
      data: {
        id: "pvdr0",
        nombre: "Proveedor 0",
        rif: "J-123456789",
        email: "proveedor0@example.com",
        telefono: "04121234567",
      },
    });

    await tx.iGasto.create({
      data: {
        id: "g0",
        concepto: "Some",
        proveedor: proveedor.id,
        monto: 25_00,
        moneda: Moneda.VED,
        descripcion: "Gasto 0",
        tasa: 12_50,
        registrado_por: "tester",
      },
    });

    const cuotas = await tx.cuota.createManyAndReturn({
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

    await tx.proyecto.createMany({
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
      await tx.iDeuda.createMany({
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

    const pago = await tx.iPago.create({
      data: {
        id: "p0",
        unidad: "villa-500",
        metodo: MetodoDePago.EFECTIVO,
        moneda: Moneda.VED,
        monto: 8134_50, // 25 USD
        tasa: 325_38,
        registrado_por: "tester",
        actualizado_por: "tester",
      },
      select: {
        id: true,
      },
    });

    await tx.destinoDePago.create({
      data: { deuda: "d500[c0]v[500]", pago: pago.id, destinado: 8_75 },
    });
  });

  console.log("✅ Seeding completed.");
}

if (process.env.NODE_ENV !== "test") {
  main()
    .catch((e) => {
      console.error(e);
      process.exit(1);
    })
    .finally(async () => {
      await prisma.$disconnect();
    });
}

import {
  EstadoDeProyecto,
  MetodoDePago,
  Moneda,
  TipoDeCuota,
} from "../../generated/prisma/client";
import { prisma } from "../client";

async function main() {
  console.log("Seeding database...");

  await prisma.$transaction(async (tx) => {
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

    const villa = await tx.villa.create({
      data: {
        numero: 362,
        IDeudas: {
          createMany: {
            data: cuotas.map((cuota, i) => ({
              id: `d${i}`,
              cuota: cuota.id,
            })),
          },
        },
      },
      select: {
        numero: true,
        IDeudas: {
          select: {
            id: true,
            Cuota: true,
          },
        },
      },
    });

    const pago = await tx.iPago.create({
      data: {
        id: "p0",
        villa: villa.numero,
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
      data: {
        pago: pago.id,
        deuda: villa.IDeudas[0]!.id,
        destinado: villa.IDeudas[0]!.Cuota.monto,
        fecha: new Date(),
      },
    });
  });

  console.log("Seeding completed.");
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });

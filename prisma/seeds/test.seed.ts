import fs from "node:fs";
import path from "node:path";
import {
  MetodoDePago,
  Moneda,
  TipoDeProveedor,
} from "../../generated/prisma/client";
import { prisma } from "../client";

async function main() {
  console.log("Seeding database...");

  // 1. Definir la ruta relativa desde la raíz del proyecto
  const relativeTarget = path.join("sql", "views");

  // 2. Obtener la ruta absoluta desde el root del proyecto (CWD)
  const fullPath = path.join(process.cwd(), relativeTarget);

  try {
    // 3. Leer archivos recursivamente
    const files = fs.readdirSync(fullPath, { recursive: true });

    // 4. Filtrar solo archivos .sql y reconstruir la ruta relativa al root
    const sqlFiles = files
      .filter((file) => (file as string).endsWith(".sql"))
      .map((file) => path.join(relativeTarget, file as string));

    for (const file of files) {
      const filePath = path.join(fullPath, file as string);

      // 2. Leer el contenido SQL del archivo
      const sql = fs.readFileSync(filePath, "utf8");

      console.log(`Ejecutando: ${file}...`);

      // 3. Ejecutar en la base de datos usando Prisma
      // Se usa $executeRawUnsafe porque el SQL viene de un string dinámico (el archivo)
      await prisma.$executeRawUnsafe(sql);
    }
  } catch (error) {
    console.error("Asegúrate de que la carpeta 'sql/views' exista en la raíz.");
  }
  await prisma.$transaction(async (tx) => {
    await tx.usuario.create({
      data: {
        id: "tester",
        email: "tester@example.com",
        password: "password",
      },
    });

    const proveedor = await tx.proveedor.create({
      data: {
        id: "pvdr0",
        nombre: "Proveedor 0",
        rif: "J-123456789",
        tipo: TipoDeProveedor.PERSONA_NATURAL,
        email: "proveedor0@example.com",
        telefono: "04121234567",
      },
    });

    await tx.iGasto.create({
      data: {
        id: "g0",
        proveedor: proveedor.id,
        monto: 25_00,
        moneda: Moneda.VED,
        descripcion: "Gasto 0",
        tasa: 12_50,
        registrado_por: "tester",
        actualizado_por: "tester",
      },
    });

    const cuotas = await tx.cuota.createManyAndReturn({
      data: [
        {
          id: "c0",
          monto: 8_75,
          mes: 1,
          anio: 2026,
        },
        {
          id: "c1",
          monto: 9_22,
          mes: 2,
          anio: 2026,
        },
        {
          id: "c3",
          monto: 8_94,
          mes: 3,
          anio: 2026,
        },
      ],
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

import fs from "node:fs";
import path from "node:path";
import { prisma } from "../client";

async function main() {
  console.log("⏳ Creating views...");

  const relativeTarget = path.join("apps", "api", "sql", "views");
  const fullPath = path.join(process.cwd(), relativeTarget);

  try {
    const files = fs.readdirSync(fullPath, { recursive: true });

    const sqlFiles = files
      .filter((file) => (file as string).endsWith(".sql"))
      .map((file) => path.join(relativeTarget, file as string));

    for (const file of sqlFiles) {
      const filePath = path.join(process.cwd(), file);
      const sql = fs.readFileSync(filePath, "utf8");

      console.log(`🌱 Executing: ${file}...`);
      await prisma.$executeRawUnsafe(sql);
    }

    console.log("✅ Views created.");
  } catch (error) {
    console.error("❌ Ensure 'sql/views' folder exists.");
    console.error(error);
  }
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });

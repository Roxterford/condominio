import "dotenv/config";
import { execSync } from "node:child_process";
import fs from "node:fs";
import path from "node:path";

console.log("⏳ Creating views...");

const dbPath = process.env.DATABASE_URL!.replace("file:", "");
const relativeTarget = path.join("sql", "views");
const fullPath = path.join(process.cwd(), relativeTarget);
console.log("exec path", fullPath);

try {
  const files = fs.readdirSync(fullPath, { recursive: true });

  const sqlFiles = files
    .filter((file: string | Buffer) => file.toString().endsWith(".sql"))
    .map((file: string | Buffer) => path.join(relativeTarget, file.toString()));

  for (const file of sqlFiles) {
    const filePath = path.join(process.cwd(), file);
    const sql = fs.readFileSync(filePath, "utf8");

    console.log(`🌱 Executing: ${file}...`);
    execSync(`sqlite3 "${dbPath}" "${sql.replace(/"/g, '\\"')}"`, {
      stdio: "inherit",
    });
  }

  console.log("✅ Views created.");
} catch (error) {
  console.error("❌ Ensure 'sql/views' folder exists.");
  console.error(error);
  process.exit(1);
}

#!/usr/bin/env tsx

import { Command } from 'commander';
import chalk from 'chalk';
import { promises as fs } from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

interface SeedModule {
  main?: () => Promise<void>;
}

const SEEDS_DIR = path.dirname(fileURLToPath(import.meta.url));
const SEED_PATTERN = /\.seed\.ts$/;

async function discoverSeeds(): Promise<string[]> {
  const files = await fs.readdir(SEEDS_DIR, { withFileTypes: true });
  const seeds: string[] = [];

  for (const file of files) {
    if (file.isFile() && SEED_PATTERN.test(file.name) && file.name !== 'index.ts') {
      seeds.push(file.name.replace(SEED_PATTERN, ''));
    }

    if (file.isDirectory()) {
      const subDir = path.join(SEEDS_DIR, file.name);
      const subFiles = await fs.readdir(subDir);

      for (const subFile of subFiles) {
        if (SEED_PATTERN.test(subFile)) {
          seeds.push(`${file.name}/${subFile.replace(SEED_PATTERN, '')}`);
        }
      }
    }
  }

  return seeds.sort();
}

async function loadSeed(seedName: string): Promise<SeedModule | null> {
  const isNested = seedName.includes('/');
  const modulePath = isNested
    ? path.join(SEEDS_DIR, seedName.replace('/', '/'))
    : path.join(SEEDS_DIR, seedName);

  const fullPath = `${modulePath}.seed.ts`;
  const mod = await import(fullPath);

  if (typeof mod.main === 'function') {
    return mod as SeedModule;
  }

  return null;
}

async function runSeed(seedName: string): Promise<void> {
  console.log(chalk.blue(`\n▶ Running seed: ${seedName}...`));

  try {
    const mod = await loadSeed(seedName);

    if (!mod) {
      console.log(chalk.yellow(`⚠ Seed '${seedName}' no exporta una función main().`));
      return;
    }

    await mod.main();

    console.log(chalk.green(`✓ Seed completed: ${seedName}`));
  } catch (error) {
    console.error(chalk.red(`✗ Seed failed: ${seedName}`));
    if (error instanceof Error) {
      console.error(chalk.yellow(error.message));
    }
    process.exit(1);
  }
}

async function promptSelectSeeds(availableSeeds: string[]): Promise<string[]> {
  const { default: inquirer } = await import('inquirer');

  const choices = availableSeeds.map((seed) => ({
    name: seed,
    value: seed,
    checked: false,
  }));

  const { seeds } = await inquirer.prompt([
    {
      type: 'checkbox',
      name: 'seeds',
      message: 'Select seeds to run:',
      choices,
      default: [],
    },
  ]);

  return seeds;
}

async function promptOrderSeeds(selectedSeeds: string[]): Promise<string[]> {
  if (selectedSeeds.length <= 1) {
    return selectedSeeds;
  }

  const { default: inquirer } = await import('inquirer');

  const { order } = await inquirer.prompt([
    {
      type: 'list',
      name: 'order',
      message: 'Select execution order:',
      choices: ['As selected', 'Reverse order'],
      default: 'As selected',
    },
  ]);

  return order === 'Reverse order' ? [...selectedSeeds].reverse() : selectedSeeds;
}

function printHelp(availableSeeds: string[]): void {
  console.log(`
${chalk.bold('🌱 Condominio Seeder CLI')}

${chalk.cyan('Usage:')}
  tsx prisma/seeds/index.ts [options] [seeds...]

${chalk.cyan('Options:')}
  -i, --interactive    Run in interactive mode
  -s, --seed <names>   Specify seeds to run (comma-separated)
  -o, --order <order>  Specify execution order (name or reverse)
  -h, --help          Show this help message
  -v, --version      Show version

${chalk.cyan('Available seeds:')}
${availableSeeds.map((s) => `  - ${s}`).join('\n')}

${chalk.cyan('Examples:')}
  tsx prisma/seeds/index.ts
  tsx prisma/seeds/index.ts --interactive
  tsx prisma/seeds/index.ts -s views mocks
  tsx prisma/seeds/index.ts -s views --order reverse
  `);
}

async function main() {
  const availableSeeds = await discoverSeeds();

  const program = new Command();

  program
    .name('prisma/seeds/index')
    .description('Condominio database seeder CLI')
    .version('1.0.0')
    .argument('[seeds...]', 'Seeds to run (space-separated)')
    .option('-i, --interactive', 'Run in interactive mode')
    .option('-s, --seed <names>', 'Seeds to run (comma-separated)')
    .option('-o, --order <order>', 'Execution order: name or reverse')
    .option('-h, --help', 'Show help')
    .parse(process.argv);

  const opts = program.opts();
  const args = program.args;

  if (opts.help) {
    printHelp(availableSeeds);
    return;
  }

  let seedsToRun: string[];

  if (opts.interactive) {
    console.log(chalk.yellow('🌱 Interactive mode'));

    const { default: inquirer } = await import('inquirer');

    const choices = availableSeeds.map((seed) => ({
      name: seed,
      value: seed,
    }));

    const { seeds } = await inquirer.prompt([
      {
        type: 'checkbox',
        name: 'seeds',
        message: 'Select seeds to run:',
        choices,
      },
    ]);

    seedsToRun = await promptOrderSeeds(seeds);
  } else if (args.length > 0) {
    const requested = args.map((s: string) => s.trim());
    const invalid = requested.filter((r: string) => !availableSeeds.includes(r));

    if (invalid.length > 0) {
      console.error(chalk.red(`✗ Unknown seeds: ${invalid.join(', ')}`));
      console.error(chalk.red(`Available: ${availableSeeds.join(', ')}`));
      process.exit(1);
    }

    seedsToRun = opts.order === 'reverse' ? [...requested].reverse() : requested;
  } else if (opts.seed) {
    const requested = opts.seed.split(',').map((s: string) => s.trim());
    const invalid = requested.filter((r: string) => !availableSeeds.includes(r));

    if (invalid.length > 0) {
      console.error(chalk.red(`✗ Unknown seeds: ${invalid.join(', ')}`));
      console.error(chalk.red(`Available: ${availableSeeds.join(', ')}`));
      process.exit(1);
    }

    seedsToRun = opts.order === 'reverse' ? [...requested].reverse() : requested;
  } else {
    console.log(chalk.yellow('No seeds specified. Use --help for usage.'));
    printHelp(availableSeeds);
    return;
  }

  if (seedsToRun.length === 0) {
    console.log(chalk.yellow('No seeds selected.'));
    process.exit(0);
  }

  console.log(chalk.cyan(`\n🚀 Running ${seedsToRun.length} seed(s)...`));

  for (const seed of seedsToRun) {
    await runSeed(seed);
  }

  console.log(chalk.green('\n✅ All seeds completed successfully!'));
}

main().catch((error) => {
  console.error(chalk.red('\n✗ Fatal error:'), error);
  process.exit(1);
});
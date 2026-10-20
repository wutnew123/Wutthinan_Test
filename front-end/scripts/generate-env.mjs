import { existsSync, readFileSync, writeFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = resolve(import.meta.dirname, '..');
const envPath = resolve(root, '.env');
const outputPath = resolve(root, 'src/environments/environment.ts');

if (!existsSync(envPath)) throw new Error('.env file not found');

const values = {};
for (const line of readFileSync(envPath, 'utf8').split(/\r?\n/)) {
  const match = line.match(/^\s*([A-Z0-9_]+)\s*=\s*(.*?)\s*$/);
  if (match) values[match[1]] = match[2].replace(/^['"]|['"]$/g, '');
}

const apiUrl = values.API_URL ?? values.URL;
if (!apiUrl) throw new Error('API_URL is missing in .env');

const content = [
  'export const environment = {',
  '  production: false,',
  `  apiUrl: ${JSON.stringify(apiUrl)}`,
  '};',
  ''
].join('\n');

writeFileSync(outputPath, content);
console.log(`Generated ${outputPath}`);

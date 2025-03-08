import esbuild from 'esbuild'

const banner = `import * as module__ from 'module';
import * as url__ from 'url';
import * as path__ from 'path';

const __filename = url__.fileURLToPath(import.meta.url);
const __dirname = path__.dirname(__filename);
const require = module__.createRequire(import.meta.url);
`

esbuild
  .build({
    entryPoints: ['src/main.ts'],
    bundle: true,
    minify: true,
    platform: 'node',
    outfile: 'bin/phishell-js.mjs',
    target: ['node18'],
    format: 'esm',
    external: ['readline/promises'],
    legalComments: 'none',
    banner: {
      js: banner,
    },
  })
  .then(() => {
    console.log('Build complete.')
  })
  .catch(() => process.exit(1))

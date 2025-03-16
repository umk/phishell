import { exec } from 'child_process'
import { promisify } from 'util'

import { PackageInfo } from './getPackageInfo'

const execPromise = promisify(exec)

async function preparePackage(packageDir: string, packageInfo: PackageInfo) {
  if (!packageInfo.scripts.has('build')) {
    throw new Error('No build script found in package.json')
  }
  try {
    process.stderr.write('Installing npm packages...\n')
    await execPromise('npm install', { cwd: packageDir })

    process.stderr.write('Running build command...\n')
    await execPromise('npm run build', { cwd: packageDir })

    process.stderr.write('Build completed successfully\n')
  } catch (error) {
    throw new Error(`Build failed: ${error}`)
  }
}

export default preparePackage

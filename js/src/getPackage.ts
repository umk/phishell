import fs from 'fs/promises'
import { createRequire } from 'module'
import path from 'path'

async function getPackage(packageDir: string): Promise<Record<string, unknown>> {
  const packageJsonPath = path.join(packageDir, 'package.json')

  const packageJsonContent = await fs.readFile(packageJsonPath, 'utf-8')
  const packageJson = JSON.parse(packageJsonContent)

  const isESModule = packageJson.type === 'module'
  const mainFile = packageJson.main || 'index.js'
  const modulePath = path.join(packageDir, mainFile)

  let module: unknown

  if (isESModule) {
    module = await import(modulePath)
  } else {
    const require = createRequire(import.meta.url)
    module = require(modulePath)
  }

  if (typeof module !== 'object' || module === null) {
    throw new Error(`Module at ${modulePath} did not export an object`)
  }

  return module as Record<string, unknown>
}

export default getPackage

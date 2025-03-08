import { createRequire } from 'module'

import { PackageInfo } from './getPackageInfo'

async function getPackage(packageInfo: PackageInfo): Promise<Record<string, unknown>> {
  const { isESModule, modulePath } = packageInfo

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

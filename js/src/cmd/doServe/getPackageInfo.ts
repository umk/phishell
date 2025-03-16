import fs from 'fs/promises'
import path from 'path'

import Ajv from '../../common/Ajv'

export type PackageInfo = {
  isESModule: boolean
  modulePath: string
  scripts: Set<string>
  /** Path to the TypeScript definition file for the package */
  types: string | undefined
}

const PACKAGE_JSON_SCHEMA = {
  type: 'object',
  properties: {
    type: { type: 'string' },
    main: { type: 'string' },
    scripts: { type: 'object', additionalProperties: { type: 'string' } },
    types: { type: 'string' },
  },
  additionalProperties: true,
}

async function getPackageInfo(packageDir: string): Promise<PackageInfo> {
  const packageJsonPath = path.join(packageDir, 'package.json')
  const packageJsonContent = await fs.readFile(packageJsonPath, 'utf-8')
  const packageJson = JSON.parse(packageJsonContent)

  if (!Ajv.validate(PACKAGE_JSON_SCHEMA, packageJson)) {
    throw new Error('package.json has errors.')
  }

  const { type, main, scripts, types } = packageJson
  return {
    isESModule: type === 'module',
    modulePath: path.join(packageDir, main || 'dist/index.js'),
    scripts: new Set<string>(Object.keys(scripts ?? {})),
    types,
  }
}

export default getPackageInfo

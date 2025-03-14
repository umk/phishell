import fs from 'fs/promises'
import path from 'path'

export type PackageInfo = {
  isESModule: boolean
  modulePath: string
  scripts: Set<string>
}

async function getPackageInfo(packageDir: string): Promise<PackageInfo> {
  const packageJsonPath = path.join(packageDir, 'package.json')

  const packageJsonContent = await fs.readFile(packageJsonPath, 'utf-8')
  const packageJson = JSON.parse(packageJsonContent)

  const isESModule = packageJson.type === 'module'
  const mainFile = packageJson.main || 'dist/index.js'
  const modulePath = path.join(packageDir, mainFile)

  const scripts = new Set<string>(Object.keys(packageJson.scripts ?? {}))

  return { isESModule, modulePath, scripts }
}

export default getPackageInfo

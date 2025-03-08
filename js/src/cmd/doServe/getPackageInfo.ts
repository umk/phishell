import fs from 'fs/promises'
import path from 'path'

export type PackageInfo = {
  isESModule: boolean
  modulePath: string
}

async function getPackageInfo(packageDir: string): Promise<PackageInfo> {
  const packageJsonPath = path.join(packageDir, 'package.json')

  const packageJsonContent = await fs.readFile(packageJsonPath, 'utf-8')
  const packageJson = JSON.parse(packageJsonContent)

  const isESModule = packageJson.type === 'module'
  const mainFile = packageJson.main || 'index.js'
  const modulePath = path.join(packageDir, mainFile)

  return { isESModule, modulePath }
}

export default getPackageInfo

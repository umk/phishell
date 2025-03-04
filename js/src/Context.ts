import { getPackageInfo, PackageInfo } from 'fn-json-schema'

import getFunctions, { FunctionInfo } from './getFunctions.js'
import getPackage from './getPackage.js'

type Context = {
  functions: Array<FunctionInfo>
}

export async function createContext(): Promise<Context> {
  return {
    ...(await createFunctionsContext(process.cwd())),
  }
}

async function createFunctionsContext(packageDir: string): Promise<Pick<Context, 'functions'>> {
  let packageInfo: PackageInfo
  try {
    packageInfo = await getPackageInfo(packageDir)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } catch (error: any) {
    throw new Error(`Cannot get package information: ${error.message}`)
  }
  let p
  try {
    p = await getPackage(packageDir)
  } catch (error) {
    throw new Error(`Cannot load package: ${error}`)
  }
  const functions = await getFunctions(packageDir, packageInfo, p)
  return { functions }
}

export default Context

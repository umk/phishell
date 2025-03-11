import getFunctions, { FunctionInfo } from './getFunctions.js'
import getPackage from './getPackage.js'
import { PackageInfo } from './getPackageInfo.js'

type Context = {
  functions: Array<FunctionInfo>
  modulePath: string
}

export async function createContext(
  packageDir: string,
  packageInfo: PackageInfo,
): Promise<Context> {
  return {
    ...(await createFunctionsContext(packageDir, packageInfo)),
  }
}

async function createFunctionsContext(
  packageDir: string,
  packageInfo: PackageInfo,
): Promise<Pick<Context, 'functions' | 'modulePath'>> {
  let p
  try {
    p = await getPackage(packageInfo)
  } catch (error) {
    throw new Error(`Cannot load package: ${error}`)
  }
  const functions = await getFunctions(packageDir, p)
  return {
    functions,
    modulePath: packageInfo.modulePath,
  }
}

export default Context

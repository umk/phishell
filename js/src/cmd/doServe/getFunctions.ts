import {
  getPackageInfo,
  getPackageSchema,
  JsonSchema,
  JsonSchemaObject,
  PackageFunction,
  PackageFunctionSignature,
  PackageInfo,
} from 'fn-json-schema'

const DEFAULT_OBJECT_SCHEMA: JsonSchema & JsonSchemaObject = {
  type: 'object',
  properties: {},
}

export type FunctionParameterSchema = {
  schema: JsonSchema & JsonSchemaObject
}

export type FunctionInfo = {
  name: string
  f: PackageFunction & PackageFunctionSignature
  parameter: FunctionParameterSchema
  handler: (...params: Array<unknown>) => unknown
}

class FunctionError extends Error {
  constructor(public readonly details: string) {
    super('Unable to process the function definition.')
  }
}

function processFunctionDefinition(
  f: PackageFunction,
  functions: Array<FunctionInfo>,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  p: any,
): boolean {
  try {
    functions.push(getFunction(f, p))
    return true
  } catch (error) {
    if (error instanceof FunctionError) {
      process.stderr.write(`${f.name} - ${error.details}\n`)
      return false
    }
    throw error
  }
}

function getFunction(
  f: PackageFunction,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  p: any,
): FunctionInfo {
  if ('error' in f) {
    throw new FunctionError(f.error)
  }
  const handler = p[f.name]
  if (handler === undefined) {
    throw new FunctionError('handler is not exported')
  }
  if (typeof handler !== 'function') {
    throw new FunctionError('handler is not a function')
  }
  if (f.signature.parameters.length > 1) {
    throw new FunctionError('handler cannot define more than one parameter')
  }
  const parameter = f.signature.parameters[0]
  return {
    name: f.name,
    f: f as PackageFunction & PackageFunctionSignature,
    parameter: getParameterSchema(parameter ? parameter.schema : DEFAULT_OBJECT_SCHEMA),
    handler,
  }
}

function getParameterSchema(schema: JsonSchema | undefined): FunctionParameterSchema {
  if (schema && 'type' in schema && schema.type === 'object') {
    return { schema }
  }
  throw new FunctionError('parameter schema could not be determined')
}

async function getFunctions(
  packageDir: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  p: any,
): Promise<Array<FunctionInfo>> {
  let packageInfo: PackageInfo
  try {
    packageInfo = await getPackageInfo(packageDir)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } catch (error: any) {
    throw new Error(`Cannot get package information: ${error.message}`)
  }
  if (!packageInfo.types) {
    throw new Error("The package doesn't have a reference to types")
  }
  const functions: Array<FunctionInfo> = []
  try {
    const fns = await getPackageSchema(packageDir, packageInfo.types)
    const results = fns.map((f) => processFunctionDefinition(f, functions, p))
    const hasErrors = results.includes(false)
    if (hasErrors) {
      throw new Error('Some functions could not be loaded')
    }
    if (functions.length === 0) {
      throw new Error("The package doesn't export any functions or all of them could not be loaded")
    }
  } catch (error) {
    throw new Error(`Cannot process package type definitions: ${error}`)
  }
  return functions
}

export default getFunctions

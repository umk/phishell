import Ajv from '../../common/Ajv'
import { ToolRequest } from '../../protocol'

import Context from './Context'
import { FunctionInfo } from './getFunctions'

function createInvoke(context: Context) {
  const functions = context.functions.reduce((prev, cur) => {
    prev.set(cur.name, cur)
    return prev
  }, new Map<string, FunctionInfo>())
  return async function (request: ToolRequest) {
    const func = functions.get(request.function.name)
    if (!func) {
      throw new Error(`Function "${request.function.name}" not found.`)
    }

    const argumentz = JSON.parse(request.function.arguments)
    if (!Ajv.validate(func.parameter, argumentz)) {
      throw new Error("The request didn't match the schema.")
    }

    return await func.handler(argumentz)
  }
}

export default createInvoke

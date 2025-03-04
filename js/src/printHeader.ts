import Context from './Context'
import { FunctionInfo } from './getFunctions'

function printHeader(context: Context) {
  for (const f of context.functions) {
    printFunction(f)
  }
  process.stdout.write('\n\n') // Print empty line to indicate the end of header
}

function printFunction(info: FunctionInfo) {
  process.stdout.write(
    JSON.stringify({
      type: 'function',
      function: {
        name: info.name,
        description: info.f.signature.description,
        parameters: info.parameter.schema,
      },
    }),
  )
}

export default printHeader

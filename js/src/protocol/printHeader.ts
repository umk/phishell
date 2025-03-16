import Tool, { ToolFunction } from './Tool'

function printHeader(functions: Array<ToolFunction>) {
  for (const f of functions) {
    const tool: Tool = { type: 'function', function: f }
    process.stdout.write(JSON.stringify(tool))
  }
  process.stdout.write('\n\n') // Print empty line to indicate the end of header
}

export default printHeader

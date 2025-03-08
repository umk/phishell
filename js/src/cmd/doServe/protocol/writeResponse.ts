import { ToolResponse } from './ToolRequest'

function writeResponse(response: ToolResponse) {
  process.stdout.write(JSON.stringify(response))
  process.stdout.write('\n')
}

export default writeResponse

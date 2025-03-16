import { ToolResponse } from './ToolRequest'

function printResponse(response: ToolResponse) {
  process.stdout.write(JSON.stringify(response))
  process.stdout.write('\n')
}

export default printResponse

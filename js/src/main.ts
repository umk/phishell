import Context, { createContext } from './Context'
import createInvoke from './createInvoke'
import getLines from './getLines'
import getRequest from './getRequest'
import printHeader from './printHeader'
import { ToolResponse } from './ToolRequest'

async function main() {
  try {
    const context = await createContext()

    printHeader(context)
    await cli(context)
    process.exit(0)
  } catch (error) {
    process.stderr.write(String(error))
    process.stderr.write('\n')
    process.exit(1)
  }
}

async function cli(context: Context) {
  const invoke = createInvoke(context)
  for await (const line of getLines(process.stdin)) {
    const message = line.trim()
    if (message) {
      try {
        const request = getRequest(message)

        try {
          const result = await invoke(request)

          const response: ToolResponse = {
            call_id: request.call_id,
            content: JSON.stringify(result),
          }

          process.stdout.write(JSON.stringify(response))
          process.stdout.write('\n')
        } catch (error) {
          const message =
            error instanceof Error
              ? error.message
              : 'An error has occurred when running the function'

          const response: ToolResponse = {
            call_id: request.call_id,
            content: JSON.stringify(message),
          }

          process.stdout.write(JSON.stringify(response))
          process.stdout.write('\n')
        }
      } catch (error) {
        process.stdout.write(String(error))
        process.stdout.write('\n')
      }
    }
  }
}

main()

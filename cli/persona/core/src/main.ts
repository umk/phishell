import { parseArguments } from './Arguments'
import { getConfig } from './Config'
import Context, { createContext } from './Context'
import getLines from './getLines'
import getRequest, { getRequestArgs, ToolResponse } from './getRequest'

async function main() {
  try {
    const argumentz = await parseArguments()
    const config = await getConfig(argumentz)
    const context = await createContext(config)

    printHeader()
    await cli(context)
    process.exit(0)
  } catch (error) {
    process.stderr.write(String(error))
    process.stderr.write('\n')
    process.exit(1)
  }
}

async function cli(_context: Context) {
  for await (const line of getLines(process.stdin)) {
    const message = line.trim()
    if (message) {
      try {
        const request = getRequest(message)
        const args = getRequestArgs(request.function.arguments)

        const response: ToolResponse = {
          call_id: request.call_id,
          content: args.message,
        }

        process.stdout.write(JSON.stringify(response))
        process.stdout.write('\n')
      } catch (error) {
        process.stdout.write(String(error))
        process.stdout.write('\n')
      }
    }
  }
}

function printHeader() {
  // A bogus function that is meant to accept the user message and
  // process it according to the logic of persona.
  process.stdout.write(
    JSON.stringify({
      type: 'function',
      function: {
        name: 'Post',
        description: 'Processes a given message.',
        parameters: {
          type: 'object',
          properties: {
            message: {
              type: 'string',
              description: 'The message to be processed.',
            },
          },
          required: ['message'],
        },
      },
    }),
  )
  process.stdout.write('\n\n') // Print empty line to indicate the end of header
}

main()

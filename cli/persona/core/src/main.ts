import { parseArguments } from './Arguments'
import { getConfig } from './Config'
import Context, { createContext } from './Context'
import RL from './RL'

async function main() {
  try {
    const argumentz = await parseArguments()
    const config = await getConfig(argumentz)
    const context = await createContext(config)
    await cli(context)
    process.exit(0)
  } catch (error) {
    process.stderr.write(String(error))
    process.stderr.write('\n')
    process.exit(1)
  }
}

async function cli(_context: Context) {
  for (;;) {
    const raw = await RL.question('')
    const message = raw.trim()
    if (message) {
      try {
      } catch (error) {
        process.stderr.write(String(error))
        process.stderr.write('\n')
      }
    }
  }
}

main()

import { runCommand } from './cmd'
import initClient from './initClient'

async function main() {
  try {
    initClient()
    await runCommand()
  } catch (error) {
    process.stderr.write(String(error))
    process.stderr.write('\n')
    process.exit(1)
  }
}

main()

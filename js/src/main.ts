import { runCommand } from './cmd'

async function main() {
  try {
    await runCommand()
  } catch (error) {
    process.stderr.write(String(error))
    process.stderr.write('\n')
    process.exit(1)
  }
}

main()

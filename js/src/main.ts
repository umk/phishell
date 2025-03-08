import parseArguments from './parseArguments'

async function main() {
  try {
    await parseArguments()
    process.exit(0)
  } catch (error) {
    process.stderr.write(String(error))
    process.stderr.write('\n')
    process.exit(1)
  }
}

main()

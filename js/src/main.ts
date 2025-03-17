import { runCommand } from './cmd'
import initClient from './initClient'

async function main() {
  try {
    initClient()

    const response = await getClient().embeddings.create({
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      model: process.env.PHISHELL_PROFILE!,
      input: 'my input',
    })
    console.error(response.data[0].embedding)

    await runCommand()
  } catch (error) {
    process.stderr.write(String(error))
    process.stderr.write('\n')
    process.exit(1)
  }
}

main()

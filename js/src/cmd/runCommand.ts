import yargs from 'yargs'
import { hideBin } from 'yargs/helpers'

import doServe from './doServe'

async function runCommand(): Promise<void> {
  await yargs(hideBin(process.argv))
    .command(
      'serve',
      'Serve exported functions to the host',
      () => {
        // No arguments
      },
      () => doServe(),
    )
    .demandCommand(1, 'Command name is required')
    .strict()
    .help()
    .parse()
}

export default runCommand

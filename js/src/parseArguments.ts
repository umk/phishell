import yargs from 'yargs'
import { hideBin } from 'yargs/helpers'

import doCheck from './cmd/doCheck'
import doServe from './cmd/doServe'

async function parseArguments(): Promise<void> {
  await yargs(hideBin(process.argv))
    .command(
      'serve',
      'Serve exported functions to the host',
      () => {
        // No arguments
      },
      () => doServe(),
    )
    .command(
      'check',
      'Read module source from Stdin and check for errors',
      () => {
        // No arguments
      },
      () => doCheck(),
    )
    .demandCommand(1, 'Command name is required')
    .strict()
    .help()
    .parse()
}

export default parseArguments

import readline from 'readline/promises'

const RL = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: true,
  prompt: '',
  historySize: 0,
})

export function setupRL() {
  RL.on('SIGINT', () => {
    RL.close()
    process.exit()
  })
}

export function readKey(prompt: string) {
  RL.write(prompt)
  return new Promise<string>((resolve) => {
    const handler = (chunk: string | undefined) => {
      process.stdin.removeListener('keypress', handler)
      chunk && resolve(chunk)
    }
    process.stdin.on('keypress', handler)
  })
}

export async function readKeyOf(prompt: string, ...keys: Array<string>) {
  let result: string | undefined
  const keys_ = keys.map((k) => k.toLowerCase())
  for (;;) {
    try {
      result = await readKey(prompt)
      if (keys_.includes(result.toLowerCase())) {
        return result
      }
    } finally {
      if (result !== '\r') {
        RL.write('\n')
      }
    }
  }
}

export async function readKeyOrDefaultOf(
  prompt: string,
  default_: string,
  ...other: Array<string>
) {
  const result = await readKeyOf(prompt, '\r', default_, ...other)
  return result === '\r' ? default_ : result
}

export default RL

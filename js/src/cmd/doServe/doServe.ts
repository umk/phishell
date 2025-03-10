import chokidar from 'chokidar'

import Context, { createContext } from './Context'
import createInvoke from './createInvoke'
import { getRequest, readRequests, writeHeader, writeResponse } from './protocol'

async function doServe() {
  const context = await createContext()
  const invoke = createInvoke(context)

  writeHeader(context)

  restartOnChange(context, 3000)

  for await (const line of readRequests(process.stdin)) {
    const message = line.trim()
    if (!message) {
      continue
    }
    try {
      const request = getRequest(message)

      try {
        const result = await invoke(request)

        writeResponse({
          call_id: request.call_id,
          content: JSON.stringify(result),
        })
      } catch (error) {
        const message =
          error instanceof Error ? error.message : 'An error has occurred when running the function'

        writeResponse({
          call_id: request.call_id,
          content: JSON.stringify(message),
        })
      }
    } catch (error) {
      process.stderr.write(String(error))
      process.stderr.write('\n')
    }
  }
}

function restartOnChange(context: Context, debounceMs: number) {
  const watcher = chokidar.watch(context.modulePath, {
    persistent: false,
    ignoreInitial: true,
    depth: Infinity,
  })
  let timeout: NodeJS.Timeout | undefined = undefined
  watcher.on('all', () => {
    if (timeout) {
      clearTimeout(timeout)
    }
    timeout = setTimeout(() => {
      // Code 99 indicates that tools provider wants to be restarted
      process.exit(99)
    }, debounceMs)
  })
}

export default doServe

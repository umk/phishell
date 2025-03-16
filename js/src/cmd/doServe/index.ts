import path from 'path'

import chokidar from 'chokidar'

import { getRequest, printHeader, printResponse, readRequests } from '../../protocol'

import { createContext } from './Context'
import createInvoke from './createInvoke'
import getPackageInfo, { PackageInfo } from './getPackageInfo'
import preparePackage from './preparePackage'

async function doServe() {
  const packageDir = process.cwd()
  let packageInfo: PackageInfo
  try {
    packageInfo = await getPackageInfo(packageDir)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } catch (error: any) {
    throw new Error(`Cannot get package information: ${error.message}`)
  }

  await preparePackage(packageDir, packageInfo)

  const context = await createContext(packageDir, packageInfo)
  const invoke = createInvoke(context)

  printHeader(
    context.functions.map((f) => ({
      name: f.name,
      description: f.f.signature.description,
      parameters: f.parameter.schema,
    })),
  )

  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  restartOnChange(path.dirname(packageInfo.types!), 3000)

  for await (const line of readRequests(process.stdin)) {
    const message = line.trim()
    if (!message) {
      continue
    }
    try {
      const request = getRequest(message)

      try {
        const result = await invoke(request)

        printResponse({
          call_id: request.call_id,
          content: JSON.stringify(result),
        })
      } catch (error) {
        const message =
          error instanceof Error ? error.message : 'An error has occurred when running the function'

        printResponse({
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

function restartOnChange(packageDir: string, debounceMs: number) {
  const watcher = chokidar.watch([`${packageDir}/**/*.js`, `${packageDir}/**/*.ts`], {
    persistent: false,
    ignoreInitial: true,
    ignored: [`${packageDir}/node_modules/**`],
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

import OpenAI from 'openai'
import { Fetch } from 'openai/core'
import * as undici from 'undici'

function initClient() {
  const socketPath = process.env.PHISHELL_SOCKET
  let client: OpenAI | undefined
  if (socketPath) {
    client = createClient(socketPath)
  }

  globalThis.getClient = () => {
    if (!client) {
      throw new Error('Client is not set up.')
    }
    return client
  }
}

function createClient(socketPath: string) {
  return new OpenAI({
    fetch: (async (
      url: undici.RequestInfo,
      init?: undici.RequestInit,
    ): Promise<undici.Response> => {
      const { signal, method, body, headers } = init ?? {}
      return await undici.fetch(url, {
        signal,
        method,
        body,
        headers,
        dispatcher: new undici.Agent({
          connect: { socketPath },
        }),
      })
    }) as unknown as Fetch,
    baseURL: 'http://localhost', // force HTTP request instead of HTTPS
    apiKey: '1',
  })
}

export default initClient

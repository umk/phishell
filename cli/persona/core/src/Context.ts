import Config from './Config.js'

type Context = {
  config: Config
}

export async function createContext(config: Config): Promise<Context> {
  return {
    config,
  }
}

export default Context

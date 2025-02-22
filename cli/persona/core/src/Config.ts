import Joi from 'joi'

import Arguments from './Arguments'

type Config = {}

const configSchema: Joi.Schema<Config> = Joi.object({})

export async function getConfig(argumentz: Arguments): Promise<Config> {
  const config = await configSchema.validateAsync(argumentz)
  return config
}

export default Config

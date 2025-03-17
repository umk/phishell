/* eslint-disable no-var */
import OpenAI from 'openai'

export {}

declare global {
  var getClient: () => OpenAI
}
/* eslint-enable no-var */

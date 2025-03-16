import Ajv from '../common/Ajv'

import ToolRequest from './ToolRequest'
import { REQUEST_SCHEMA } from './schemas'

function getRequest(value: string): ToolRequest {
  const request = JSON.parse(value)
  if (!Ajv.validate(REQUEST_SCHEMA, request)) {
    throw new Error("The content didn't match the request schema.")
  }

  return request as ToolRequest
}

export default getRequest

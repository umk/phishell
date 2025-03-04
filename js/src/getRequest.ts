import Ajv from './Ajv'
import { REQUEST_SCHEMA } from './schemas'
import ToolRequest from './ToolRequest'

function getRequest(value: string): ToolRequest {
  const request = JSON.parse(value)
  if (!Ajv.validate(REQUEST_SCHEMA, request)) {
    throw new Error("The content didn't match the request schema.")
  }

  return request as ToolRequest
}

export default getRequest

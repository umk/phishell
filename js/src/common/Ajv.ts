import Ajv from 'ajv'

const ajv = new Ajv({
  allErrors: false,
  allowUnionTypes: true,
  useDefaults: true,
})

export default ajv

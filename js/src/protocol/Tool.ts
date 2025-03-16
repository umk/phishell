type Tool = {
  /** The type of the tool. Currently, only `function` is supported. */
  type: 'function'

  /** Description of the function. */
  function: ToolFunction
}

export type ToolFunction = {
  /**
   * The name of the function to be called. Must be a-z, A-Z, 0-9, or
   * contain underscores and dashes, with a maximum length of 64.
   */
  name: string

  /**
   * A description of what the function does, used by the model to
   * choose when and how to call the function.
   */
  description: string | undefined

  /**
   * The parameters the function accepts, described as a JSON Schema object.
   */
  parameters: unknown
}

export default Tool

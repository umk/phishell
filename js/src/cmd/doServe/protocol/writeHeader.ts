import Context from '../Context'

function writeHeader(context: Context) {
  for (const f of context.functions) {
    process.stdout.write(
      JSON.stringify({
        type: 'function',
        function: {
          name: f.name,
          description: f.f.signature.description,
          parameters: f.parameter.schema,
        },
      }),
    )
  }
  process.stdout.write('\n\n') // Print empty line to indicate the end of header
}

export default writeHeader

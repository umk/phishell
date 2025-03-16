async function* readRequests(stream: NodeJS.ReadStream): AsyncGenerator<string> {
  let leftover: Buffer | undefined = undefined

  for await (const chunk of stream) {
    let data: Buffer = leftover ? Buffer.concat([leftover, chunk]) : chunk
    let newlineIndex: number

    while ((newlineIndex = data.indexOf(10)) !== -1) {
      yield data.subarray(0, newlineIndex).toString('utf8') // Convert to string only before yielding
      data = data.subarray(newlineIndex + 1) // Keep remaining buffer
    }

    leftover = data.length > 0 ? data : undefined // Store remaining partial buffer
  }

  if (leftover) {
    yield leftover.toString('utf8') // Convert remaining data to string before yielding
  }
}

export default readRequests

# 𝛗 shell

**Phi Shell** is a command-line application designed to serve as a command processor with enhanced capabilities for integrating external tool providers and interacting with Large Language Models (LLMs). These tool providers are standalone programs created by users. Providers can be implemented in any programming language, taking advantage of thousands of libraries to enable features ranging from integrations with external systems to the implementation of sophisticated business logic. The host program communicates with tool providers through their `Stdin` and `Stdout`, reducing the need for extensive boilerplate code.

Phi Shell is shipped with type declarations for the [Go](provider) programming language, allowing users to implement their tools with few imports and function calls.

Phi Shell offers two distinct modes of operation:

- **Basic Shell**: This mode replicates the basic features of traditional shells while enabling users to integrate external programs (tool providers) with the host program. Additionally, the output generated by console commands can be added to the chat history, facilitating inference and interaction with LLMs.
- **Chat**: Interacts with LLMs directly to perform inference and leverages attached tools for enhanced functionality during LLM interactions.

The modes can be switched between using the <kbd>Tab</kbd> key.

When a tool call initiates an asynchronous operation, or if a provider starts an asynchronous process independently of a tool call, the provider can output messages to `Stdout` in the required format. This makes the host program capture the messages, making them available in the Inbox. The Inbox can be accessed by typing the `inbox` command or pressing the <kbd>Esc</kbd> key while in Basic Shell mode. The messages can then be reviewed by the user, who has the option to delete them or forward them to the chat for inference.

## Configuration

Phi Shell uses two YAML files as a source of configuration:

- **Global**: `.phishell.yaml` file in the user's home directory;
- **Local**: `.phishell.yaml` file in the current directory. This directory can be changed at the time of running Phi Shell (see the `-dir` argument for details).

Both files are not mandatory, but at least one of them must exist in order to run Phi Shell. The local configuration file takes precedence over the global one. The configuration file specifies profiles, which are basically the configuration of the LLM client and few other options that affect the interaction, and the name of default profile. On a high level the configuration YAML is defined as follows:

```yaml
profiles:
    profileA:
        # configuration of profile A
    profileB:
        # configuration of profile B

default: profileB
```

Each of the profiles specify the following data:

<table>
    <thead>
        <tr>
            <th align="left" width="175">Property</th>
            <th align="left">Description</th>
            <th width="125">Default</th>
        </tr>
    </thead>
    <tbody>
        <tr valign="top">
            <td><code>preset</code></td>
            <td>Prefills options with values suitable for particular use case. Valid values are <code>openai</code> and <code>ollama</code>.</td>
            <td></td>
        </tr>
        <tr valign="top">
            <td><code>baseurl</code></td>
            <td>Base URL of an LLM service</td>
            <td>OpenAI API</td>
        </tr>
        <tr valign="top">
            <td><code>key</code></td>
            <td>Key to use with the API. Must be specified even for services, which don't require the key (use a random value in this case). Can be replaced with <code>PHI_KEY</code> environment variable, which is applied to all profiles, unless the profile specifies key explicitly. If key is not specified, Phi Shell will ask for the key and preserve it in the key chain under the profile's name.</td>
            <td></td>
        </tr>
        <tr valign="top">
            <td><code>model</code></td>
            <td>LLM model to use</td>
            <td><code>GPT-4o</code></td>
        </tr>
        <tr valign="top">
            <td><code>prompt</code></td>
            <td>Additional instructions to LLM that fed to LLM in a system prompt to adjust interaction. When multiple profiles are specified for the session, only the prompt of the first profile is respected.</td>
            <td></td>
        </tr>
        <tr valign="top">
            <td><code>retries</code></td>
            <td>The number of times that Phi Shell will call LLM for inference in case if the response didn't pass validation.</td>
            <td>5</td>
        </tr>
        <tr valign="top">
            <td><code>concurrency</code></td>
            <td>How many concurrent calls can be made to LLM</td>
            <td>1</td>
        </tr>
        <tr valign="top">
            <td><code>compactionToks</code></td>
            <td>The number of tokens when Phi Shell will start compaction of the conversation history. This number can be exceeded anyway when running long chains of interaction with LLM, which involve several tools calls in a row.</td>
            <td>2000</td>
        </tr>
    </tbody>
</table>

Configuration example:

```yaml
profiles:
  qwen2.5:
    preset: ollama
    model: qwen2.5:14b

  openai: {}

default: qwen2.5
```

## Arguments

The following command line arguments are supported when running Phi Shell:

<table>
    <thead>
        <tr>
            <th align="left" width="175">Argument</th>
            <th align="left">Description</th>
        </tr>
    </thead>
    <tbody>
        <tr valign="top">
            <td><code>-debug</code></td>
            <td>Run program in debug mode to display additional information about interaction with LLM server.</td>
        </tr>
        <tr valign="top">
            <td><code>-dir [dir]</code></td>
            <td>Base directory. Defaults to current directory.</td>
        </tr>
        <tr valign="top">
            <td><code>-profile [name]</code></td>
            <td>Specify LLM profile to use for chat mode. If not specified, the default profile is used. There can be
                multiple profiles specified, and all of them will have their distinct mode, which can be cycled through
                using the <kbd>Tab</kbd> key. All the profiles will operate over the same chat history, but only the first profile will be used for LLM interactions behind the scenes (for example, the chat history compaction).</td>
        </tr>
        <tr valign="top">
            <td><code>-v</code></td>
            <td>Show version and quit.</td>
        </tr>
    </tbody>
</table>

## Commands

The Basic Shell mode supports the following built-in commands:

<table>
    <thead>
        <tr>
            <th align="left" width="175">Command</th>
            <th align="left">Description</th>
        </tr>
    </thead>
    <tbody>
        <tr valign="top">
            <td><code>attach [cmd]</code></td>
            <td>Run background process that provides tools and events. This actually attaches the tool provider. The
                command line of attaching program may include both the path to the program and command line arguments.
            </td>
        </tr>
        <tr valign="top">
            <td><code>cd [dir]</code></td>
            <td>Change the current directory</td>
        </tr>
        <tr valign="top">
            <td><code>export &lt;var&gt;</code></td>
            <td>If variable is specified, set value of an environment variable. Otherwise show the list of environment
                variables. To set the variable, it must be specified in a <code>name=value</code> format. The exported variables are propagated to the child processes.</td>
        </tr>
        <tr valign="top">
            <td><code>help</code></td>
            <td>Display the help message.</td>
        </tr>
        <tr valign="top">
            <td><code>history</code></td>
            <td>Display the chat messages history. Available only when running the Phi Shell with <code>-debug</code>
                flag.</td>
        </tr>
        <tr valign="top">
            <td><code>inbox</code></td>
            <td>Show the list of incoming messages from background process.</td>
        </tr>
        <tr valign="top">
            <td><code>jobs</code></td>
            <td>List attached background processes. Includes one terminated or failed background process, if any, for
                diagnostic purposes.</td>
        </tr>
        <tr valign="top">
            <td><code>kill [pid]</code></td>
            <td>Kill a background process with the given PID.</td>
        </tr>
        <tr valign="top">
            <td><code>push &lt;cmd&gt;</code></td>
            <td>Run non-interactive command and push result to chat history. If command is not provided, push the output
                of the recent invocation to chat history. The recent invocation output is present only when the command
                was completed with the zero exit code.</td>
        </tr>
        <tr valign="top">
            <td><code>pwd</code></td>
            <td>Print the current directory.</td>
        </tr>
        <tr valign="top">
            <td><code>reset</code></td>
            <td>Reset chat history.</td>
        </tr>
    </tbody>
</table>

These commands cannot be combined with other commands, either built-in or external, by using the pipe operator, and their output is not considered for push into the chat history.

## Providers

Developing custom tools for the Phi Shell involves creating a standalone program, referred to as the provider. This provider communicates with the Phi Shell, which acts as the tool host, using `Stdin` and `Stdout`. Any diagnostic output generated by the provider must be directed to `Stderr`, ensuring the Phi Shell can distinguish useful payload from diagnostic output.

To establish communication between the Phi Shell and the provider, the provider must follow a defined protocol. The steps in this protocol are:

1. **Listing Supported Tools**: The provider must print a list of the supported tools in JSON format, with each tool definition appearing on a separate line.

2. **Signaling Readiness**: After listing the tools, the provider must print an empty line to indicate that it is ready to process tool calls.

3. **Processing Tool Calls**: The provider must read input from `Stdin`, one line at a time, where each line represents a tool call. Once the tool call is processed, the provider must print the response to `Stdout`, with each response appearing on a separate line and preserving the call ID to ensure matching of requests and responses.

4. **Posting Asynchronous Messages**: The provider may optionally post asynchronous messages to `Stdout`, one message per line. These messages can represent events occurring outside the context of tool calls. Providers that do not implement tools but only post events are also valid.

Tool calls may be processed concurrently, and responses can be returned out of order. The Phi Shell applies timeouts while waiting for responses from the provider.

### Schemas

The communication contracts between the host and provider processes are outlined below as JSON schemas in the order of their typical usage by provider. The types can be converted to language-specific type definitions using variety of tools or by using LLMs.

<details>
<summary>Tool declaration</summary>

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "type": {
      "type": "string",
      "description": "The type of the tool. Currently, only `function` is supported.",
      "enum": ["function"]
    },
    "function": {
      "$ref": "#/definitions/ToolFunction",
      "description": "Description of the function."
    }
  },
  "required": ["type"],
  "definitions": {
    "ToolFunction": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "description": "The name of the function to be called. Must be a-z, A-Z, 0-9, or contain underscores and dashes, with a maximum length of 64.",
          "pattern": "^[a-zA-Z0-9_-]{1,64}$"
        },
        "description": {
          "type": "string",
          "description": "A description of what the function does, used by the model to choose when and how to call the function."
        },
        "parameters": {
          "type": "object",
          "additionalProperties": true,
          "description": "The parameters the functions accepts, described as a JSON Schema object."
        }
      },
      "required": ["name", "parameters"]
    }
  }
}
```
</details>

<details>
<summary>Tool request</summary>

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "call_id": {
      "type": "string",
      "description": "Correlation ID of the tool call."
    },
    "function": {
      "$ref": "#/definitions/ToolRequestFunction",
      "description": "The function that the model called."
    },
    "context": {
      "$ref": "#/definitions/ToolRequestContext",
      "description": "Context of the tool call request."
    }
  },
  "required": ["call_id", "function", "context"],
  "definitions": {
    "ToolRequestFunction": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "description": "Name of the function to call."
        },
        "arguments": {
          "type": "string",
          "description": "Arguments to call the function with, as generated by the model in JSON format."
        }
      },
      "required": ["name", "arguments"]
    },
    "ToolRequestContext": {
      "type": "object",
      "properties": {
        "dir": {
          "type": "string",
          "description": "Full path to the current directory."
        }
      },
      "required": ["dir"]
    }
  }
}
```
</details>

<details>
<summary>Tool response</summary>

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "call_id": {
      "type": "string",
      "description": "ID of the tool call that this message is responding to."
    },
    "content": {
      "type": "object",
      "description": "The contents of the tool response."
    }
  },
  "required": ["call_id"]
}
```
</details>

<details>
<summary>Standalone message</summary>

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "id": {
      "type": "string",
      "description": "Unique ID of the message used for deduplication."
    },
    "content": {
      "type": "string",
      "description": "Content of the message."
    },
    "dir": {
      "type": "string",
      "description": "Working directory to use in connection with the message."
    },
    "date": {
      "type": "string",
      "format": "date-time",
      "description": "Date and time when the message was created by provider. Can be different from the date and time when the message was actually sent to the host."
    }
  },
  "required": ["content"]
}
```
</details>

## Built-In Tools

Phi Shell offers several basic tools to use with daily tasks. These tools are available immediately upon the startup.

<table>
    <thead>
        <tr>
            <th align="left" width="175">Tool</th>
            <th align="left">Description</th>
        </tr>
    </thead>
    <tbody>
        <tr valign="top">
            <td><a href="cli/tool/builtin/schemas/exec_command.json">Execute Command</a></td>
            <td>Executes a console command or chain of piped commands</td>
        </tr>
        <tr valign="top">
            <td><a href="cli/tool/builtin/schemas/exec_http_call.json">Send HTTP Request</a></td>
            <td>Sends an HTTP request, given a URL, method, headers, query parameters and body of request</td>
        </tr>
        <tr valign="top">
            <td><a href="cli/tool/builtin/schemas/fs_create_update.json">Create Or Update File</a></td>
            <td>Creates or updates a file. When updating the file, the user is provided with diff between previous and current versions to approve or reject the changes.</td>
        </tr>
        <tr valign="top">
            <td><a href="cli/tool/builtin/schemas/fs_read.json">Read File</a></td>
            <td>Reads a file</td>
        </tr>
        <tr valign="top">
            <td><a href="cli/tool/builtin/schemas/fs_delete.json">Delete File Or Directory</a></td>
            <td>Deletes a file or directory. Can delete directories recursively if explicitly asked to.</td>
        </tr>
        <tr valign="top">
            <td><a href="cli/tool/builtin/schemas/fs_stat.json">Get File System Entry Info</a></td>
            <td>Gets information about a file system entry, whether it exists and what kind of entry it is.</td>
        </tr>
    </tbody>
</table>

## License

Phi Shell is licensed under the MIT License. See the [LICENSE.txt](LICENSE.txt) file for details.
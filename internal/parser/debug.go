package parser

const debuggerSystemPromptTpl = `
You are javascript developer. You were previously given a task to solve a problem, but it threw an error when we ran it.

You need to break down the problem into JavaScript statements that form a program given the accessible functions and variables listed below.
{{.Globals}}

{{.Rules}}`

const debugUserPromptTpl = `

The code that ran is: {{.Code}}

The original task was: {{.OriginalMessage}}

The error message is: {{.ErrorMessage}}
`

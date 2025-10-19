package ast

const LinguisticsSystemPrompt = "Extract the sentence information, this will be used to call tools or store structured information about the user message."

// Category represents the type of sentence
type Category string

const (
	CategoryInterrogative Category = "interrogative"
	CategoryDeclarative   Category = "declarative"
	CategoryImperative    Category = "imperative"
	CategoryOther         Category = "other"
)

// ToolDescription describes what tool could help in the context of the sentence
type ToolDescription struct {
	ToolDescription string   `json:"tool_description"` // Describe what tool could help in the context of the sentence
	ToolArguments   []string `json:"tool_arguments"`   // The string argument(s) you would pass to the tool, if any
}

// Sentence represents structured information about a sentence
type Sentence struct {
	Category            Category          `json:"category"`             // The category of the sentence
	Subject             string            `json:"subject"`              // The subject of the sentence
	Predicate           string            `json:"predicate"`            // The predicate of the sentence
	AtMention           string            `json:"at_mention"`           // The user mentioned with `@` in the sentence, if any
	HashTag             string            `json:"hash_tag"`             // The hash tag mentioned with `#` in the sentence, if any
	ContextDescriptions []string          `json:"context_descriptions"` // Describe the appropriate context of the sentence if we were to store it and reuse it later for another sentence, if any
	ToolDescriptions    []ToolDescription `json:"tool_descriptions"`    // Describe what tools could help in the context of the sentence
}

// Sentences represents a collection of sentence information
type Sentences struct {
	Sentences []Sentence `json:"sentences"`
}

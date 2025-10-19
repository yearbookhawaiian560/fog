package ast

/*
const SystemPrompt = "Parse the user message into a structured Abstract Syntax Tree (AST) for fog language parsing. Each node represents a semantic unit that can contain other nodes."

// NodeType represents the type of semantic node
type NodeType string

const (
	NodeTypeSentence            NodeType = "sentence"
	NodeTypeNounPhrase          NodeType = "noun_phrase"
	NodeTypeVerbPhrase          NodeType = "verb_phrase"
	NodeTypePrepositionalPhrase NodeType = "prepositional_phrase"
	NodeTypeDeterminer          NodeType = "determiner"
	NodeTypeNoun                NodeType = "noun"
	NodeTypeVerb                NodeType = "verb"
	NodeTypeAdjective           NodeType = "adjective"
	NodeTypeAdverb              NodeType = "adverb"
	NodeTypePreposition         NodeType = "preposition"
	NodeTypeConjunction         NodeType = "conjunction"
	NodeTypeComplementizer      NodeType = "complementizer"
)

// NLNode represents a node in an Abstract Syntax Tree for fog language parsing.
// Each node represents a semantic unit that can contain other nodes.
type NLNode struct {
	NodeType NodeType  `json:"node_type"` // Type of the semantic node
	Content  string    `json:"content"`   // Text content of the node
	Children []*NLNode `json:"children"`  // Child nodes
}

// ASTNodeType represents the type of AST node
type ASTNodeType string

const (
	NodeTypeProgram             ASTNodeType = "program"
	NodeTypeLetStatement        ASTNodeType = "let_statement"
	NodeTypeReturnStatement     ASTNodeType = "return_statement"
	NodeTypeExpressionStatement ASTNodeType = "expression_statement"
	NodeTypeBlockStatement      ASTNodeType = "block_statement"
	NodeTypeIdentifier          ASTNodeType = "identifier"
	NodeTypeIntegerLiteral      ASTNodeType = "integer_literal"
	NodeTypeBooleanLiteral      ASTNodeType = "boolean_literal"
	NodeTypeStringLiteral       ASTNodeType = "string_literal"
	NodeTypeArrayLiteral        ASTNodeType = "array_literal"
	NodeTypeHashLiteral         ASTNodeType = "hash_literal"
	NodeTypePrefixExpression    ASTNodeType = "prefix_expression"
	NodeTypeInfixExpression     ASTNodeType = "infix_expression"
	NodeTypeIfExpression        ASTNodeType = "if_expression"
	NodeTypeElseExpression      ASTNodeType = "else_expression"
)

// Node represents a node in an Abstract Syntax Tree for a simple interpreted language.
// Based on the Monkey programming language interpreter implementation.
type Node struct {
	NodeType    ASTNodeType `json:"node_type"`             // Type of AST node
	Value       string      `json:"value"`                 // Token literal value
	Operator    *string     `json:"operator,omitempty"`    // Operator for prefix/infix expressions
	Statements  []*Node     `json:"statements,omitempty"`  // List of statements in a block
	Condition   *Node       `json:"condition,omitempty"`   // Condition in if expression
	Consequence *Node       `json:"consequence,omitempty"` // Consequence block in if expression
	Alternative *Node       `json:"alternative,omitempty"` // Alternative block in if expression
	Parameters  []*Node     `json:"parameters,omitempty"`  // Function parameters
	Body        *Node       `json:"body,omitempty"`        // Function body
}
*/

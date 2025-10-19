package ast

type NodeType string

const (
	NodeTypeProgram         NodeType = "Program"
	NodeTypeLetStatement    NodeType = "LetStatement"
	NodeTypeIdentifier      NodeType = "Identifier"
	NodeTypeReturnStatement NodeType = "ReturnStatement"
	NodeTypeInfixExpression NodeType = "InfixExpression"
	NodeTypeFunctionLiteral NodeType = "FunctionLiteral"
	NodeTypeCallExpression  NodeType = "CallExpression"
	NodeTypeIntegerLiteral  NodeType = "IntegerLiteral"
	NodeTypeFloatLiteral    NodeType = "FloatLiteral"
	NodeTypeStringLiteral   NodeType = "StringLiteral"
	NodeTypeJSCode          NodeType = "JSCode"
)

func NodeTypesAsString() string {
	nodeTypes := []string{
		string(NodeTypeProgram),
		string(NodeTypeLetStatement),
		string(NodeTypeIdentifier),
		string(NodeTypeReturnStatement),
		string(NodeTypeInfixExpression),
		string(NodeTypeFunctionLiteral),
		string(NodeTypeCallExpression),
		string(NodeTypeIntegerLiteral),
		string(NodeTypeFloatLiteral),
		string(NodeTypeStringLiteral),
	}
	result := ""
	for _, nt := range nodeTypes {
		result += "- " + nt + "\n"
	}
	return result
}

type Node struct {
	Type        NodeType `json:"type"`                  // e.g. "Program", "LetStatement", "Identifier", etc.
	Name        string   `json:"name,omitempty"`        // For identifiers, variable names, etc.
	Description string   `json:"description,omitempty"` // The description of the node.
	JS          string   `json:"js,omitempty"`          // The JS code of the node.
	Distance    float64  `json:"distance,omitempty"`    // The distance of the object description from the embedding.
	Value       string   `json:"value,omitempty"`       // For literals, or the value of an assignment.
	Children    []*Node  `json:"children,omitempty"`    // For children nodes, can be arguments of a function call, or the body of a function.
}

/*
// Node is a generic AST node for LLM-friendly structured output.
type Node struct {
	Type        NodeType `json:"type"`                  // e.g. "Program", "LetStatement", "Identifier", etc.
	Name        string   `json:"name,omitempty"`        // For identifiers, variable names, etc.
	Description string   `json:"description,omitempty"` // The description of the node.
	Value       string   `json:"value,omitempty"`       // For literals, or the value of an assignment.
	NodeValue   *Node    `json:"node_value,omitempty"`  // For literals, or the value of an assignment.

	Arguments []*Node `json:"arguments,omitempty"` // For function call arguments.
	Body      []*Node `json:"body,omitempty"`      // For block statements, function bodies, etc.

	// For infix/prefix expressions.
	Operator string `json:"operator,omitempty"`

	// For binary expressions.
	Left  *Node `json:"left,omitempty"`  // For binary expressions.
	Right *Node `json:"right,omitempty"` // For binary expressions.

	// For conditional statements.
	Condition   *Node   `json:"condition,omitempty"`   // For if/while conditions.
	Consequence []*Node `json:"consequence,omitempty"` // For if/else blocks.
	Alternative []*Node `json:"alternative,omitempty"` // For else/else-if blocks.
}

// Example: Program that defines a function add(a, b), assigns result = add(2, 3),
// and uses an if statement to return "below four" if result < 4, else "above four".
var ProgramExample = &Node{
	Type: NodeTypeProgram,
	Body: []*Node{
		// function add(a, b) { return a + b; }
		{
			Type: NodeTypeFunctionLiteral,
			Name: "add",
			Arguments: []*Node{
				{Type: NodeTypeIdentifier, Name: "a"},
				{Type: NodeTypeIdentifier, Name: "b"},
			},
			Body: []*Node{
				{
					Type: NodeTypeReturnStatement,
					NodeValue: &Node{
						Type:     NodeTypeInfixExpression,
						Operator: "+",
						Left:     &Node{Type: NodeTypeIdentifier, Name: "a"},
						Right:    &Node{Type: NodeTypeIdentifier, Name: "b"},
					},
				},
			},
		},
		// let result = add(2, 3);
		{
			Type: NodeTypeLetStatement,
			Name: "result",
			NodeValue: &Node{
				Type: NodeTypeCallExpression,
				Name: "add",
				Arguments: []*Node{
					{Type: NodeTypeIdentifier, Value: "2"},
					{Type: NodeTypeIdentifier, Value: "3"},
				},
			},
		},
		// if (result < 4) { return "below four"; } else { return "above four"; }
		{
			Type: "IfStatement",
			Condition: &Node{
				Type:     NodeTypeInfixExpression,
				Operator: "<",
				Left:     &Node{Type: NodeTypeIdentifier, Name: "result"},
				Right:    &Node{Type: NodeTypeIdentifier, Value: "4"},
			},
			Consequence: []*Node{
				{
					Type:      NodeTypeReturnStatement,
					NodeValue: &Node{Type: NodeTypeStringLiteral, Value: "below four"},
				},
			},
			Alternative: []*Node{
				{
					Type:      NodeTypeReturnStatement,
					NodeValue: &Node{Type: NodeTypeStringLiteral, Value: "above four"},
				},
			},
		},
	},
}
*/

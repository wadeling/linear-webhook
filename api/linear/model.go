package linear

var QueryIssue = `
    query ($key: String!) {
        issue (id:$key) {
			id
            title 
			assignee {
				displayName
			}
			state {
				name
			}
			updatedAt
			history {
				nodes {
					toState {
						name
						type
					}
					toAssignee {
						displayName
					}
				}
			}
        }
    }
`

// User define linear user
type User struct {
	Name        string `json:"name"`        // full name,eg. user@email.com
	DisplayName string `json:"displayName"` // username without email
}

// WorkFlowState define linear workflow state,eg.'canceled','progressing'
type WorkFlowState struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	UpdatedAt string `json:"updatedAt"`
}

// Node history record
type Node struct {
	ToAssignee User          `json:"toAssignee"`
	ToState    WorkFlowState `json:"toState"`
}

// History define issue change history
type History struct {
	Nodes []Node `json:"nodes"`
}

type Issue struct {
	Id        string        `json:"id"`
	Title     string        `json:"title"`
	Assignee  User          `json:"assignee"`
	State     WorkFlowState `json:"state"`
	History   History       `json:"history"`
	UpdatedAt string        `json:"updatedAt"`
}

type IssueInfo struct {
	Issue Issue `json:"issue"`
}

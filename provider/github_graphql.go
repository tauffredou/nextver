package provider

type CommitNode struct {
	Oid     string
	Message string
}

type PageInfo struct {
	HasNextPage bool
}

type tagEdge struct {
	Node struct {
		Target struct {
			Tag struct {
				Message string
				Target  struct {
					Commit struct {
						Oid string
					} `graphql:"... on Commit"`
				}
			} `graphql:"... on Tag"`
		}
	}
}

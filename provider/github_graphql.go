package provider

import "time"

type CommitNode struct {
	Oid     string
	Message string
	Author  Author
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

type Author struct {
	Name  string
	Email string
	Date  time.Time
}

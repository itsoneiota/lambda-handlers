package aws

import "context"

// Context is the aws request context.
type Context struct {
	context.Context
	sourceIP string
}

// SourceIP returns the source ip that has made the request.
func (c Context) SourceIP() string {
	return c.sourceIP
}

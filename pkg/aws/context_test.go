package aws

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestContexterInterface(t *testing.T) {
	ip := "SourceIP"
	unixNow := int64(1234567)
	e := events.APIGatewayProxyRequestContext{
		Identity: events.APIGatewayRequestIdentity{
			SourceIP: ip,
		},
		RequestTimeEpoch: unixNow,
	}

	ctx := Context{e}

	assert.Equal(t, ip, ctx.SourceIP())
	assert.Equal(t, unixNow, ctx.UnixNow())
}

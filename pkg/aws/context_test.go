package aws

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestContexterInterface(t *testing.T) {
	ip := "SourceIP"
	stage := "dev"
	unixNow := int64(1234567)
	e := events.APIGatewayProxyRequestContext{
		Identity: events.APIGatewayRequestIdentity{
			SourceIP: ip,
		},
		RequestTimeEpoch: unixNow,
		HTTPMethod:       http.MethodGet,
		Stage:            stage,
	}

	ctx := Context{e}

	assert.Equal(t, ip, ctx.SourceIP())
	assert.Equal(t, unixNow, ctx.UnixNow())
	assert.Equal(t, http.MethodGet, ctx.HttpMethod())
	assert.Equal(t, stage, ctx.Stage())
}

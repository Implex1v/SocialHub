package integrationtest

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
	IntegrationTest(t, func() {
		resp, err := http.Get("http://localhost:8000/health")
		if err != nil {
			t.Error("failed to get status", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "{\"Status\":\"Ok\"}\n", bodyToString(t, resp.Body))
	})
}

func bodyToString(t *testing.T, r io.ReadCloser) string {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		t.Fatal(err)
	}

	return buf.String()
}

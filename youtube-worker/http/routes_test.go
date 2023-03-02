package http

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	e := echo.New()
	Status(e)
	req := httptest.NewRequest(http.MethodGet, "/status", strings.NewReader(""))
	rec := httptest.NewRecorder()
	e.NewContext(req, rec)

	assert.Equal(t, "Ok", rec.Body.String())
}

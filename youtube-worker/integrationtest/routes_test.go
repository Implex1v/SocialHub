package integrationtest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

type RoutesTestSuite struct {
	suite.Suite
	wrapper Wrapper
}

func TestRoutes(t *testing.T) {
	suite.Run(t, &RoutesTestSuite{})
}

func (s *RoutesTestSuite) SetupSuite() {
	if err := s.wrapper.RunContainer(); err != nil {
		log.Fatalln(err)
	}
}

func (s *RoutesTestSuite) TearDownSuite() {
	s.wrapper.CleanUp()
}

func (s *RoutesTestSuite) Test_Status_Endpoint() {
	resp, err := http.Get("http://localhost:" + s.wrapper.httpPort + "/health")
	if err != nil {
		s.T().Error("failed to get status", err)
	}

	assert.Equal(s.T(), 200, resp.StatusCode)
	assert.Equal(s.T(), "{\"Status\":\"Ok\"}\n", bodyToString(s.T(), resp.Body))
}

func (s *RoutesTestSuite) Test_Metrics_Endpoint() {
	resp, err := http.Get("http://localhost:" + s.wrapper.httpPort + "/metrics")
	if err != nil {
		s.T().Error("failed to get status", err)
	}

	assert.Equal(s.T(), 200, resp.StatusCode)
}

func bodyToString(t *testing.T, r io.ReadCloser) string {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		t.Fatal(err)
	}

	return buf.String()
}

package integrationtest

import (
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type ConsumerTestSuite struct {
	suite.Suite
	wrapper Wrapper
}

func TestConsumer(t *testing.T) {
	suite.Run(t, &ConsumerTestSuite{})
}

func (s *ConsumerTestSuite) SetupSuite() {
	if err := s.wrapper.RunContainer(); err != nil {
		log.Fatalln(err)
	}
}

func (s *ConsumerTestSuite) TearDownSuite() {
	s.wrapper.CleanUp()
}

func (s *ConsumerTestSuite) Test_Youtube_Poll() {
	s.wrapper.NewMessage(s.T(), "youtube.poll", "foo")
}

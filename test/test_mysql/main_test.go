package test_mysql

import (
	"github.com/andrewapj/redwing/test"
	"os"
	"testing"
)

func TestMain(m *testing.M)  {

	identifier := test.StartContainer([]string{"../../db/mysql/docker-compose.yml"})
	exitVal := m.Run()
	test.StopContainer([]string{"../../db/mysql/docker-compose.yml"}, identifier)

	os.Exit(exitVal)
}

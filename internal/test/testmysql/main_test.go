package testmysql

import (
	"github.com/andrewapj/redwing/internal/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	identifier := test.StartContainer([]string{"../../../db/mysql/docker-compose.yml"})
	exitVal := m.Run()
	test.StopContainer([]string{"../../../db/mysql/docker-compose.yml"}, identifier)

	os.Exit(exitVal)
}

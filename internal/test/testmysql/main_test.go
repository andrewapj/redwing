package testmysql

import (
	"github.com/andrewapj/redwing/internal"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	identifier := internal.StartContainer([]string{"../../../db/mysql/docker-compose.yml"})
	exitVal := m.Run()
	internal.StopContainer([]string{"../../../db/mysql/docker-compose.yml"}, identifier)

	os.Exit(exitVal)
}

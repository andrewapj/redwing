package test

import (
	"fmt"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	"os"
	"strings"
	"time"
)

func StartContainer(compose []string) string {
	identifier := strings.ToLower(uuid.New().String())
	c := tc.NewLocalDockerCompose(compose, identifier)

	execError := c.WithCommand([]string{"up", "-d"}).Invoke()
	time.Sleep(time.Second * 20)

	if execError.Error != nil {
		fmt.Printf("Error starting container: %s : %v", compose, execError.Error)
		os.Exit(1)
	}
	return identifier
}

func StopContainer(compose []string, identifier string) {
	c := tc.NewLocalDockerCompose(compose, identifier)
	execError := c.Down()
	if execError.Error != nil {
		fmt.Printf("Error starting container: %s : %v", compose, execError.Error)
		os.Exit(1)
	}
}

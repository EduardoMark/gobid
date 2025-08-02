package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	ctx := context.TODO()

	cmd := exec.CommandContext(
		ctx,
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":    err,
			"output": string(output),
		}).Error("Failed to execute command")
		return
	}

	logrus.WithField(
		"output", string(output),
	).Info("Migration completed successfully.")
}

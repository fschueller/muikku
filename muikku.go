package main

import (
  "os"
  "log"
  "github.com/urfave/cli/v2"
)

func main() {
  app := &cli.App{
    Name: "muikku",
    Usage: "CLI for basic photo management",
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}

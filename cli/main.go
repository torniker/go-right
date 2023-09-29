package main

import (
	"github.com/spf13/cobra"
	"github.com/torniker/go-right/cli/country"
	"github.com/torniker/go-right/env"
	"github.com/torniker/go-right/pkg/repeat"
)

func main() {
	env.New("cli")
	var rootCmd = &cobra.Command{Use: "right"}
	rootCmd.PersistentFlags().StringVarP(&repeat.Repeat, "repeat", "r", "", "repeat command")

	country.Register(rootCmd)

	_ = rootCmd.Execute()
}


package main

import (
	"context"
	"os"

	"sewz/core"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "sewz",
		Short: "A pretty docker ps wrapper",
		RunE:  core.RunPs,
	}

	cmd.Flags().BoolVarP(&core.ShowAll, "all", "a", false, "Show all containers")
	cmd.Flags().BoolVar(&core.ShowID, "id", false, "Show container ID")
	cmd.Flags().BoolVarP(&core.ShowImage, "image", "i", false, "Show image")
	cmd.Flags().BoolVar(&core.ShowCmd, "cmd", false, "Show command")
	cmd.Flags().BoolVar(&core.ShowCreated, "created", false, "Show created date")

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}

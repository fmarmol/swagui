package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fmarmol/jin"
	"github.com/fmarmol/swagui/templates"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use: "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("file argument expected before flags")
		}
		return nil

	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing file argument")
		}
		port := cmd.Flag("port").Value.String()
		token := cmd.Flag("token").Value.String()
		filename := args[0]
		docFd, err := os.Open(filename)
		if err != nil {
			return err
		}

		docBytes, err := io.ReadAll(docFd)
		if err != nil {
			return err
		}

		jin.SetRealeaseMode()
		app := jin.New()
		app.GET("/docs", func(c jin.Context) (any, error) {
			c.Data(200, "application/yaml", docBytes)
			return nil, nil

		})
		app.GET("/", func(c jin.Context) (any, error) {
			return c.Templ(templates.Index("/docs", token))
		})
		log.Println("http://localhost:" + port)
		return app.Run(":" + port)
	},
}

func main() {
	cmd.Flags().StringP("port", "p", "4242", "port")
	cmd.Flags().StringP("token", "t", "", "jwt token")
	err := cmd.Execute()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

}

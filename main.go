package main

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/fmarmol/jin"
	"github.com/spf13/cobra"
)

//go:embed templates/index.template
var indexTmpl string

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
		filename := args[0]
		docFd, err := os.Open(filename)
		if err != nil {
			return err
		}

		docBytes, err := io.ReadAll(docFd)
		if err != nil {
			return err
		}

		t := template.New("index")
		t, err = t.Parse(indexTmpl)
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
			err := t.Execute(c.Writer, map[string]any{"URL": "/docs"})
			if err != nil {
				return nil, err
			}
			return nil, nil
		})
		log.Println("http://localhost:" + port)
		return app.Run(":" + port)
	},
}

func main() {
	cmd.Flags().StringP("port", "p", "4242", "port")
	err := cmd.Execute()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

}

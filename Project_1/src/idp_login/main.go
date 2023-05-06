package main

import (
    "fmt"
    "log"
    "os"

    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "idp_login",
        Usage: "Manage Identity Providers (IdPs) and identity attributes for users",
        
		Commands: []*cli.Command{
			{
				Name:    "manage-idp",
				Usage:       "Set, change or delete IdPs and their operational parameters",
				Description: "For host administrators",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "idp",
						Aliases: []string{"i"},
						Usage:   "IdP name",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "params",
						Aliases: []string{"p"},
						Usage:   "IdP operational parameters for the IdP",
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("manage-idp")
					return nil
				},
			},
			{
				Name:    "manage-attributes",
				Usage:   "Set, change or delete identity attributes for a given IdP for the current user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "idp",
						Aliases: []string{"i"},
						Usage:   "IdP name",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "attributes",
						Aliases: []string{"a"},
						Usage:   "Identity attributes for the IdP",
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("manage-attributes")
					return nil
				},
			},
			{
				Name:    "list-users",
				Usage:   "List the users registered for the current IdP and the identity parameters registered for each user",
				Action: func(c *cli.Context) error {
					fmt.Println("list-users")
					return nil
				},
			},
			{
				Name:    "list-idps",
				Usage:   "List the IdPs registered for the current user and the identity parameters registered for each IdP",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "user",
						Aliases: []string{"u"},
						Usage:   "User to list IdPs for (defaults to current user)",
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("list-idps")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

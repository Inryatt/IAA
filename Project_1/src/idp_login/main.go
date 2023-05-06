package main

import (
    "fmt"
    "log"
    "os"
	"os/user"
	"database/sql"

    "github.com/urfave/cli/v2"
	_ "github.com/mattn/go-sqlite3"
)

func isAdministrator() bool {
	// Check if the current user is the idpadmins group
	usr, err := user.Current()
	if err != nil {
		fmt.Println("[!] Error: could not get current user")
		return false
	}

	// Get group information
	group, err := user.LookupGroup("idpadmins")
	if err != nil {
		fmt.Println("[!] Error: could not get idpadmins group information")
		return false
	}

	// Get the list of user groups
	userGroups, err := usr.GroupIds()
	if err != nil {
		fmt.Println("[!] Error: could not get user groups")
		return false
	}

	// Check if the user is in the idpadmins group
	for _, userGroup := range userGroups {
		if userGroup == group.Gid {
			return true
		}
	}

	return false
}

func check_requirements() {
	// Check if the program has the setuid bit set
	if os.Geteuid() != 0 {
		fmt.Println("[!] Error: program must be run with the 4750 permissions")
		return
	}

	// Check if database exists
	if _, err := os.Stat("/etc/project_1.db"); os.IsNotExist(err) {
		fmt.Println("[-] Database does not exist, creating it...")

		// Create database
		db, err := sql.Open("sqlite3", "/etc/project_1.db")
		if err != nil {
			fmt.Println("[!] Error: could not create database")
			return
		}
		defer db.Close()

		// Create users table
		_, err = db.Exec("CREATE TABLE users (username TEXT PRIMARY KEY, idps TEXT)")
		if err != nil {
			fmt.Println("[!] Error: could not create users table")
			return
		}

		// Create idps table
		_, err = db.Exec("CREATE TABLE idps (name TEXT PRIMARY KEY, params TEXT)")
		if err != nil {
			fmt.Println("[!] Error: could not create idps table")
			return
		}

		// Create attributes table (username, idp, attributes)
		_, err = db.Exec("CREATE TABLE attributes (username TEXT, idp TEXT, attributes TEXT, PRIMARY KEY (username, idp))")
		if err != nil {
			fmt.Println("[!] Error: could not create attributes table")
			return
		}
	}
}

func main() {
	// Check requirements
	check_requirements()

    app := &cli.App{
        Name:  "idp_login",
        Usage: "Manage Identity Providers (IdPs) and identity attributes for users",
        
		Commands: []*cli.Command{
			{
				Name:    	 "manage-idp",
				Usage:       "manage-idp [--operation set|change|delete] [--idp IDP_NAME] [--params PARAMS]",
				Description: "Set, change or delete operational parameters for a given IdP, only users belonging to the idpadmins group can perform this operation",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "operation",
						Aliases: []string{"o"},
						Usage:   "Operation to perform (set, change, delete)",
						Required: true,
					},
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
					// Check if the current user is the idpadmins group
					if !isAdministrator() {
						fmt.Println("[!] Error: current user is not an administrator")
						return nil
					}

					fmt.Println("manage-idp")
					return nil
				},
			},
			{
				Name:    	 "manage-attributes",
				Usage:   	 "manage-attributes [--operation set|change|delete] [--idp IDP_NAME] [--attributes ATTRIBUTES]",
				Description: "Set, change or delete identity attributes for a given IdP, the changes are applied only to the current user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "operation",
						Aliases: []string{"o"},
						Usage:   "Operation to perform (set, change, delete)",
						Required: true,
					},
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
				Name:    	 "list-users",
				Usage:   	 "list-users",
				Description: "List all users with registered IdPs, only users belonging to the idpadmins group can perform this operation",
				Action: func(c *cli.Context) error {
					// Check if the current user is the idpadmins group
					if !isAdministrator() {
						fmt.Println("[!] Error: current user is not an administrator")
						return nil
					}

					fmt.Println("list-users")
					return nil
				},
			},
			{
				Name:    	 "list-idps",
				Usage:   	 "list-idps",
				Description: "List all registered IdPs, only for the current user",
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

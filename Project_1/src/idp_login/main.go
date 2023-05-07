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

const DATABASE_PATH = "/tmp/project_1.sqlite"
//const DATABASE_PATH = "/etc/project_1.sqlite"

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
	// Print requirements message
	fmt.Println("[*] Checking requirements...")
	/*
	// Check if the program has the setuid bit set
	if os.Geteuid() != 0 {
		fmt.Println("[!] Error: program must be run with the 4750 permissions")
		return
	}
	*/

	// Check if database exists
	if _, err := os.Stat(DATABASE_PATH); os.IsNotExist(err) {
		fmt.Println("[-] Database does not exist, creating it...")

		// Create database
		db, err := sql.Open("sqlite3", DATABASE_PATH)
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

	// Print success message
	fmt.Println("[+] All requirements met")
}

func idp_exists(idp string) bool {
	// Open database connection
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return false
	}

	// Close database connection
	defer db.Close()

	// Check if the IdP exists
	rows, err := db.Query("SELECT name FROM idps WHERE name = ?", idp)
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return false
	}

	// Close rows
	defer rows.Close()

	// Check if the IdP exists
	if rows.Next() {
		return true
	}

	return false
}

func manage_idp(operation string, idp string, params string) {
	// If the operation is set, check if the IdP already exists
	if operation == "set" {
		// Check if params is empty
		if params == "" {
			fmt.Println("[!] Error: params cannot be empty for set operation")
			return
		}

		// Check if the IdP already exists
		if idp_exists(idp) {
			fmt.Println("[!] Error: IdP already exists")
			return
		}

		// Open database connection
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		// Close database connection
		defer db.Close()

		// Insert the IdP into the database
		_, err = db.Exec("INSERT INTO idps (name, params) VALUES (?, ?)", idp, params)
		if err != nil {
			fmt.Println("[!] Error: could not insert IdP into database")
			return
		}

		// Print success message
		fmt.Println("[+] IdP successfully added")
	} else if operation == "change" {
		// Check if params is empty
		if params == "" {
			fmt.Println("[!] Error: params cannot be empty for change operation")
			return
		}
		// Check if the IdP exists
		if !idp_exists(idp) {
			fmt.Println("[!] Error: IdP does not exist")
			return
		}

		// Open database connection
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		// Close database connection
		defer db.Close()

		// Update the IdP in the sqlite database
		_, err = db.Exec("UPDATE idps SET params = ? WHERE name = ?", params, idp)
		if err != nil {
			// print error
			fmt.Println(err)
			return
		}

		// Print success message
		fmt.Println("[+] IdP successfully updated")
	} else if operation == "delete" {
		// Check if the IdP exists
		if !idp_exists(idp) {
			fmt.Println("[!] Error: IdP does not exist")
			return
		}

		// Open database connection
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		// Close database connection
		defer db.Close()

		// Delete the IdP from the database
		_, err = db.Exec("DELETE FROM idps WHERE name = ?", idp)
		if err != nil {
			fmt.Println("[!] Error: could not delete IdP from database")
			return
		}

		// Print success message
		fmt.Println("[+] IdP successfully deleted")
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
					/*
					if !isAdministrator() {
						fmt.Println("[!] Error: current user is not an administrator")
						return nil
					}
					*/

					// Get the operation to perform
					operation := c.String("operation")
					if operation != "set" && operation != "change" && operation != "delete" {
						fmt.Println("[!] Error: invalid operation")
					} else {
						// Get the IdP name
						idp := c.String("idp")
						
						// Get the IdP params
						params := c.String("params")

						// Manage the IdP
						manage_idp(operation, idp, params)

					}
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

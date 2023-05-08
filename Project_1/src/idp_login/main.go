package main

import (
    "fmt"
    "log"
    "os"
	"os/user"
	"database/sql"
    "encoding/json"
	"bytes"

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
			fmt.Println("[!] Error: could not update IdP in database")
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

func getCurrentUser() string {
	// Get the current user
	user, err := user.Current()
	if err != nil {
		fmt.Println("[!] Error: could not get current user")
		return ""
	}

	return user.Username
}

func list_available_idps() {
	// Open database connection
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return
	}

	// Close database connection
	defer db.Close()

	// Query the database
	rows, err := db.Query("SELECT name FROM idps")
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	// Close rows
	defer rows.Close()

	// Print IdPs
	fmt.Println("[+] Available IdPs:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println("    - " + name)
	}
}

func list_users() {
	// Open database connection
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return
	}

	// Close database connection
	defer db.Close()

	// Query the database
	rows, err := db.Query("SELECT username FROM attributes")
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	// Close rows
	defer rows.Close()

	// Print users
	fmt.Println("[+] Users:")
	for rows.Next() {
		var username string
		rows.Scan(&username)
		fmt.Println("    - " + username)
	}
}

func list_idps(username string) {
	// Open database connection
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return
	}

	// Close database connection
	defer db.Close()

	// Query the database for the user's IdPs and attributes
	rows, err := db.Query("SELECT idp, attributes FROM attributes WHERE username = ?", username)
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	// Close rows
	defer rows.Close()

	// Print IdPs and pretty print attributes
	for rows.Next() {
		var idp string
		var attributes string
		var out bytes.Buffer

		rows.Scan(&idp, &attributes)

		err := json.Indent(&out, []byte(attributes), "", "\t")
		if err != nil {
			fmt.Println("[!] Error: could not indent JSON")
			fmt.Println("[+] IdP name: " + idp)
			fmt.Println("[+] IdP params:")
			fmt.Println(attributes)
			return
		}

		fmt.Println("[+] IdP name: " + idp)
		fmt.Println("[+] IdP attributes:")
		fmt.Println(out.String())
	}
}

func print_attributes(idp string) {
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

	// Query the database abd query for the idp params
	rows, err := db.Query("SELECT params FROM idps WHERE name = ?", idp)
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	// Close rows
	defer rows.Close()

	// Print JSON params
	for rows.Next() {
		var params string
		rows.Scan(&params)
		var out bytes.Buffer
		err := json.Indent(&out, []byte(params), "", "\t")
		if err != nil {
			fmt.Println("[!] Error: could not indent JSON")
			fmt.Println("[+] IdP params:")
			fmt.Println(params)
			return
		}

		fmt.Println("[+] IdP params:")
		fmt.Println(out.String())
	}
}

func manage_attributes(username string, operation string, idp string, attributes string) {
	// Check if the IdP exists
	if !idp_exists(idp) {
		fmt.Println("[!] Error: IdP does not exist")
		return
	}

	// Check if attributes is empty
	if attributes == "" && operation != "delete" {
		fmt.Println("[!] Error: attributes cannot be empty")
		return
	}

	// Check if the operation is valid
	if operation != "set" && operation != "change" && operation != "delete" {
		fmt.Println("[!] Error: invalid operation")
		return
	}

	if operation == "set" {
		// Open database connection
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		// Close database connection
		defer db.Close()

		// Insert the attributes into the database
		_, err = db.Exec("INSERT INTO attributes (username, idp, attributes) VALUES (?, ?, ?)", username, idp, attributes)
		if err != nil {
			fmt.Println("[!] Error: could not insert attributes into database")
			return
		}

		// Print success message
		fmt.Println("[+] Attributes successfully added")
	} else if operation == "change" {
		// Check if attributes is empty
		if attributes == "" {
			fmt.Println("[!] Error: attributes cannot be empty for change operation")
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

		// Update the attributes in the sqlite database
		_, err = db.Exec("UPDATE attributes SET attributes = ? WHERE username = ? AND idp = ?", attributes, username, idp)
		if err != nil {
			fmt.Println("[!] Error: could not update attributes in database")
			return
		}

		// Print success message
		fmt.Println("[+] Attributes successfully updated")
	} else if operation == "delete" {
		// Open database connection
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		// Close database connection
		defer db.Close()

		// Delete the attributes from the sqlite database
		_, err = db.Exec("DELETE FROM attributes WHERE username = ? AND idp = ?", username, idp)
		if err != nil {
			fmt.Println("[!] Error: could not delete attributes from database")
			return
		}

		// Print success message
		fmt.Println("[+] Attributes successfully deleted")
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
					},
					&cli.StringFlag{
						Name:    "attributes",
						Aliases: []string{"a"},
						Usage:   "Identity attributes for the IdP",
					},
				},
				Action: func(c *cli.Context) error {
					// Get the operation to perform
					operation := c.String("operation")
					if operation != "set" && operation != "change" && operation != "delete" {
						fmt.Println("[!] Error: invalid operation")
					} else {
						// Get username
						username := getCurrentUser()

						// Check if idp is set
						idp := c.String("idp")
						if idp == "" {
							// List all IdPs available
							list_available_idps()
						} else {
							// Get the attributes
							attributes := c.String("attributes")

							if attributes == "" && operation != "delete" {
								// Print the attributes for the given IdP
								print_attributes(idp)
							} else {
								// Manage the attributes
								manage_attributes(username, operation, idp, attributes)
							}
						}
					}

					return nil
				},
			},
			{
				Name:    	 "list-users",
				Usage:   	 "list-users",
				Description: "List all users with registered IdPs, only users belonging to the idpadmins group can perform this operation",
				Action: func(c *cli.Context) error {
					// Check if the current user is the idpadmins group
					/*
					if !isAdministrator() {
						fmt.Println("[!] Error: current user is not an administrator")
						return nil
					}
					*/

					// List users
					list_users()

					return nil
				},
			},
			{
				Name:    	 "list-idps",
				Usage:   	 "list-idps",
				Description: "List all registered IdPs, only for the current user",
				Action: func(c *cli.Context) error {
					// Get the current user
					username := getCurrentUser()

					// List IdPs
					list_idps(username)

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
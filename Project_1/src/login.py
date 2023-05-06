import site

site.main()

import os
import qrcode
import sqlite3

def get_idps():
    DB_LOCATION = "/etc/project_1.db"

    # Check if the database exists
    if not os.path.exists(DB_LOCATION):
        return None

    # Connect to the database
    conn = sqlite3.connect(DB_LOCATION)
    c = conn.cursor()

    # Get the IDPs
    c.execute("SELECT * FROM idps")
    idps = c.fetchall()

    # Close the connection
    conn.close()

    return idps

def qr_code(url):
    # Create a QR code instance
    qr = qrcode.QRCode(version=1, box_size=1, border=1)

    # Add the data to the QR code
    qr.add_data(url)

    # Compile the QR code
    qr.make(fit=True)

    # Get the QR code data as a list of lists
    data = qr.get_matrix()

    # Print the URL
    print("\033[1;33m[?]\033[0m Visit the following URL to authenticate: " + url)

    # Print the QR code to the terminal with colors and spacing
    print("\033[33m[-]\033[0m Or scan the following QR code:")
    for row in data:
        for cell in row:
            if cell:
                print "  ",
            else:
                print "##",
        print

def auth(username):
    if username != None:
        # print with colors
        print("\033[1;32m[+]\033[0m Authenticating user: " + username)
    
    # Get the IDPs
    idps = get_idps()

    # Check if there are any IDPs
    if idps == None:
        print("\033[1;31m[!]\033[0m No IDPs configured")
        # print with orange color
        print("\033[33m[-]\033[0m Falling back to local authentication")
        return False

    # Print the menu
    while True:
        print("\033[1;33m[?]\033[0m Select an IDP:")
        for i in range(len(idps)):
            print(str(i) + ") " + idps[i][0])
        print("q) Quit")

        # Get the user's choice
        choice = raw_input("> ")

        # Check if the user wants to quit
        if choice == "q":
            return False

        # Check if the user's choice is valid
        try:
            choice = int(choice)
            if choice < 0 or choice >= len(idps):
                raise ValueError
        except ValueError:
            print("[-] Invalid choice")
            continue

        # Get the IDP
        idp = idps[choice]

        # Print the IDP's information
        print("\033[1;32m[+]\033[0m IDP: " + idp[0])
        break

    return True

def get_user(pamh):
    try:
        return pamh.get_user(None)
    except pamh.exception, e:
        return e.pam_result

def pam_sm_authenticate(pamh, flags, argv):
    user = get_user(pamh)
    if user == None:
        return pamh.PAM_USER_UNKNOWN

    try:
        if auth(user) == True:
            return pamh.PAM_SUCCESS
        else:
            return pamh.PAM_AUTH_ERR
    except pamh.exception, e:
        return pamh.PAM_AUTH_ERR

def pam_sm_open_session(pamh, flags, argv):
    user = get_user(pamh)

    if user == None:
        return pamh.PAM_USER_UNKNOWN

    home_dir = pathlib.Path("/home/" + user)

    if not home_dir.exists():
        home_dir.mkdir()

    return pamh.PAM_SUCCESS

def pam_sm_close_session(pamh, flags, argv):
    return pamh.PAM_SUCCESS

def pam_sm_setcred(pamh, flags, argv):
    return pamh.PAM_SUCCESS

def pam_sm_acct_mgmt(pamh, flags, argv):  
    return pamh.PAM_SUCCESS

def pam_sm_chauthtok(pamh, flags, argv):
    return pamh.PAM_SUCCESS
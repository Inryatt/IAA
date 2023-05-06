import site

site.main()

import os
import qrcode
import sqlite3

def get_idps():
    DB_LOCATION = "/etc/project_1.db"

    if not os.path.exists(DB_LOCATION):
        return None

    conn = sqlite3.connect(DB_LOCATION)
    c = conn.cursor()

    c.execute("SELECT * FROM idps")
    idps = c.fetchall()

    conn.close()

    return idps

def qr_code(url):
    qr = qrcode.QRCode(version=1, box_size=1, border=1)

    qr.add_data(url)

    qr.make(fit=True)

    data = qr.get_matrix()

    print("\033[1;33m[?]\033[0m Visit the following URL to authenticate: " + url)

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
        print("\033[1;32m[+]\033[0m Authenticating user: " + username)
    
    idps = get_idps()

    if idps == None:
        print("\033[1;31m[!]\033[0m No IDPs configured")
        print("\033[33m[-]\033[0m Falling back to local authentication")
        return False

    while True:
        print("\033[1;33m[?]\033[0m Select an IDP:")
        for i in range(len(idps)):
            print(str(i) + ") " + idps[i][0])
        print("q) Quit")

        choice = raw_input("> ")

        if choice == "q":
            return False

        try:
            choice = int(choice)
            if choice < 0 or choice >= len(idps):
                raise ValueError
        except ValueError:
            print("[-] Invalid choice")
            continue

        idp = idps[choice]

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
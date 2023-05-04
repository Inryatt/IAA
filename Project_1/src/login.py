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
        resp = pamh.conversation(pamh.Message(pamh.PAM_PROMPT_ECHO_OFF, "Password:"))
        if resp.resp == "password":
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



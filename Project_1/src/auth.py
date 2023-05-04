import sqlite3

DB_LOCATION = "/etc/project_1.db"

def db_exists():
    try:
        conn = sqlite3.connect(DB_LOCATION)
        conn.close()
        return True
    except:
        return False

def create_db():
    conn = sqlite3.connect(DB_LOCATION)
    c = conn.cursor()
    c.execute("CREATE TABLE users (username text, password text)")
    c.execute("CREATE TABLE tokens (username text, token text)")
    conn.commit()
    conn.close()

def check_user_token(username, token):
    conn = sqlite3.connect(DB_LOCATION)
    c = conn.cursor()
    c.execute("SELECT * FROM tokens WHERE username=? AND token=?", (username, token))
    if c.fetchone() == None:
        conn.close()
        return False
    else:
        conn.close()
        return True

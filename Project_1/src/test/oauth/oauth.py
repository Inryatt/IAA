import requests
import time
import qrcode
import sqlite3
import json

# Get available idps for current user
def get_idps(username, DATABASE_PATH):
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()

    c.execute("SELECT idp FROM attributes WHERE username = ?", (username,))
    idps = c.fetchall()

    conn.close()
    return idps

# Get idp information
def get_idp(idp, username, DATABASE_PATH):
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()

    c.execute("SELECT attributes FROM attributes WHERE username = ? AND idp = ?", (username, idp))
    idp = c.fetchall()

    conn.close()

    # Get idp json
    try:
        python_dict = json.loads(idp[0][0])
        idp = {
            str(key) if isinstance(key, unicode) else key:
            str(value) if isinstance(value, unicode) else value
            for key, value in python_dict.items()
        }
    except:
        print('\033[1;31m[!]\033[0m Invalid IdP')
        return None

    return idp


# Parse API response
def parse_response(response, dict):
    response = response.replace(' ' , '') \
                       .replace('\n', '') \
                       .replace('}' , '') \
                       .replace('{' , '') \
                       .replace('":' , '=') \
                       .replace('"' , '') \
                       .replace(',' , '&') 

    for pair in response.split('&'):
        key, value = pair.split('=')
        dict[key] = value
    return dict

# Request a device code
def request_device(client_id, scope, url):
    data = {
        'client_id': client_id,
        'scope': scope
    }

    response = requests.post(url, data=data)

    if response.status_code == 200:
        response_dict = parse_response(response.text, {})
        return (response_dict['device_code'], response_dict['user_code'])
    else:
        # Print error message
        print('\033[1;31m[!]\033[0m Error requesting device code')
        return None

# Generate QR code and print to console
def generate_qr_code(url):
    qr = qrcode.QRCode(version=1, box_size=10, border=5)
    qr.add_data(url)
    qr.make(fit=True)
    qr.print_ascii()


# Poll for a user token
def poll_for_token(url, arguments):
    interval = 5
    timeout = 60

    start_time = time.time()
    while True:
        response = requests.post(url, data=arguments)
        if response.status_code == 200:
            response_dict = parse_response(response.text, {})
            
            if 'access_token' in response_dict:
                return response_dict
            else:
                time.sleep(interval)
        else:
            return None

def oauth2(idp, qr_code):
    # Request a device code
    device_code, user_code = request_device(idp['request_arguments']['client_id'], idp['request_arguments']['scope'], idp['request_url'])
    if device_code == None:
        return False

    # Update poll arguments
    idp['poll_arguments']['device_code'] = device_code

    # Print User information
    if qr_code:
        print('\033[1;32m[+]\033[0m Please scan the following QR code with your mobile device:')
        generate_qr_code(idp['user_url'])
    else:
        print('\033[1;32m[+]\033[0m Please visit the following URL in your browser: ' + idp['user_url'])
    
    print('\033[1;32m[+]\033[0m Enter the following code when prompted: ' + user_code)

    # Poll for a user token
    response_dict = poll_for_token(idp['poll_url'], idp['poll_arguments'])
    if response_dict == None:
        return False

    return True

def login(username):
    DATABASE_PATH = "/tmp/project_1.sqlite"
    # DATABASE_PATH = "/etc/project_1.sqlite"

    # Get available idps for current user
    idps = get_idps(username, DATABASE_PATH)

    # Present user with a list of available idps and prompt for selection
    print('\033[1;33m[?]\033[0m Please select an IdP:')
    for i, idp in enumerate(idps):
        print('[' + str(i) + '] ' + idp[0])

    selection = raw_input('[>] ')

    # Validate selection
    try:
        selection = int(selection)
        if selection < 0 or selection >= len(idps):
            raise ValueError
    except ValueError:
        print('\033[1;31m[!]\033[0m Invalid selection')
        return False

    # Get idp information
    idp = idps[selection][0]

    # Get idp information
    idp = get_idp(idp, username, DATABASE_PATH)

    # Ask user if they would like to scan a QR code or enter a URL
    qr_code = raw_input('\033[1;33m[?]\033[0m Would you like to scan a QR code? [y/n] ')
    if qr_code == 'y':
        qr_code = True
    elif qr_code == 'n':
        qr_code = False
    else:
        print('\033[1;31m[!]\033[0m Invalid input')
        return False

    # Request a device code
    return oauth2(idp, qr_code)

# test implementation
def main():
    username = 'pengrey'

    if login(username):
        print('\033[1;32m[+]\033[0m Login successful')
    else:
        print('\033[1;31m[!]\033[0m Login failed')

if __name__ == '__main__':
    main()



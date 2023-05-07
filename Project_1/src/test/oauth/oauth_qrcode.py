# Test implementation of multiple device flow authentications
import requests
import time
import qrcode

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
        print('Please scan the following QR code with your mobile device:')
        generate_qr_code(idp['user_url'])
    else:
        print('Please visit the following URL in your browser:' + idp['user_url'])
    
    print('Enter the following code when prompted: ' + user_code)

    # Poll for a user token
    response_dict = poll_for_token(idp['poll_url'], idp['poll_arguments'])
    if response_dict == None:
        return False

    # Print token
    print response_dict['access_token']
    return True

# test implementation
def main():
    # Github
    github = {
        'request_url': 'https://github.com/login/device/code',
        'request_arguments': {
            'client_id': 'a6b22dc869165e33cd5b',
            'scope': 'user:email'
        },
        'user_url': 'https://github.com/login/device',
        'poll_url': 'https://github.com/login/oauth/access_token',
        'poll_arguments': {
            'client_id': 'a6b22dc869165e33cd5b',
            'device_code': None,
            'grant_type': 'urn:ietf:params:oauth:grant-type:device_code'
        }
    }

    # Google
    google = {
        'request_url': 'https://accounts.google.com/o/oauth2/device/code',
        'request_arguments': {
            'client_id': '572020437555-gmf1c30oaqla0749f730udft703es94q.apps.googleusercontent.com',
            'scope': 'https://www.googleapis.com/auth/userinfo.email'
        },
        'user_url': 'https://accounts.google.com/o/oauth2/device/usercode',
        'poll_url': 'https://accounts.google.com/o/oauth2/token',
        'poll_arguments': {
            'client_id': '572020437555-gmf1c30oaqla0749f730udft703es94q.apps.googleusercontent.com',
            'client_secret': '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',
            'device_code': None,
            'grant_type': 'urn:ietf:params:oauth:grant-type:device_code'
        }
    }

    # Test Github
    if oauth2(github, True):
        print('Github authentication successful')
    else:
        print('Github authentication failed')

    # Test Google
    if oauth2(google, True):
        print('Google authentication successful')
    else:
        print('Google authentication failed')
    
if __name__ == '__main__':
    main()
        

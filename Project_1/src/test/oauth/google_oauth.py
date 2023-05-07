# Test implementation of google device flow authentication
import requests
import time

# Request a device code from google
def request_device_code(client_id, scope):
    url = 'https://accounts.google.com/o/oauth2/device/code'
    data = {
        'client_id': client_id,
        'scope': scope
    }

    response = requests.post(url, data=data)

    if response.status_code == 200:
        return response.text
    else:
        return None

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

# Poll google for a user token
def poll_for_token(client_id, client_secret, device_code, interval, timeout):
    url = 'https://accounts.google.com/o/oauth2/token'
    data = {
        'client_id': client_id,
        'client_secret': client_secret,
        'device_code': device_code,
        'grant_type': 'urn:ietf:params:oauth:grant-type:device_code'
    }

    start_time = time.time()
    while True:
        response = requests.post(url, data=data)
        if response.status_code == 200:
            response_dict = parse_response(response.text, {})

            if 'access_token' in response_dict:
                return response_dict
            else:
                time.sleep(interval)
        else:
            return None

# test implementation
def main():
    client_id = '572020437555-gmf1c30oaqla0749f730udft703es94q.apps.googleusercontent.com'
    client_secret = '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'
    scope = 'https://www.googleapis.com/auth/userinfo.email'
    interval = 5
    timeout = 60

    device_code_text = request_device_code(client_id, scope)
    if device_code_text == None:
        print 'Failed to request device code'
        return None

    # Parse device code response
    device_code = parse_response(device_code_text, {})

    print('Device code: ' + device_code['device_code'])
    print('User code: ' + device_code['user_code'])
    print('Verification URL: ' + device_code['verification_url'])

    # Ask user to continue
    #raw_input('Press enter to continue...')

    # Poll for user token
    user_token = poll_for_token(client_id, client_secret, device_code['device_code'], interval, timeout)
    if user_token == None:
        print 'Failed to get user token'
        return None

    print('User token: ' + user_token['access_token'])

if __name__ == '__main__':
    main()

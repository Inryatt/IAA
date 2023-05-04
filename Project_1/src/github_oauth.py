# Test implementation of github device flow authentication
import requests
import json
import time

# Request a device code from github
def request_device_code(client_id, scope):
    url = 'https://github.com/Pengrey/IAA'
    data = {
        'client_id': client_id,
        'scope': scope
    }

    response = requests.post('https://github.com/login/device/code', data=data)

    if response.status_code == 200:
        print response.text
        return response.text
    else:
        return None

# Parse API response
def parse_response(response, dict):
    for pair in response.split('&'):
        key, value = pair.split('=')
        dict[key] = value
    
    return dict

# Poll github for a user token
def poll_for_token(client_id, device_code, interval, timeout):
    url = 'https://github.com/Pengrey/IAA'
    data = {
        'client_id': client_id,
        'device_code': device_code,
        'grant_type': 'urn:ietf:params:oauth:grant-type:device_code'
    }

    start_time = time.time()
    while True:
        response = requests.post('https://github.com/login/oauth/access_token', data=data)
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
    client_id = 'a6b22dc869165e33cd5b'
    scope = 'user:email'
    interval = 5
    timeout = 60

    device_code_text = request_device_code(client_id, scope)
    if device_code_text == None:
        print('Failed to get device code')
        return

    # Parse device code response
    device_code = parse_response(device_code_text, {})
    
    print('Device code: ' + device_code['device_code'])
    print('User code: ' + device_code['user_code'])
    print('Verification URL: https://github.com/login/device')

    token = poll_for_token(client_id, device_code['device_code'], interval, timeout)
    if token == None:
        print('Failed to get token')
        return

    print('Token: ' + token['access_token'])

if __name__ == '__main__':
    main()

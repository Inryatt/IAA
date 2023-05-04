# IAA Project 1

## Project Description
This project is a PAM module that allows users to authenticate using external Identity Providers (IdPs). The module uses OAuth 2.0 with the Device Authorization Grant strategy to authenticate users.

## Project Structure

## Project Setup
### Pam Python Installation
To install the Pam Python module, run the following command:
```bash
sudo apt-get install libpam0g libpam-runtime libpam0g-dev python2 python2-dev
git clone https://github.com/Ralnoc/pam-python.git
cd pam-python
vim ./src/setup.py # Change the python version to 2.7
vim ./src/test.py # Change the python version to 2.7
sudo make && make install
ls /lib/security # Check if pam_python.so is present
```

### Project Installation
To install the project, run the following command:
```bash
```

Edit /etc/pam.d/common-auth and add this line to the top of the file:
```bash
auth sufficient pam_python.so /path/to/Project_1/src/login.py
```

Edit /etc/pam.d/common-session and add this line to the top of the file:
```bash
session	sufficient pam_python.so /path/to/Project_1/src/login.py
```

### Project Configuration
To configure the project, run the following command:
```bash
```

## Project Usage
To use the project, run the following command:
```bash
```

## Testing
You can test the module by switching users:
```bash
su <testuser>
```
# IAA Project 1

## Project Description
This project is a PAM module that allows users to authenticate using external Identity Providers (IdPs). The module uses OAuth 2.0 with the Device Authorization Grant strategy to authenticate users.

## Project Structure

## idp_login
`idp_login` is a command-line application to manage Identity Providers (IdPs) and identity attributes for users in a protected repository.

### Usage
```bash
idp_login <action> [options]
```
- `<action>`: The action to perform, such as setting, changing, deleting, or listing IdPs and identity attributes.
- `[options]`: Optional arguments to specify additional information for the action.

### Actions and Options

1. **Set, change, or delete IdPs and their operational parameters (for host administrators):**

```bash
idp_login manage-idp [--operation set|change|delete] [--idp IDP_NAME] [--params PARAMS]
```

- `--operation`: The operation to perform, e.g. set, change, or delete.
- `--idp IDP_NAME`: The name of the IdP to be managed (required for set and change operations).
- `--params PARAMS`: The operational parameters for the IdP (required for set and change operations).

2. **Set, change, or delete identity attributes for a given IdP for the current user:**

```bash
idp_login manage-attributes [--operation set|change|delete] [--idp IDP_NAME] [--attributes ATTRIBUTES]
```

- `--operation`: The operation to perform, e.g. set, change, or delete.
- `--idp IDP_NAME`: The name of the IdP whose attributes need to be managed (required for set and change operations).
- `--attributes ATTRIBUTES`: The identity attributes for the IdP (required for set and change operations).

3. **List all users with registered IdPs (for host administrators):**

```bash
idp_login list-users
```

4. **List the IdPs registered for the current user and the identity parameters for each IdP:**

```bash
idp_login list-idps
```

## Examples

- To set an IdP with its operational parameters:

```bash
sudo idp_login manage-idp --operation set --idp IDP_NAME --params PARAMS
```

- To delete an IdP:

```bash
sudo idp_login manage-idp --operation delete --idp IDP_NAME
```

- To set identity attributes for a given IdP:

```bash
idp_login manage-attributes --operation set --idp IDP_NAME --attributes ATTRIBUTES
```

- To list all users with registered IdPs:

```bash
sudo idp_login list-users
```

- To list the IdPs registered for the current user:

```bash
idp_login list-idps
```

## Project Setup
### Pam Python Installation
To install the Pam Python module, run the following command:
```bash
sudo apt-get install libpam0g libpam-runtime libpam0g-dev python2 python2-dev
git clone https://github.com/Ralnoc/pam-python.git
cd pam-python
vim ./src/setup.py # Change the python version to 2.7
vim ./src/test.py  # Change the python version to 2.7
sudo make && make install
ls /lib/security   # Check if pam_python.so is present
```

### Pam Python Module
To install the project, run the following command:
```bash
git clone https://github.com/Pengrey/IAA.git
cd Project_1
sudo python pip install -r ./src/pam/requirements.txt
```

Edit /etc/pam.d/common-auth and add this line to the top of the file:
```bash
auth sufficient pam_python.so /path/to/Project_1/src/login.py
```

Edit /etc/pam.d/common-session and add this line to the top of the file:
```bash
session	sufficient pam_python.so /path/to/Project_1/src/login.py
```

### IdP Login
To install the IdP Login command-line application, run the following command:
```bash
cd Project_1/src/idp_login
go build -o idp_login main.go       // Build the application
sudo groupadd idpadmins             // Create a group for IdP administrators
sudo usermod -a -G idpadmins root   // Add root to the group
sudo chown root idp_login           // Set the owner to root
sudo chmod 4750 idp_login           // Set the setuid bit
sudo mv idp_login /usr/bin          // Move the application to /usr/bin
```
import qrcode

# Define the URL
url = 'https://example.com'

# Create a QR code instance
qr = qrcode.QRCode(version=1, box_size=1, border=1)

# Add the data to the QR code
qr.add_data(url)

# Compile the QR code
qr.make(fit=True)

# Get the QR code data as a list of lists
data = qr.get_matrix()

# Print the QR code to the terminal
for row in data:
    for cell in row:
        if cell:
            print '##',
        else:
            print '  ',
    print

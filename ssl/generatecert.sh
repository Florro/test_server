echo Generate CA key:
openssl genrsa -passout pass:1111 -des3 -out ./localhost/ca.key 4096

echo Generate CA certificate:
openssl req -passin pass:1111 -new -x509 -days 365 -key ./localhost/ca.key -out ./localhost/ca.crt -subj  "/C=US/ST=CA/L=Cupertino/O=YourCompany/OU=YourApp/CN=MyRootCA"

echo Generate server key:
openssl genrsa -passout pass:1111 -des3 -out ./localhost/server.key 4096

echo Generate server signing request:
openssl req -passin pass:1111 -new -key ./localhost/server.key -out ./localhost/server.csr -subj  "/C=US/ST=CA/L=Cupertino/O=YourCompany/OU=YourApp/CN=localhost"

echo Self-sign server certificate:
openssl x509 -req -passin pass:1111 -days 365 -in ./localhost/server.csr -CA ./localhost/ca.crt -CAkey ./localhost/ca.key -set_serial 01 -out ./localhost/server.crt

echo Remove passphrase from server key:
openssl rsa -passin pass:1111 -in ./localhost/server.key -out ./localhost/server.key

echo Generate client key
openssl genrsa -passout pass:1111 -des3 -out ./localhost/client.key 4096

echo Generate client signing request:
openssl req -passin pass:1111 -new -key ./localhost/client.key -out ./localhost/client.csr -subj  "/C=US/ST=CA/L=Cupertino/O=YourCompany/OU=YourApp/CN=%CLIENT-COMPUTERNAME%"

echo Self-sign client certificate:
openssl x509 -passin pass:1111 -req -days 365 -in ./localhost/client.csr -CA ./localhost/ca.crt -CAkey ./localhost/ca.key -set_serial 01 -out ./localhost/client.crt

echo Remove passphrase from client key:
openssl rsa -passin pass:1111 -in ./localhost/client.key -out ./localhost/client.key

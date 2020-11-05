# Tile38 Webserver
A docker image with a http web server written in Go that routes all http calls to the tile38 server.
The server will check every requests using one of the available methods: jwt (using rsa or hmac) and basic (using a secret)
and only authorized requests will be routed. Since you can only using tile38 AUTH command by a cli interface,
this docker image allows you ensure http calls only come from authorized clients.

# Environment variables
```
# URI for the tile38 server
# !!IMPORTANT!! DO NOT INCLUDE TRAILING SLASHES
TILE38_URI="http://localhost:1234"
# valid types: jwt_rsa, jwt_hmac and basic
#jwt_rsa and jwt_hmac expect a header Authorization in the form 'Authorization: Bearer jwt_token'
#basic expects a header Authorization in the form 'Authorization: secret'
VALIDATION_TYPE="jwt_rsa"
# VALIDATION_SECRET is required if VALIDATION_TYPE is equal to jwt_hmac or basic
#VALIDATION_SECRET="mysecret"
# URI for a RSA Public key pem file used during request validation if VALIDATION_TYPE is jwt_rsa
PUBKEY_URI="https://my.public.key/key.pem"
# Address and port that the server will bind to
SERVER_ADDR="0.0.0.0:4433"
#Server .cert file
SERVER_CERT=server.crt
#Server key file
SERVER_KEY=server.key
# Timeouts are expressed in seconds
SERVER_WRITE_TIMEOUT=15
SERVER_READ_TIMEOUT=15
# Public key cache duration in minutes
CACHE_DURATION=10
```

# Generating server.crt and server.key
```bash
openssl req -x509 -nodes -newkey ec -pkeyopt ec_paramgen_curve:secp384r1 -keyout server.key -out server.crt -days 3650
```

# Auth Api Service
This is a server that manages the users of applications.

## routes 

### Application Routes
```
### POST https://localhost:4444/application/register
Auth-Service-Authorization {a secret key generated from a common secret key}
{
    "name": "App2",
    "password": "password",
    "description": "web app 2"
}
returns {
    "base_secret_key": "38baec0b17b2a98d4d8c0fa237e979ddc5fbf908b4077e8e949ada1c28a06b67",
    "message": "Application created successfully"
}
```

```
### POST https://localhost:4444/application/login
Auth-Service-Authorization {a secret key generated from a common secret key}
{
    "name": "App2",
    "password": "password",
}
returns {
    "base_secret_key": "38baec0b17b2a98d4d8c0fa237e979ddc5fbf908b4077e8e949ada1c28a06b67",
    "message": "Application login successfully"
}
```
### User Proxied Routes Thrue application server

```
### POST https://localhost:4444/user/register
Auth-Service-Authorization {a secret key generated from a common secret key}
Application-Secret {base_secret_key}
{
    "name":"tim1",
    "email":"tim1@gmail.com",
    "password":"tim123"
}
returns {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTYwNTU2NDksInRva2VuX3R5cGUiOiJhY2Nlc3NfdG9rZW4iLCJ1c2VyX2lkIjoiNDJvaUhvYXpFdE5rOVZ5UUVRZXVrelZhcEVUbnpBWjl0TExrIiwidXNlcl90eXBlIjoxfQ.x9DfZ5O9R-OShoYCrZa5vavNeuQ0dQbzd6EJzBi1mwY",
    "message": "User created successfully",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg2NDU4NDksInRva2VuX3R5cGUiOiJyZWZyZXNoX3Rva2VuIiwidXNlcl9pZCI6IjQyb2lIb2F6RXROazlWeVFFUWV1a3pWYXBFVG56QVo5dExMayIsInVzZXJfdHlwZSI6MX0.DQhnSIv_tD8V5kyWX_-uq-VuhemgDXW6dt9ObwbUdjA"
}
```

```
### POST https://localhost:4444/user/login
Auth-Service-Authorization {a secret key generated from a common secret key}
Application-Secret {base_secret_key}
{
    "email":"tim1@gmail.com",
    "password":"tim123"
}
returns {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTYwNTU3NzUsInRva2VuX3R5cGUiOiJhY2Nlc3NfdG9rZW4iLCJ1c2VyX2lkIjoiNDJvaUhvYXpFdE5rOVZ5UUVRZXVrelZhcEVUbnpBWjl0TExrIiwidXNlcl90eXBlIjoxfQ.RjXfGtHudbJT8tvWoyHfUI5uxAqh5WH2lmXIctdoFFo",
    "message": "User login successfully",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg2NDU5NzUsInRva2VuX3R5cGUiOiJyZWZyZXNoX3Rva2VuIiwidXNlcl9pZCI6IjQyb2lIb2F6RXROazlWeVFFUWV1a3pWYXBFVG56QVo5dExMayIsInVzZXJfdHlwZSI6MX0.E1_7E-9WX4Bw1umZLx9ewrRvB7u0kiZsHMR9Mb5O7lg"
}
```

```
### POST https://localhost:4444/user/logout
Auth-Service-Authorization {a secret key generated from a common secret key}
Application-Secret {base_secret_key}
Authorization {access_token (dont add "bearer" in front) }
returns {
    "message": "User logout successfully"
}
```

```
### GET https://localhost:4444/user/token
Auth-Service-Authorization {a secret key generated from a common secret key}
Application-Secret {base_secret_key}
Authorization {access_token (dont add "bearer" in front) }
returns {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTYwNTYxOTksInRva2VuX3R5cGUiOiJhY2Nlc3NfdG9rZW4iLCJ1c2VyX2lkIjoidzlka21KTDI5d0ZIOEdKd0ZXQ1hrU2pGUlUzRG54eGRNc004IiwidXNlcl90eXBlIjoxfQ.SdWS3liumzSQHzTwfE_yR2KRrs9klLar9jHqpDlZS-c",
    "message": "New access token generated successfully"
}
```

# Build
```Dockerfile
sudo docker compose up
```

# Localy
Install dependencies
```bash
go mod tidy
```

run with air and .env
```bash
export $(egrep -v ‘^#’ .env | xargs) && air
```

# Intresting Topics
Rate limiter more specific examples : https://artursiarohau.medium.com/go-chi-rate-limiter-useful-examples-8277dc4d4ff5

# Database setup
???
```sql
create database authservicedb;
use authservicedb;

CREATE USER 'username'@'localhost' IDENTIFIED BY 'password';

-- Grant privileges to the new user for a specific database
GRANT ALL PRIVILEGES ON authservicedb.* TO 'username'@'localhost';


-- Flush privileges to apply changes
FLUSH PRIVILEGES;
```
# SSL
on wsl
```bash
openssl req -nodes -new -x509 -keyout server.key -out server.cert
```
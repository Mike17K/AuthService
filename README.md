# Auth Api Service
This is a server that manages the users of applications.


# Build

## Docker Network 
```bash 
sudo docker network create backend-network
```

## Create Mysql docker volumes
```bash
docker volume create custom-mysql-volume
```

## Run Mysql Container
```bash
sudo docker run -d \
  --network backend-network \
  --name custom-mysql \
  -e MYSQL_ROOT_PASSWORD=root_password \
  -e MYSQL_DATABASE=authservicedb \
  -e MYSQL_USER=username \
  -e MYSQL_PASSWORD=password \
  -p 3333:3306 \
  -v mysql_data_volume:/var/lib/mysql \
  -v $(pwd)/AuthService/database/init.sql:/docker-entrypoint-initdb.d/init.sql \
  --restart always \
  mysql:8.3.0
```

### Setup Enviroment Variables
In the .env file

```bash
cp example.env .env && sudo nano .env
```

```
SERVICE_SECRET_KEY=<random numbers and chars>
PORT=<the port you want the auth service to run>

# Database
DB_USER=root
DB_PASSWORD=root_password
DB_HOST=custom-mysql
DB_PORT=3306
DB_NAME=authservicedb
```

## Docker image
```bash
docker build -t auth-service .
```

```bash
docker run -d \
  --network backend-network \
  --env-file .env \
  --name custom-auth-service \
  -p 4444:4444 \
  --restart always \
  auth-service
```


## Docker compose
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

# Trubleshooting
copying the ssh key to get it to my secrets in the repo for the CICD
```bash
clip < ~/.ssh/....pub
```
Steps to Resolve Docker Permission Issues When CICD is deploying to server over ssh
Add User to Docker Group on server:
```
sudo usermod -aG docker <username>
```
and the logout and login again

Run a database mysql for the service to connect
```bash
docker run -d \
  --name mysql_service_container \
  -e MYSQL_ROOT_PASSWORD=root_password \
  -e MYSQL_DATABASE=authservicedb \
  -e MYSQL_USER=username \
  -e MYSQL_PASSWORD=password \
  -p 3309:3306 \
  -v $(pwd)/AuthService/database/mysql_data:/var/lib/mysql \
  -v $(pwd)/AuthService/database/init.sql:/docker-entrypoint-initdb.d/init.sql \
  --restart always \
  mysql:8.3.0
```

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
# Intresting Topics
Rate limiter more specific examples : https://artursiarohau.medium.com/go-chi-rate-limiter-useful-examples-8277dc4d4ff5

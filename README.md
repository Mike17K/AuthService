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
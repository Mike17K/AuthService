-- Grant permission for the user to connect from the Docker network
GRANT ALL PRIVILEGES ON authservicedb.* TO 'username'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON randevousservicedb.* TO 'username'@'%' IDENTIFIED BY 'password';

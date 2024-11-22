## Running Database 
RUN=docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_USER=myuser -e MYSQL_PASSWORD=mypassword -e MYSQL_DATABASE=mydatabase -p 3306:3306 -d mysql:latest 
docker exec -it my-mysql bash
mysql -u myuser -p
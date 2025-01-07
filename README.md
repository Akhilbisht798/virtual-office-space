## Running Database 
RUN=docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_USER=myuser -e MYSQL_PASSWORD=mypassword -e MYSQL_DATABASE=mydatabase -p 3306:3306 -d mysql:latest 
docker exec -it my-mysql bash
mysql -u myuser -p


### Running MinIo
docker run -d --name minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e "MINIO_ROOT_USER=ROOT" \
  -e "MINIO_ROOT_PASSWORD=password" \
  quay.io/minio/minio server /data --console-address ":9001"



- Start both minio server or s3 and a database put the map in minio so it can import that.
- should have entry for the map and the room and should have it on s3 in zip form (later automate it).
- Maps stored :- maps/map1.zip maps/map2.zip etc.
- Make a database entry for that map.

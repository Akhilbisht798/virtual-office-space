## Development

### Running Database 
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

- When creating a user also create a sprite entry for it (haven't put it on use in client though). 
- create a sprite add it to s3 with entry on db.

### Maps
- Create a map using tiled. 
- make sure to embed you tileset on your map.
- export it as a json and make sure path are right. 
- make a zip and upload it to s3 and put database entry for it. 
- add thumbnail directory in map directory and add map.png in it.

### Spirites
- have a spirites in spirites/charecterId
- can create using chatgpt.

### TODO
- [ ] detailing
  - [ ] maps have thumbnail icons and spaces shows it.
  - [ ] have one more good map and improve space creation page.
  - [ ] same charecter for same user. 
  - [ ] good ui like for shadcn etc. 
  - [ ] show recent spaces. 
- [ ] hosting it. 

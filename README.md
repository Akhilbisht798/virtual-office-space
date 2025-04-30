# 🌐 Nuzo

A **real-time interactive virtual space** where users can:  
- 🚶‍♂️ **Move avatars** in a shared 2D world.  💬 **Chat with others** via video calls.  
- 🎨 **Customize rooms** add space and maps which they like. 

**Live Demo:** [virtual-office-space.vercel.app](https://virtual-office-space.vercel.app/)  

**Tech Stack:**  
- **Frontend:** React.js + Phaser.js (for 2D rendering)  
- **Backend:** Go + WebSockets + docker 
- **Database:** mysql (aiven)  
- **Deployment:** Vercel (Frontend) + render (Backend) + tigirs(s3 fly.io)

## Development

**For local development** use **Docker for database and object storage**
**Steps for spining up docker database and object storage**

**Database**
- `docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_USER=myuser -e MYSQL_PASSWORD=mypassword -e       MYSQL_DATABASE=mydatabase -p 3306:3306 -d mysql:latest`
- `docker exec -it my-mysql bash`
- `mysql -u myuser -p`

**Object Storage**
```bash
docker run -d --name minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e "MINIO_ROOT_USER=ROOT" \
  -e "MINIO_ROOT_PASSWORD=password" \
  quay.io/minio/minio server /data --console-address ":9001"
```

**After starting up the Docker container**
- Start both minio server or s3 and a database put the map in minio so it can import that.
- should have entry for the map and the room and should have it on s3 in zip form (later automate it).
- Maps stored :- maps/map1.zip maps/map2.zip etc.
- Make a database entry for that map.

- When creating a user also create a sprite entry for it (haven't put it on use in client though). 
- create a sprite add it to s3 with entry on db.

**How to add new Map?**
- Create a map using tiled. 
- make sure to embed you tileset on your map.
- export it as a json and make sure path are right. 
- make a zip and upload it to s3 and put database entry for it. 
- add thumbnail directory in map directory and add map.png in it.

**How to add new Spirites?**
- have a spirites in spirites/charecterId
- can create using chatgpt.

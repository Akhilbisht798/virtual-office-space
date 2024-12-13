import Phaser from "phaser";
import { getSocketInstance } from "../../hooks/useSocket";

class GameScene extends Phaser.Scene {
  static SCENE_KEY = "GameScene";
  constructor() {
    super({ key: GameScene.SCENE_KEY });
    this.player;
    this.cursors;
    this.players = {};
    this.ws;
    this.prevX = 80;
    this.prevY = 30;
    this.roomId = "90";
    this.userId;
  }

  init(data) {
    this.files = data.files;
    this.mapName;
    for (let key in this.files) {
      const keyArr = key.split("/")
      this.mapName = keyArr[0]
      break
    }
  }

  preload() {
    const jsonKeys = Object.keys(this.files).filter(key => key.endsWith('.json'));
    const mapJson = this.files[jsonKeys[0]];
    const map = JSON.parse(mapJson)
    
    this.load.tilemapTiledJSON("map", map);

    map.tilesets.forEach(tileset => {
      const assetFileURl = this.mapName + tileset.image;
      this.load.image(tileset.name, this.files[assetFileURl])
    });

    this.load.spritesheet("run", "assets/spirite/Run.png", {
      frameWidth: 42,
      frameHeight: 42,
    });
    this.load.spritesheet("idle", "assets/spirite/Idle.png", {
      frameWidth: 42,
      frameHeight: 42,
    });
  }

  create() {
    const map = this.make.tilemap({ key: "map" });
    this.map = map;
    const tilesets = this.map.tilesets.map(tileset => {
      return this.map.addTilesetImage(tileset.name, tileset.name)
    })

    let collisionLayer;

    map.layers.map((layer) => {
      try {
        const createLayer = map.createLayer(
          layer.name,
          tilesets,
          0,
          0,
        );
        if (layer.name === "boundary") {
          collisionLayer = createLayer;
          console.log(collisionLayer)
          createLayer.setCollisionByExclusion([-1]);
        }
      } catch (error) {
        console.log("error creating the layer: ", error);
      }
    });

    //Player
    const idle = {
      key: "idle",
      frames: this.anims.generateFrameNumbers("idle", {
        frames: [0, 1, 2, 3],
      }),
      frameRate: 8,
      repeat: -1,
    };
    const running = {
      key: "running",
      frames: this.anims.generateFrameNumbers("run", {
        frames: [0, 1, 2, 3, 4, 5],
      }),
      frameRate: 8,
      repeat: -1,
    };
    this.player = this.physics.add.sprite(this.prevX, this.prevY, "run");
    this.player.anims.create(idle);
    this.player.anims.create(running);

    // Collision
    if (collisionLayer) {
      this.physics.add.collider(this.player, collisionLayer);
    }
    console.log("Create function done.");

    this.cameras.main.startFollow(this.player);
    this.cameras.main.setZoom(2.4);
    this.cameras.main.setBounds(
      0,
      0,
      this.map.widthInPixels,
      this.map.heightInPixels,
    );

    this.ws = getSocketInstance();
    this.ws.onopen = () => {
      console.log("connected to socket server");

      const token = localStorage.getItem("jwt")

      this.ws.send(
        JSON.stringify({
          type: "join",
          payload: {
            roomId: this.roomId,
            x: this.prevX,
            y: this.prevY,
            jwt: token,
          },
        }),
      );
    };
    this.ws.onmessage = (message) => {
      const data = JSON.parse(message.data);
      const type = data["type"];
      const payload = data["payload"];
      switch (type) {
        case "space-joined":
          console.log("Space Joined working");
          this.spaceJoined(payload);
          break;
        case "user-join":
          this.userJoin(payload);
          break;
        case "user-left":
          this.userLeft(payload);
          break;
        case "movement":
          this.movement(payload);
          break;
        case "movement-rejected":
          console.log(payload);
          break;
      }
    };
    this.ws.onclose = () => {
      console.log("Socket disconnect");
    };
    this.cursor = this.input.keyboard.createCursorKeys();
  }

  update() {
    const speed = 160;

    if (this.prevX !== this.player.x || this.prevY !== this.player.y) {
      this.prevX = this.player.x;
      this.prevY = this.player.y;
      this.ws.send(
        JSON.stringify({
          type: "move",
          payload: {
            roomId: this.roomId,
            userId: this.userId,
            x: this.player.x,
            y: this.player.y,
          },
        }),
      );
    }

    this.player.setVelocity(0);
    if (
      this.cursor.left.isDown ||
      this.cursor.right.isDown ||
      this.cursor.up.isDown ||
      this.cursor.down.isDown
    ) {
      this.player.play("running", true);
    } else {
      this.player.play("idle", true);
    }

    if (this.cursor.left.isDown) {
      this.player.setVelocityX(-speed);
    } else if (this.cursor.right.isDown) {
      this.player.setVelocityX(speed);
    }

    if (this.cursor.up.isDown) {
      this.player.setVelocityY(-speed);
    } else if (this.cursor.down.isDown) {
      this.player.setVelocityY(speed);
    }
  }

  spaceJoined(payload) {
    const spawn = payload["spawn"];
    const users = payload["users"];
    this.userId = payload["userId"]
    const sprite = payload["sprite"]
    console.log("sprite ", sprite)
    console.log(this.userId)
    console.log("users: ", users)

    users.forEach((u) => {
      const p = this.physics.add.sprite(u.x, u.y, "run");
      p.play("idle", true);
      this.players[u.userId] = p;
    });
    this.player.x = spawn.x;
    this.player.y = spawn.y;
  }

  userJoin(payload) {
    const p = this.physics.add.sprite(payload.x, payload.y, "run");
    this.players[payload.userId] = p;
  }

  userLeft(payload) {
    const uId = payload["userId"];
    this.players[uId].destroy();
    delete this.players[uId];
  }

  movement(payload) {
    this.players[payload.userId].setPosition(payload.x, payload.y);
  }

  movementRejected(payload) {
    this.player.setPosition(payload.x, payload.y);
  }
}

export default GameScene;

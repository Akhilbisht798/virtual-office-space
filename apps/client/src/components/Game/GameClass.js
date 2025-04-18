import Phaser from "phaser";
import { getSocketInstance } from "../../hooks/useSocket";
import AgoraRTC from "agora-rtc-sdk-ng";
import {v4 as uuidv4 } from "uuid";

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
    this.roomId;
    this.userId;
    this.remoteVideo = {};
    this.onCall = false;

    this.localTracks = {
      videoTrack: null,
      audioTrack: null,
    };
    this.agoraClient = null;
    this.agoraAppID = import.meta.env.VITE_AGORA_APP_ID;
    this.agoraToken = null;
    this.agoraChannel = null;
  }

  init(data) {
    this.files = data.files;
    this.roomId = data.spaceId
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
    this.load.spritesheet("ch1", "/spirite/charecter1.png", {
      frameWidth: 142/3,
      frameHeight: 424/4,
    })
    this.load.spritesheet("ch2", "/spirite/charecter2.png", {
      frameWidth: 148/3,
      frameHeight: 371/4,
    })
    this.load.spritesheet("ch3", "/spirite/charecter3.png", {
      frameWidth: 143/3,
      frameHeight: 380/4,
    })
    this.load.spritesheet("ch4", "/spirite/charecter4.png", {
      frameWidth: 144/3,
      frameHeight: 360/4,
    })
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

    const randomId = Math.floor(Math.random() * 4) + 1;
    const charecterId = "ch" + randomId
    this.player = this.createSpirite(charecterId)

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

      console.log("joining to roomId", this.roomId)
      this.ws.send(
        JSON.stringify({
          type: "join",
          payload: {
            roomId: this.roomId,
            x: this.prevX,
            y: this.prevY,
            jwt: token,
            spirite: charecterId,
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
        case "call-req":
          this.onCallRequest(payload)
          break;
      }
    };
    this.ws.onclose = () => {
      console.log("Socket disconnect");
    };
    this.cursor = this.input.keyboard.createCursorKeys();

    //init agora here.
    this.initAgora()
  }

  update() {
    const speed = 160;
    let isNear = false;

    if (this.onCall) {
      this.showLeaveButton()
      console.log("LEAVE BUTTON ADDED")
    } else {
      const btn = document.getElementById("video-stop-btn")
      if (btn) {
        btn.remove()
      }
    }

    Object.values(this.players).forEach(player => {
      const distance = Phaser.Math.Distance.Between(this.player.x, this.player.y, player.x, player.y);
      if (distance <= 40) {
        if (!this.onCall) {
          this.showVideoCallBtn(player.userId)
        }
        isNear = true;
      } 
    })

    if (!isNear) {
      const btn = document.getElementById("video-call-btn")
      if (btn) {
        btn.remove()
      }
    }

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

    if (this.cursor.left.isDown) {
      this.player.setVelocityX(-speed);
      this.player.play("chLeft", true);
    } else if (this.cursor.right.isDown) {
      this.player.setVelocityX(speed);
      this.player.play("chRight", true);
    } else if (this.cursor.up.isDown) {
      this.player.setVelocityY(-speed);
      this.player.play("chUp", true);
    } else if (this.cursor.down.isDown) {
      this.player.setVelocityY(speed);
      this.player.play("chDown", true);
    } else {
      this.player.play("chIdle", true)
    }
  }

  async initAgora() {
    this.agoraClient = AgoraRTC.createClient({ mode: "rtc", codec: "vp8" });
    this.agoraClient.on("user-published", async (user, mediaType) => {
      try {
        await this.agoraClient.subscribe(user, mediaType)
        if (mediaType === "video") {
          console.log("remote users joined: ", user.uid)
          await this.displayRemoteVideo(user)
        }

        if (mediaType === 'audio') {
          user.audioTrack.play()
        }
        const remoteUsers = this.agoraClient.remoteUsers
      } catch (error) {
        console.error("Error subscribing to remote user: ", error)
      }
    })

    this.agoraClient.on("user-unpublished", async (user) => {
      const video = document.getElementById(user.uid)
      video && video.remove()
      delete this.remoteVideo[user.uid]
    })

    this.agoraClient.on("user-left", async (user) => {
      const remoteUsers = this.agoraClient.remoteUsers
      if (remoteUsers.length === 0) {
        this.leaveChannel()
      }
    })
  }

  async displayLocalVideo() {
    console.log("Display local track method has been called.")
    const video = document.createElement("video");
    video.id = "local";
    const stream = new MediaStream()
    stream.addTrack(this.localTracks.videoTrack.getMediaStreamTrack());
    video.srcObject = stream;
    video.muted = true;
    video.autoplay = true;
    video.width = 160;
    video.height = 120;
    const videoContainer = document.getElementById("video-container");
    if (videoContainer === null) {
      const videoContainer = document.createElement("div");
      videoContainer.id = "video-container";
      videoContainer.style.position = "absolute";
      videoContainer.style.right = "0";
      videoContainer.style.top = "50%";
      videoContainer.style.display = "flex";
      videoContainer.style.flexDirection = "column";
      videoContainer.style.alignItems = "center";
      videoContainer.style.background = "rgba(0, 0, 0, 0.5)";
      videoContainer.style.padding = "10px";

      document.body.appendChild(videoContainer);
    }

    const container = document.getElementById("video-container");
    container.appendChild(video);
  }

  async displayRemoteVideo(user) {
    const video = document.createElement("video");
    video.id = user.uid;
    video.autoplay = true;
    video.width = 160;
    video.height = 120;
    video.muted = true;

    const videoContainer = document.getElementById("video-container");
    if (videoContainer === null) {
      const videoContainer = document.createElement("div");
      videoContainer.id = "video-container";
      videoContainer.style.position = "absolute";
      videoContainer.style.right = "0";
      videoContainer.style.top = "50%";
      videoContainer.style.display = "flex";
      videoContainer.style.flexDirection = "column";
      videoContainer.style.alignItems = "center";
      videoContainer.style.background = "rgba(0, 0, 0, 0.5)";
      videoContainer.style.padding = "10px";

      document.body.appendChild(videoContainer);
    }

    const container = document.getElementById("video-container");
    container.appendChild(video);
    user.videoTrack.play(video)
    this.remoteVideo[user.uid] = video;
  }

  async showVideoCallBtn(userId) {
    const checkBtn = document.getElementById("video-call-btn");
    if (checkBtn) {
      return;
    }
    const btn = document.createElement("button");
    btn.id = "video-call-btn"
    btn.innerText = "Call"
    btn.style.position = "fixed"
    btn.style.bottom = "20px"
    btn.style.left = "50%"
    btn.style.fontSize = "16px";
    btn.style.padding = "10px 20px";
    btn.style.cursor = "pointer";

    btn.style.backgroundColor = "#007bff"; 
    btn.style.color = "#fff";
    btn.style.border = "none";
    btn.style.borderRadius = "5px";

    btn.addEventListener("click", () => {
      console.log("Lets video call")
      this.makeCall(userId)
    });

    document.body.appendChild(btn)
  }

  async showLeaveButton() {
    const checkBtn = document.getElementById("video-stop-btn");
    if (checkBtn) {
      return;
    }
    const btn = document.createElement("button");
    btn.id = "video-stop-btn"
    btn.innerText = "Leave"
    btn.style.position = "fixed"
    btn.style.bottom = "20px"
    btn.style.left = "50%"
    btn.style.fontSize = "16px";
    btn.style.padding = "10px 20px";
    btn.style.cursor = "pointer";

    btn.style.backgroundColor = "#007bff"; 
    btn.style.color = "#fff";
    btn.style.border = "none";
    btn.style.borderRadius = "5px";

    btn.addEventListener("click", () => {
      this.leaveChannel();
    });

    document.body.appendChild(btn)
  }

  async makeCall(remoteUserId) {
    console.log("Make a call")
    const callID = uuidv4()

    this.ws.send(
      JSON.stringify({
        type: "make-call",
        payload: {
          roomId: this.roomId,
          userId: this.userId,
          remoteUserId: remoteUserId,
          callId: callID, 
        }
      })
    );
  }

  async createLocalTrack() {
    this.localTracks.audioTrack = await AgoraRTC.createMicrophoneAudioTrack();
    this.localTracks.videoTrack = await AgoraRTC.createCameraVideoTrack();
    console.log("Local tracks recived")
  }

  async joinChannel() {
    try {
      await this.agoraClient.join(this.agoraAppID, this.agoraChannel, null, 0);
      await this.createLocalTrack();
      await this.displayLocalVideo();
      await this.agoraClient.publish([this.localTracks.videoTrack, this.localTracks.audioTrack]);
      this.onCall = true;
    } catch (error) {
      console.error("Error joining the channel: ", error)
    }
  }

  spaceJoined(payload) {
    const spawn = payload["spawn"];
    const users = payload["users"];
    this.userId = payload["userId"]
    //const sprite = payload["spirite"]

    users.forEach((u) => {
      const spirite = u.sprite;
      const username = u.username;
      const p = this.createSpirite(spirite);
      const pName = this.add.text(p.x - 10, p.y - 40, username);
      p.play("chIdle", true)
      p.userId = u.userId;
      p.username = pName
      this.players[u.userId] = p;
    });
    this.player.x = spawn.x;
    this.player.y = spawn.y;
  }

  userJoin(payload) {
    const spirite = payload["spirite"]
    const username = payload["username"];
    const x = payload["x"];
    const y = payload["y"];

    const p = this.createSpirite(spirite)
    const pName = this.add.text(x - 10, y - 40, username);
    p.play("chIdle", true)
    p.userId = payload.userId;
    p.username = pName;
    this.players[payload.userId] = p;
  }

  userLeft(payload) {
    const uId = payload["userId"];
    this.players[uId].username.destroy();
    this.players[uId].destroy();
    delete this.players[uId];
  }

  movement(payload) {
    const prevX = this.players[payload.userId].x
    const prevY = this.players[payload.userId].y
    const x = payload.x;
    const y = payload.y;
    if (x > prevX) {
      this.players[payload.userId].anims.play({
        key: "chRight",
        repeat: 0
      })
    } else if (x < prevX) {
      this.players[payload.userId].anims.play({
        key: "chLeft",
        repeat: 0
      })
    } else if (y > prevY) {
      this.players[payload.userId].anims.play({
        key: "chDown",
        repeat: 0
      })
    } else if (y < prevY) {
      this.players[payload.userId].anims.play({
        key: "chUp",
        repeat: 0
      })
    }
    this.players[payload.userId].setPosition(payload.x, payload.y);
    this.players[payload.userId].username.setPosition(payload.x - 10, payload.y - 40);
  }

  movementRejected(payload) {
    this.player.setPosition(payload.x, payload.y);
  }

  onCallRequest(payload) {
    this.agoraChannel = payload.channel;
    this.joinChannel()
    this.ws.send(
      JSON.stringify({
        type: "call-accept",
        payload: {
          userId: this.userId,
          channelId: this.agoraChannel,
        }
      })
    )
    console.log("Channel Joinned")
  }

  leaveChannel() {
    this.agoraClient.leave(() => {
      console.log("Client leaves channel");
    }, (err) => {
      console.log("Client leave failed", err);
    });
    if (this.localTracks.videoTrack) this.localTracks.videoTrack.close();
    if (this.localTracks.audioTrack) this.localTracks.audioTrack.close();
    const videoContainer = document.getElementById("video-container");
    videoContainer && videoContainer.remove();

    this.ws.send(
      JSON.stringify({
        type: "leave-call",
        payload: {
          userId: this.userId,
          channelId: this.agoraChannel,
        }
      })
    )
    this.agoraChannel = null;
    const btn = document.getElementById("video-call-btn")
    if (btn) btn.remove()
    this.onCall = false;
  }

  createSpirite(charecterId) {
    const chDown = {
      key: "chDown", 
      frames: this.anims.generateFrameNumbers(charecterId, {
        frames: [0, 1, 2],
      }),
      frameRate: 8,
      repeat: -1,
    }
    const chLeft = {
      key: "chLeft", 
      frames: this.anims.generateFrameNumbers(charecterId, {
        frames: [3, 4, 5],
      }),
      frameRate: 8,
      repeat: -1,
    }
    const chRight = {
      key: "chRight", 
      frames: this.anims.generateFrameNumbers(charecterId, {
        frames: [6, 7, 8],
      }),
      frameRate: 8,
      repeat: -1,
    }
    const chUp = {
      key: "chUp", 
      frames: this.anims.generateFrameNumbers(charecterId, {
        frames: [9, 10, 11],
      }),
      frameRate: 8,
      repeat: -1,
    }
    const idle = {
      key: "chIdle",
      frames: this.anims.generateFrameNumbers(charecterId, {
        frames: [1],
      }),
      frameRate: 8,
      repeat: -1,
    }
    const player = this.physics.add.sprite(this.prevX, this.prevY, "chRight");
    player.anims.create(chDown);
    player.anims.create(chLeft);
    player.anims.create(chRight);
    player.anims.create(chUp);
    player.anims.create(idle);

    player.setScale(0.5, 0.5);
    player.body.setSize(40, 80, true); // Adjust hitbox
    player.setOffset(5, 7);

    return player;
  }
}

export default GameScene;

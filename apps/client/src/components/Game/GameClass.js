import Phaser from "phaser";

// Maybe player are already in the room. 
// Maybe they join later. 
// so we need to have this class to manage the players
// So I need to display spirte at a certain x,y coordinates.
// this can change on socket update.
class GameScene extends Phaser.Scene {
    constructor() {
        super({ key: "GameScene" });
        this.player;
        this.cursors;
        this.players = [];
    }
    preload() {
        this.load.tilemapTiledJSON("map", "assets/map.json");
        this.load.image("flower", "assets/Pixel 16 Village/Fences/flower fence.png");
        this.load.image("wood", "assets/Pixel 16 Village/Fences/wooden fence.png");
        this.load.image(
            "floor",
            "assets/Pixel 16 Village/Walls & Floor Tiles [Update 1.1]/floor-tiles.png",
        );
        this.load.image(
            "walls",
            "assets/Pixel 16 Village/Walls & Floor Tiles [Update 1.1]/walls.png",
        );

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
        const tileset = {
            walls: map.addTilesetImage("walls", "walls"),
            floor: map.addTilesetImage("floor-tiles", "floor"),
            wood: map.addTilesetImage("wooden fence", "wood"),
            flower: map.addTilesetImage("flower fence", "flower"),
        };

        let collisionLayer;
        console.log("map:", map)

        map.layers.map((layer) => {
            try {
                const createLayer = map.createLayer(
                    layer.name,
                    [tileset.floor, tileset.walls, tileset.wood, tileset.flower],
                    0,
                    0,
                );
                if (layer.name === "boundary") {
                    console.log("collision layer");
                    collisionLayer = createLayer;
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
        }
        const running = {
            key: "running",
            frames: this.anims.generateFrameNumbers("run", {
                frames: [0, 1, 2, 3, 4, 5],
            }),
            frameRate: 8,
            repeat: -1,
        }
        this.player = this.physics.add.sprite(30, 30, "run");
        this.player.anims.create(idle);
        this.player.anims.create(running);


        // Collision
        if (collisionLayer) {
            this.physics.add.collider(this.player, collisionLayer);
        }

        this.cameras.main.startFollow(this.player);
        this.cameras.main.setZoom(2.4);
        this.cameras.main.setBounds(0, 0, this.map.widthInPixels, this.map.heightInPixels);
        this.cursor = this.input.keyboard.createCursorKeys();
    }

    update() {
        const speed = 160;

        this.player.setVelocity(0);
        if (this.cursor.left.isDown || this.cursor.right.isDown ||
            this.cursor.up.isDown || this.cursor.down.isDown) {
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

    addPlayer() {
        this.players.push(this.physics.add.sprite(30, 30, "run"));
    }
}

export default GameScene;
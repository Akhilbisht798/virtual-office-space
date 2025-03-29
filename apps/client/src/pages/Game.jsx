import Phaser from "phaser";
import GameScene from "../components/Game/GameClass";
import { useEffect } from "react";
import { SERVER } from "../global";
import axios from "axios";
import JSZip from "jszip";
import { getSocketInstance } from "../hooks/useSocket";

const sizes = {
  width: 400,
  height: 300,
};

const Game = ({ spaceId }) => {
  async function downloadZip(url) {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error("Failed to download the file");
    }

    return await response.blob();
  }

  async function unzipFile(blob) {
    const zip = new JSZip();
    const unzipped = await zip.loadAsync(blob);

    const files = {};
    for (const [filename, file] of Object.entries(unzipped.files)) {
      //console.log("FileName: ", filename);
      if (!file.dir) {
        const content = await file.async("blob");
        files[filename] = content;
      }
    }
    return files;
  }

  async function processFiles(files) {
    const assetsFile = {};
    for (const [filename, fileBlob] of Object.entries(files)) {
      if (filename.endsWith(".json")) {
        const text = await fileBlob.text();
        assetsFile[filename] = text;
        //console.log(`content of ${filename}`, text);
      } else if (filename.endsWith(".png") || filename.endsWith(".jpg")) {
        const url = URL.createObjectURL(fileBlob);
        assetsFile[filename] = url;
        //console.log(`Image URL for ${filename}`, url);
      }
    }
    return assetsFile;
  }

  async function HandleZipFromS3(presignedURl) {
    try {
      //console.log("Presigned Url: ", presignedURl);
      const zipBlob = await downloadZip(presignedURl);
      const files = await unzipFile(zipBlob);
      const processedFiles = await processFiles(files);
      //console.log("files: ", processedFiles);
      return processedFiles;
    } catch (error) {
      console.log("Error handling zip file:", error);
    }
  }

  async function getAssets() {
    const url = SERVER + "/api/v1/joinroom";

    const data = {
      roomID: spaceId,
    };

    try {
      const res = await axios.post(url, data);
      const files = await HandleZipFromS3(res.data["url"]);
      return files;
    } catch (err) {
      console.error("Error fetching assets: ", err);
      throw err;
    }
  }

  useEffect(() => {
    let game;
    async function start() {
      const files = await getAssets();

      const config = {
        type: Phaser.CANVAS,
        width: sizes.width,
        height: sizes.height,
        canvas: GameCanvas,
        pixelArt: true,
        scale: {
          mode: Phaser.Scale.RESIZE,
          autoCenter: Phaser.Scale.CENTER_BOTH,
        },
        physics: {
          default: "arcade",
          arcade: {
            gravity: { y: 0 },
            debug: true,
          },
        },
        scene: [GameScene],
      };
      game = new Phaser.Game(config);
      game.scene.start("GameScene", { files, spaceId });
    }

    start();

    return () => {
      console.log("destroying game");
      const ws = getSocketInstance()
      if (ws) {
        ws.close();
      }
      game.input.keyboard.destroy();
      game.destroy(true, true);
      //TODO: check if all the assets will be deleted.
    };
  }, []);
  return <canvas id="GameCanvas" />;
};

export default Game;

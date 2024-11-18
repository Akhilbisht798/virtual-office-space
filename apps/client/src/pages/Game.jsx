import Phaser from "phaser";
import GameScene from "../components/Game/GameClass";
import { useEffect } from "react";
import useSocket from "../hooks/useSocket";

const sizes = {
    width: 400,
    height: 300,
}

const Game = () => {
    // const ws = useSocket();
    useEffect(() => {
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
            scene: [GameScene]
        }
        const game = new Phaser.Game(config);

        return () => {
            console.log("destroying game");
            game.input.keyboard.destroy();
            game.destroy(true, true);
        };
    }, []);
    return <canvas id="GameCanvas" />;
};

export default Game;

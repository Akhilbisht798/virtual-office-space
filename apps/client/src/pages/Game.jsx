import Phaser from "phaser";
import GameScene from "../components/Game/GameClass";
import { useEffect } from "react";

const sizes = {
    width: 400,
    height: 300,
}

const Game = () => {
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
            game.destroy(true, true);
            game.input.keyboard.destroy();
        };
    }, []);
    return <canvas id="GameCanvas" />;
};

export default Game;

import { SOCKET_SERVER } from "../global";
import { useEffect, useRef } from "react";

let socketInstance = null;

const useSocket = () => {
    const ws = useRef(null);
    useEffect(() => {
        function init() {
            if (!socketInstance) {
                ws.current = new WebSocket(SOCKET_SERVER);
                ws.current.onopen = () => {
                    console.log("connected to socket server");
                    const userId = "userID-" + Math.floor(Math.random() * 100).toString()
                    console.log(userId)
                    ws.current.send(JSON.stringify({
                        type: "join",
                        payload: {
                            roomId: "8",
                            userId: userId,
                            x: 50,
                            y: 80,
                        },
                    }));
                }
                ws.current.onclose = () => {
                    console.log("disconnected from socket server");
                }
            }
            socketInstance = ws.current;
        }
        init();

        return () => {
            if (socketInstance) {
                ws.current.close();
                socketInstance = null;
            }
        }
    }, []);
    return ws;
}

export default useSocket;

export const getSocketInstance = () => {
    if (!socketInstance) {
        socketInstance = new WebSocket(SOCKET_SERVER)
    }
    // socketInstance.onopen = () => {
    //     console.log("connected to socket server");
    //     const userId = "userID-" + Math.floor(Math.random() * 100).toString()
    //     console.log(userId)
    //     socketInstance.send(JSON.stringify({
    //         type: "join",
    //         payload: {
    //             roomId: "8",
    //             userId: userId,
    //             x: 50,
    //             y: 80,
    //         },
    //     }));

    // }
    // socketInstance.onclose = () => {
    //     console.log("disconnected from socket server");
    // }
    return socketInstance;
}
import { CssBaseline, ThemeProvider } from "@mui/material";
import React from "react";
import ReactDOM from "react-dom/client";
import theme from "./Themes";
import io from "socket.io-client";
import {RouterProvider} from "react-router-dom";
import {BackendResponseMessage, EventEmitter} from "./Helper/backend_events";
import router from "./router";
const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
    <ThemeProvider theme={theme}>
      <CssBaseline />
        <App/>
    </ThemeProvider>
);

function App() {
    const socketIO = io(
        process.env.NODE_ENV === "development" ? "ws://127.0.0.1:8000" : "/", {
            transports: ["websocket"],
            reconnect: true,
        }
    );
    React.useEffect(()=>{
        socketIO.on("connect", () => {});
        let emitter = EventEmitter.getInstance();
        socketIO.on("push", (message:BackendResponseMessage) => {
            console.debug(`监听到push事件 methodName：${message.methodName} data：${message.data}`)
            emitter.emit(message.methodName,message)
        });
        emitter.on("hello",(message: BackendResponseMessage)=>{})
    },[])
    return <RouterProvider router={router} />;
}
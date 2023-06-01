"use client";

import React, { useEffect, useState } from "react";
import { Websocket, WebsocketBuilder} from 'websocket-ts';
import Message from "./message";

let socket: Websocket;

interface Msg {
    username: string;
    message: string;
}

export default function Home() {
    const [message, setMessage] = useState("");
    const [username, setUsername] = useState("");
    let userMessages: Msg[] = []
    const [allMessages, setAllMessages] = useState(userMessages);

    useEffect(() => {
        socketInitializer();
        return () => {
            socket.close();
        };
    }, []);

    async function socketInitializer() {

        socket = new WebsocketBuilder('ws://localhost:8080/ws/123?v=1.0')
            .onOpen((i, ev) => { console.log("opened") })
            .onClose((i, ev) => { console.log("closed") })
            .onError((i, ev) => { console.log("error") })
            .onRetry((i, ev) => { console.log("retry") })
            .onMessage((i,ev) => {
                let msg: Msg = JSON.parse(ev.data);
                setAllMessages(messages => ([...messages, msg]))
                console.log("MSG: ")
                console.log(msg)
            })
            .build();
    }

    function handleSubmit(e: any) {
        e.preventDefault();

        let msg: Msg = {
            username: username,
            message: message,
        }

        console.log("emitted " + msg);
        socket.send(JSON.stringify(msg));

        setMessage("");
    }

  return (
    <main className="container mx-auto px-4 py-8">
          <h1 className="text-2xl font-bold mb-4">Chat app</h1>
          <h1 className="text-lg">Enter a username</h1>
          <input className="dark:bg-gray-900 and dark:text-white w-full px-4 py-2 rounded border border-gray-300 focus:outline-none focus:ring focus:border-blue-500 mb-4" value={username} onChange={(e) => setUsername(e.target.value)} />

          <br />
          <br />

          <div>
              {allMessages.length ? allMessages.map(({ username, message }, index) => (
                Message({index, username, message})
              )): ""}

              <br />

              <form className="flex" onSubmit={handleSubmit}>
                  <input
                      className="dark:bg-gray-900 dark:text-white flex-grow px-4 py-2 rounded-l border border-gray-300 focus:outline-none focus:ring focus:border-blue-500"
                      name="message"
                      placeholder="enter your message"
                      value={message}
                      onChange={(e) => setMessage(e.target.value)}
                      autoComplete={"off"}
                  />
                  <button className="px-4 py-2 bg-blue-500 text-white font-semibold rounded-r hover:bg-blue-600" type="submit">Send</button>
              </form>
          </div>

    </main>
  )
}

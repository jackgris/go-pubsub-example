"use client";

import React, { useEffect, useState } from "react";
import { Websocket, WebsocketBuilder} from 'websocket-ts';

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
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
          <h1>Chat app</h1>
          <h1>Enter a username</h1>
          <input value={username} onChange={(e) => setUsername(e.target.value)} />

          <br />
          <br />

          <div>
              {allMessages.length ? allMessages.map(({ username, message }, index) => (
                  <div key={index}>
                      {username}: {message}
                  </div>
              )): ""}

              <br />

              <form onSubmit={handleSubmit}>
                  <input
                      name="message"
                      placeholder="enter your message"
                      value={message}
                      onChange={(e) => setMessage(e.target.value)}
                      autoComplete={"off"}
                  />
              </form>
          </div>

    </main>
  )
}

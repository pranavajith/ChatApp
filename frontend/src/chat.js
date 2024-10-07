import React, { useEffect, useState } from "react";
import "./chatStyle.css";

const Chat = () => {
  const [username, setUsername] = useState("");
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const [ws, setWs] = useState(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws?username=" + username);

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      setMessages((prevMessages) => [...prevMessages, msg]);
    };

    setWs(socket);

    return () => {
      socket.close();
    };
  }, [username]);

  const sendMessage = (e) => {
    e.preventDefault();
    if (ws && message) {
      ws.send(JSON.stringify({ content: message }));
      setMessage("");
    }
  };

  return (
    <div>
      <h1>Chat Application</h1>
      <input
        type="text"
        placeholder="Enter your username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <div>
        <form onSubmit={sendMessage}>
          <input
            type="text"
            placeholder="Type a message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
          />
          <button type="submit">Send</button>
        </form>
      </div>
      <div>
        <h2>Messages</h2>
        {messages.map((msg, index) => (
          <div key={index}>
            <strong>{msg.username}:</strong> {msg.content}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Chat;

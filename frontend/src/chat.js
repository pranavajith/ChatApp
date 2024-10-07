import React, { useEffect, useState } from "react";
import "./chat.css";

const Chat = () => {
  const [username, setUsername] = useState("");
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const [ws, setWs] = useState(null);
  const [connectedUsers, setConnectedUsers] = useState([]);
  const [typingUsers, setTypingUsers] = useState(new Set());

  const sendMessage = (e) => {
    e.preventDefault();
    if (ws && message) {
      ws.send(JSON.stringify({ content: message }));
      setMessage("");
    }
  };

  useEffect(() => {
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [ws]);

  const setWebSocket = () => {
    setMessages([]);
    if (!username) {
      alert("Input username please!");
      return;
    }
    const socket = new WebSocket("ws://localhost:8080/ws?username=" + username);

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);

      if (msg.message_type === "user_list") {
        const users = msg.content.split(", ");
        setConnectedUsers(users);
      } else if (msg.message_type === "typing") {
        if (msg.username === username) {
          return;
        }
        setTypingUsers((prevTyping) => {
          const newTyping = new Set(prevTyping);
          newTyping.add(msg.username);
          setTimeout(() => {
            newTyping.delete(msg.username);
            setTypingUsers(new Set(newTyping));
          }, 1500);
          return newTyping;
        });
      } else {
        setMessages((prevMessages) => [...prevMessages, msg]);
      }
    };
    socket.onerror = (event) => {
      console.error("WebSocket error:", event);
      alert("Failed to connect to the chat server. Please try again later.");
    };
    setWs(socket);
  };

  const sendTypingIndiactor = () => {
    if (ws) {
      ws.send(JSON.stringify({ message_type: "typing", username: username }));
    }
  };

  return (
    <div className="chat-container">
      <div className="user-list">
        <h2>Active Users</h2>
        <ul>
          {connectedUsers.map((user, index) => (
            <li key={index}>{user}</li>
          ))}
        </ul>
      </div>
      <div className="chat-area">
        <h1>Chat Application</h1>
        <input
          type="text"
          placeholder="Enter your username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              setWebSocket();
            }
          }}
        />
        <button onClick={setWebSocket}>Set username</button>
        <div>
          <form onSubmit={sendMessage}>
            <input
              type="text"
              placeholder="Type a message..."
              value={message}
              onChange={(e) => {
                setMessage(e.target.value);
                sendTypingIndiactor();
              }}
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
          {Array.from(typingUsers).map((user) => (
            <div key={user}>
              <em>{user} is typing...</em>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Chat;

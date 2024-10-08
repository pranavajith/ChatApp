import React, { useEffect, useRef, useState } from "react";
import "./chat.css";

const Chat = () => {
  const [username, setUsername] = useState("");
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const [ws, setWs] = useState(null);
  const [connectedUsers, setConnectedUsers] = useState([]);
  const [reciever, setReciever] = useState("All");
  const [typingUsers, setTypingUsers] = useState(new Set());
  const [darkMode, setDarkMode] = useState(false);

  const messagesEndRef = useRef(null);

  const sendMessage = (e) => {
    e.preventDefault();
    if (ws && message) {
      ws.send(
        JSON.stringify({
          content: message,
          reciever: reciever === "All" ? "" : reciever,
          username: username,
        })
      );
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
      console.log(msg);

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

  const sendTypingIndicator = () => {
    if (ws) {
      const msg = JSON.stringify({
        message_type: "typing",
        username: username,
        reciever: reciever === "All" ? "" : reciever,
      });
      console.log(msg);
      ws.send(msg);
    }
  };

  const toggleTheme = () => {
    setDarkMode((prev) => !prev);
  };

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  return (
    <div className={`chat-container ${darkMode ? "dark" : "light"}`}>
      <div className="user-list">
        <h2 className="active-users-title">Active Users</h2>
        <ul className="user-list-items">
          {connectedUsers.map((user, index) => (
            <li key={index} className="user-list-item">
              {user}
            </li>
          ))}
        </ul>
      </div>
      <div className="chat-area">
        <h1 className="chat-title">Chat Application</h1>
        <button className="theme-toggle-button" onClick={toggleTheme}>
          Switch to {darkMode ? "Light" : "Dark"} Mode
        </button>
        <input
          type="text"
          className="username-input"
          placeholder="Enter your username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              setWebSocket();
            }
          }}
        />
        <button className="set-username-button" onClick={setWebSocket}>
          Set username
        </button>
        <div className="messages-container">
          <h2 className={`messages-title-${darkMode ? "dark" : "light"}`}>
            Messages
          </h2>
          <div className="messages">
            {messages.map((msg, index) => (
              <div
                className={`message-display message-display-${
                  darkMode ? "dark" : "light"
                } ${msg.username === "Server" ? "message-server" : ""} ${
                  msg.reciever === username && msg.reciever ? "message-dm" : ""
                }`}
                key={index}
              >
                <strong className="message-sender">{msg.username}:</strong>{" "}
                <span className="message-content">{msg.content}</span>
              </div>
            ))}
            {Array.from(typingUsers).map((user) => (
              <div
                key={user}
                className={`message-display ${darkMode ? "dark" : "light"}`}
              >
                <em className="typing-indicator">{user} is typing...</em>
              </div>
            ))}
            <div ref={messagesEndRef}></div>
          </div>
        </div>
        <div className="message-input-container">
          <form onSubmit={sendMessage} className="message-form">
            <input
              type="text"
              className="message-input"
              placeholder="Type a message..."
              value={message}
              onChange={(e) => {
                setMessage(e.target.value);
                sendTypingIndicator();
              }}
            />
            <select
              className="recipient-select"
              value={reciever}
              onChange={(e) => setReciever(e.target.value)}
            >
              <option value="All">All</option>
              {connectedUsers.map((user, index) => (
                <option key={index} value={user}>
                  {user}
                </option>
              ))}
            </select>
            <button type="submit" className="send-button">
              Send
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Chat;

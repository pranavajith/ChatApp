import React, { useEffect, useRef, useState } from "react";
import "./chat.css";

const Chat = () => {
  const [username, setUsername] = useState("");
  const [lobbyName, setLobbyName] = useState("");
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const [filteredMessages, setFilteredMessages] = useState([]);
  const [ws, setWs] = useState(null);
  const [connectedUsers, setConnectedUsers] = useState([]);
  const [connectedLobbies, setConnectedLobbies] = useState({});
  const [selectedLobby, setSelectedLobby] = useState("None");
  const [typingUsers, setTypingUsers] = useState(new Set());
  const [darkMode, setDarkMode] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [userFocused, setUserFocused] = useState(true);
  const messagesEndRef = useRef(null);

  const sendMessage = (e) => {
    e.preventDefault();
    if (ws && message) {
      const msg = JSON.stringify({
        content: message,
        username: username,
        message_type: selectedLobby !== "None" ? "lobby_message" : "message",
        lobbyId:
          selectedLobby !== "None" ? connectedLobbies[selectedLobby] : "",
      });
      // console.log("Here is connectedLobbies: ", connectedLobbies);
      console.log("Here is the message sent: ", msg);
      ws.send(msg);
      setMessage("");
    }
  };

  const generateRandomId = () => {
    return Math.floor(10000 + Math.random() * 90000).toString();
  };

  const createLobby = (e) => {
    // e.preventDefault();
    if (ws && lobbyName) {
      ws.send(
        JSON.stringify({
          message_type: "create_lobby",
          username: username,
          lobbyId: generateRandomId().toString(),
          content: lobbyName,
        })
      );
      setLobbyName("");
    }
  };

  const joinLobby = (lobbyId) => {
    if (ws && lobbyId) {
      ws.send(
        JSON.stringify({
          message_type: "join_lobby",
          username: username,
          lobbyId: lobbyId,
        })
      );
      console.log(`Joined lobby: ${lobbyId}`);
    }
  };

  const leaveLobby = (lobbyId) => {
    if (ws && lobbyId) {
      ws.send(
        JSON.stringify({
          message_type: "leave_lobby",
          lobbyId: lobbyId,
        })
      );
      console.log(`Left lobby: ${lobbyId}`);
    }
  };

  const handleLobbyChange = (lobbyName) => {
    if (selectedLobby !== "None") {
      leaveLobby(connectedLobbies[selectedLobby]); // Leave the current lobby
      console.log("Left lobby: ", selectedLobby);
    }
    setSelectedLobby(lobbyName); // Set the new lobby
    if (lobbyName !== "None") {
      joinLobby(connectedLobbies[lobbyName]); // Join the new lobby
      console.log("Joined lobby: ", lobbyName);
    }
  };

  const setWebSocket = () => {
    if (!username) {
      alert("Input username please!");
      return;
    }

    setMessages([]);
    setFilteredMessages([]);

    const socket = new WebSocket(
      "wss://chatapp-yvlc.onrender.com/ws?username=" + username
    );

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      console.log("Here is the message received: ", msg);

      switch (msg.message_type) {
        case "user_list":
          const users = msg.content.split(", ");
          setConnectedUsers(users);
          break;

        case "lobby_list":
          const lobbies = msg.content.split(", ").reduce((acc, lobby) => {
            const [name, id] = lobby.split(":");
            acc[name] = id; // Map lobbyName to lobbyId
            return acc;
          }, {});
          console.log("Here are the available lobbies: ", lobbies);
          setConnectedLobbies(lobbies);
          break;

        case "typing":
          if (msg.username === username) return;
          setTypingUsers((prevTyping) => {
            const newTyping = new Set(prevTyping);
            newTyping.add(msg.username);
            setTimeout(() => {
              newTyping.delete(msg.username);
              setTypingUsers(new Set(newTyping));
            }, 1500);
            return newTyping;
          });
          break;

        default:
          setMessages((prevMessages) => [...prevMessages, msg]);
          setFilteredMessages((prevMessages) => [...prevMessages, msg]);
          if (!userFocused) {
            new Notification(`New message from ${msg.username}`, {
              body: msg.content,
            });
          }
          break;
      }
    };

    socket.onerror = (event) => {
      console.error("WebSocket error:", event);
      alert("Failed to connect to the chat server. Please try again later.");
    };

    setWs(socket);
  };

  const handleSearch = (e) => {
    const query = e.target.value;
    setSearchQuery(query);
    if (query) {
      const filtered = messages.filter(
        (msg) =>
          msg.content.toLowerCase().includes(query.toLowerCase()) ||
          msg.username.toLowerCase().includes(query.toLowerCase())
      );
      setFilteredMessages(filtered);
    } else {
      setFilteredMessages(messages);
    }
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
        <h2 className="active-lobbies-title">Active Lobbies</h2>
        <ul className="lobby-list-items">
          {Object.entries(connectedLobbies).map(([name, id]) => (
            <li key={id} className="lobby-list-item">
              {name}
            </li>
          ))}
        </ul>
      </div>

      <div className="chat-area">
        <h1 className="chat-title">Chat Application</h1>
        <button
          className="theme-toggle-button"
          onClick={() => setDarkMode(!darkMode)}
        >
          Switch to {darkMode ? "Light" : "Dark"} Mode
        </button>

        <input
          type="text"
          className="username-input"
          placeholder="Enter your username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") setWebSocket();
          }}
        />
        <button className="set-username-button" onClick={setWebSocket}>
          Set username
        </button>

        <input
          type="text"
          className="lobby-input"
          placeholder="Create a lobby"
          value={lobbyName}
          onChange={(e) => setLobbyName(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") createLobby();
          }}
        />
        <button className="create-lobby-button" onClick={createLobby}>
          Create Lobby
        </button>

        <select
          className="lobby-select"
          value={selectedLobby}
          onChange={(e) => handleLobbyChange(e.target.value)}
        >
          <option value="None">None</option>
          {Object.entries(connectedLobbies).map(([name, id]) => (
            <option key={id} value={name}>
              {name}
            </option>
          ))}
        </select>

        <input
          type="text"
          className="search-input"
          placeholder="Search messages..."
          value={searchQuery}
          onChange={handleSearch}
        />

        <div className="messages-container">
          <div className="messages">
            {filteredMessages.map((msg, index) => (
              <div
                className={`message-display message-display-${
                  darkMode ? "dark" : "light"
                }`}
                key={index}
              >
                <strong>{msg.username}:</strong> {msg.content}
              </div>
            ))}
            {Array.from(typingUsers).map((user) => (
              <div
                key={user}
                className={`message-display ${darkMode ? "dark" : "light"}`}
              >
                <em>{user} is typing...</em>
              </div>
            ))}
            <div ref={messagesEndRef}></div>
          </div>
        </div>

        <form onSubmit={sendMessage} className="message-form">
          <input
            type="text"
            className="message-input"
            placeholder="Type a message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
          />
          <button type="submit" className="send-button">
            Send
          </button>
        </form>
      </div>
    </div>
  );
};

export default Chat;

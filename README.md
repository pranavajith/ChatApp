# Chat App

Welcome to my first **Chat App**!  
It may be minimal and have some quirks, but it’s a start, and I’m excited to share it with you. 🎉

### [Deployed Site](https://chat-app-ebon-eta.vercel.app/)

Note: If you are loading it for the first time, you may have to give a good 1-5 mins after setting the username to see changes in the UI, as the backend server is running on the free model, and hence would take some time to boot up after loading off after 15 mins of inactivity. Thanks!

## Features in Progress

I have many plans to make this chat app better in the future. Below is a list of some enhancements and improvements I hope to implement soon.

### 1. **User Experience Enhancements**

- **User List**: Display a list of connected users so participants can see who is online. (Update: Done ✅)
- **Typing Indicator**: Show when another user is typing to improve real-time interaction. (Update: Done ✅)
- **Message Read Receipts**: Indicate when a message has been read by the recipient.
- **Message Recievers**: Allow users to pick whom to send a message to. (Update: Done ✅)
- **Emojis and Reactions**: Allow users to send emojis or react to messages for fun and better communication.
- **Dark Mode/Light Mode**: Provide a theme switcher for better visual comfort based on user preference. (Update: Done ✅)
- **Customizable User Profiles**: Allow users to set profile pictures and status messages for personalization.

### 2. **Functional Features**

- **Message History**: Store chat history on the server-side and retrieve it when users connect. (Update: Done ✅)
- **Private Messaging**: Enable direct messaging between users for private conversations. (Update: Done ✅)
- **File Sharing**: Support sending images, documents, and other file types within the chat.
- **Search Functionality**: Allow users to search through past messages for easy access to information. (Update: Done ✅)
- **Notification System**: Implement browser notifications for new messages, even when the app is not in focus. (Update: Done ✅)
- **Deploy Website Online**: Deploy the Frontend and Backend online

### 3. **Security Features**

- **Authentication**: Implement user authentication (e.g., login system) to restrict access to authorized users.
- **Input Sanitization**: Sanitize user inputs to prevent malicious attacks like Cross-Site Scripting (XSS).
- **TLS/SSL Encryption**: Secure WebSocket communication using `wss://` (WebSocket Secure) for data integrity and protection.

### 4. **Server-Side Improvements**

- **User Rooms/Channels**: Create different chat rooms or channels where users can join specific conversations.
- **Admin Controls**: Provide admin functionality to manage users, moderate the chat, or control rooms.
- **Rate Limiting**: Prevent spam by limiting the rate at which users can send messages.
- **Analytics Dashboard**: Track app usage statistics such as the number of messages sent and active users.

### 5. **Integration with Third-Party Services**

- **Bots and Integrations**: Integrate bots for automated responses or helpful functionalities (e.g., weather updates, news).
- **APIs**: Connect external APIs (like a translation service) to enable multilingual chat support.
- **Payment Processing**: If applicable, integrate payment systems for in-app transactions, such as sending virtual gifts.

### 6. **Accessibility Features**

- **Screen Reader Support**: Ensure the app is accessible to users with visual impairments.
- **Keyboard Navigation**: Allow full navigation through the app using keyboard shortcuts for ease of use.

---

## Running the Application

To run the application, follow these steps:

1. Clone this repository.
2. Make sure you have Go and Node.js installed.
3. Run the backend and frontend simultaneously using the provided Makefile:
   ```bash
   make run
   ```
   This will start both the **Go backend** and the **React frontend** concurrently.

Alternatively, you can run the frontend and backend individually:

- To run the backend (Go server):

  ```bash
  make run-backend
  ```

- To run the frontend (React app):
  ```bash
  make run-frontend
  ```

The app will be available at `http://localhost:3000`, and the WebSocket server runs on `ws://localhost:8080/ws`.

---

## Contributions

Feel free to fork this project and contribute. I’m open to feedback, bug reports, and feature requests!

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

---

Thank you for checking out my Chat App! Stay tuned for more updates and exciting new features. 😊

---

### How the Makefile Works

- **`make run`**: This target starts both the backend and frontend in parallel.
- **`make run-backend`**: This starts only the backend (Go) server.
- **`make run-frontend`**: This starts only the frontend (React) app.
- **`-j2` option**: This allows two jobs (frontend and backend) to run in parallel, speeding up the launch process.

This structure ensures your backend and frontend components are deployed seamlessly together using a single Makefile command!

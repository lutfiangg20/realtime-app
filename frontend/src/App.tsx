import { useEffect, useRef, useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";

function App() {
  const [messages, setMessages] = useState<string[]>([]);
  const scroll = useRef<HTMLDivElement>(document.createElement("div"));

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:3000/ws");

    socket.onmessage = (event: MessageEvent) => {
      setMessages((prev) => [...prev, event.data]);
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    return () => {
      socket.close();
    };
  }, []);

  useEffect(() => {
    // Scroll ke bagian paling bawah saat ada data baru
    if (scroll.current) {
      scroll.current.scrollTop = scroll.current.scrollHeight;
    }
  }, [messages]);

  return (
    <div>
      <h1>Real-time Command Output</h1>
      <div
        ref={scroll}
        style={{ overflow: "auto", height: "300px", backgroundColor: "gray" }}
      >
        <pre style={{ textDecorationColor: "black" }}>
          {messages.join("\n")}
        </pre>
      </div>
    </div>
  );
}

export default App;

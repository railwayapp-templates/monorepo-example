import { Fragment, useCallback, useEffect, useState } from "react";

import { mc } from "./assets/mc";

import './App.css'

function App() {
  const [messages, setMessages] = useState<string[]>([]);
  const [isConnectionOpen, setIsConnectionOpen] = useState(false);

  const onToggleConnection = useCallback(() => {
      setIsConnectionOpen((isOpen) => !isOpen);
  }, []);

  useEffect(() => {
      if (!isConnectionOpen) return;

      const eventSource = new EventSource(import.meta.env.VITE_BACKEND_HOST + "/sse");

      eventSource.onopen = () => {
          console.log("[SSE] Connection established");
      };

      eventSource.onmessage = (event) => {
          setMessages((messages) => [...messages, event.data]);
      };

      eventSource.onerror = (event) => {
          console.error("[SSE] Error:", event);

          if (eventSource.readyState === EventSource.CLOSED) {
              console.log("[SSE] Connection closed because of an error");
              setIsConnectionOpen(false);
          }
      };

      const cleanup = () => {
          console.log("[SSE] Closing connection");
          eventSource.close();
          window.removeEventListener("beforeunload", cleanup);
      };

      window.addEventListener("beforeunload", cleanup);

      return cleanup;
  }, [isConnectionOpen]);

  useEffect(() => {
      window.scrollTo({
          top: document.body.scrollHeight,
          behavior: "smooth",
      });
  }, [messages]);

  return (
      <div className="mx-auto flex size-full flex-col gap-4 p-10 text-center text-white">
          <h1 className="text-4xl font-semibold">Here's some unnecessary quotes for you to read...</h1>

          {messages.map((message, index, elements) => (
              <Fragment key={index}>
                  <p className={mc("duration-200", index + 1 !== elements.length ? "opacity-40" : "scale-105 font-bold")}>{message}</p>
              </Fragment>
          ))}

          {/* eslint-disable-next-line tailwindcss/no-arbitrary-value */}
          <button className={mc("hover:opacity-75 duration-200 font-bold text-lg", isConnectionOpen ? "text-[#f06b6b]" : "text-[#6bf06b]")} onClick={onToggleConnection}>
              {isConnectionOpen ? "Stop" : "Start"} Quotes
          </button>

          <div className="h-96 w-full" />
      </div>
  );
}

export default App

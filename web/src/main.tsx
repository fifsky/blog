import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import "nprogress/nprogress.css";
import "./index.css";
import { ErrorBoundary } from "@/components/ErrorBoundary";

dayjs.extend(isBetween);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ErrorBoundary>
      <App />
    </ErrorBoundary>
  </React.StrictMode>,
);

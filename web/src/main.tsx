import React from "react";
import ReactDOM from "react-dom/client";
import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import "nprogress/nprogress.css";
import "./index.css";
import { ErrorBoundary } from "@/components/ErrorBoundary";
import { Toaster } from "sonner";
import { RouterProvider } from "react-router";
import { router } from "./router";

dayjs.extend(isBetween);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ErrorBoundary>
      <Toaster position="top-center" duration={3000} />
      <RouterProvider router={router} />
    </ErrorBoundary>
  </React.StrictMode>,
);

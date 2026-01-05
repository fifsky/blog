import { RouterProvider } from "react-router";
import { router } from "./router";
import { Toaster } from "@/components/ui/sonner";

export default function App() {
  return (
    <>
      <Toaster position="top-center" duration={3000} />
      <RouterProvider router={router} />
    </>
  );
}

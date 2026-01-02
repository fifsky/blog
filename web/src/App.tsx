import { RouterProvider } from "react-router-dom";
import { router } from "./router";
import { StoreProvider } from "./store/context";

export default function App() {
  return (
    <StoreProvider>
      <RouterProvider router={router}></RouterProvider>
    </StoreProvider>
  );
}

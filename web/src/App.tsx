import { Outlet } from "react-router";
import { Suspense } from "react";
import { RouteProgress } from "@/components/RouteProgress";
import Loading from "@/components/Loading";

export default function App() {
  return (
    <Suspense fallback={<Loading />}>
      <RouteProgress />
      <Outlet />
    </Suspense>
  );
}

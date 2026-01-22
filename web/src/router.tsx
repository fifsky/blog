import { lazy, Suspense } from "react";
import { createBrowserRouter, type RouteObject } from "react-router";
import { Layout } from "@/components/Layout";
import { AdminLayout } from "@/components/AdminLayout";
import { RouteErrorBoundary } from "@/components/RouteErrorBoundary";
import App from "./App";
import { SkeletonArticle, SkeletonArticleList } from "@/components/Skeleton";

const ArticleList = lazy(() => import("@/pages/ArticleList"));
const ArticleDetail = lazy(() => import("@/pages/ArticleDetail"));
const TravelMap = lazy(() => import("@/pages/TravelMap"));
const Archive = lazy(() => import("@/pages/Archive"));
const About = lazy(() => import("@/pages/About"));
const Login = lazy(() => import("@/pages/Login"));
const AdminIndex = lazy(() => import("@/pages/admin/AdminIndex"));
const AdminArticle = lazy(() => import("@/pages/admin/AdminArticle"));
const AdminComment = lazy(() => import("@/pages/admin/AdminComment"));
const AdminMood = lazy(() => import("@/pages/admin/AdminMood"));
const AdminCate = lazy(() => import("@/pages/admin/AdminCate"));
const AdminLink = lazy(() => import("@/pages/admin/AdminLink"));
const AdminRemind = lazy(() => import("@/pages/admin/AdminRemind"));
const AdminUser = lazy(() => import("@/pages/admin/AdminUser"));
const AdminPhoto = lazy(() => import("@/pages/admin/AdminPhoto"));
const PostArticle = lazy(() => import("@/pages/admin/PostArticle"));
const PostUser = lazy(() => import("@/pages/admin/PostUser"));

const routesConfig: RouteObject[] = [
  {
    element: <App />,
    errorElement: <RouteErrorBoundary />,
    children: [
      {
        path: "/",
        element: <Layout />,
        children: [
          {
            index: true,
            element: (
              <Suspense fallback={<SkeletonArticleList />}>
                <ArticleList />
              </Suspense>
            ),
          },
          {
            path: "search",
            element: (
              <Suspense fallback={<SkeletonArticleList />}>
                <ArticleList />
              </Suspense>
            ),
          },
          {
            path: "about",
            element: (
              <Suspense fallback={<SkeletonArticle />}>
                <About />
              </Suspense>
            ),
          },
          {
            path: "archive",
            element: (
              <Suspense fallback={<SkeletonArticleList />}>
                <Archive />
              </Suspense>
            ),
          },
          {
            path: "travel",
            element: (
              <Suspense fallback={<SkeletonArticleList />}>
                <TravelMap />
              </Suspense>
            ),
          },
          {
            path: "date/:year/:month",
            element: (
              <Suspense fallback={<SkeletonArticleList />}>
                <ArticleList />
              </Suspense>
            ),
          },
          {
            path: "date/:year/:month/:day",
            element: (
              <Suspense fallback={<SkeletonArticleList />}>
                <ArticleList />
              </Suspense>
            ),
          },
          {
            path: "category/:domain",
            element: (
              <Suspense fallback={<SkeletonArticleList />}>
                <ArticleList />
              </Suspense>
            ),
          },
          {
            path: "article/:id",
            element: (
              <Suspense fallback={<SkeletonArticle />}>
                <ArticleDetail />
              </Suspense>
            ),
          },
        ],
      },
      { path: "/login", element: <Login /> },
      {
        path: "/admin",
        element: <AdminLayout />,
        children: [
          { path: "index", element: <AdminIndex /> },
          { path: "articles", element: <AdminArticle /> },
          { path: "post/article", element: <PostArticle /> },
          { path: "comments", element: <AdminComment /> },
          { path: "moods", element: <AdminMood /> },
          { path: "cates", element: <AdminCate /> },
          { path: "links", element: <AdminLink /> },
          { path: "remind", element: <AdminRemind /> },
          { path: "users", element: <AdminUser /> },
          { path: "photos", element: <AdminPhoto /> },
          { path: "post/user", element: <PostUser /> },
        ],
      },
    ],
  },
];

export const router = createBrowserRouter(routesConfig);

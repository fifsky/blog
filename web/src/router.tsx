import { useEffect, lazy, Suspense } from "react";
import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  Outlet,
  useLocation,
} from "react-router-dom";
import { Layout } from "./components/Layout";
import { AdminLayout } from "./components/AdminLayout";
const ArticleList = lazy(() =>
  import("./pages/ArticleList").then((m) => ({ default: m.ArticleList }))
);
const ArticleDetail = lazy(() =>
  import("./pages/ArticleDetail").then((m) => ({ default: m.ArticleDetail }))
);
const About = lazy(() =>
  import("./pages/About").then((m) => ({ default: m.About }))
);
const Login = lazy(() =>
  import("./pages/Login").then((m) => ({ default: m.Login }))
);
const AdminIndex = lazy(() =>
  import("./pages/AdminIndex").then((m) => ({ default: m.AdminIndex }))
);
const AdminArticle = lazy(() =>
  import("./pages/AdminArticle").then((m) => ({ default: m.AdminArticle }))
);
const AdminComment = lazy(() =>
  import("./pages/AdminComment").then((m) => ({ default: m.AdminComment }))
);
const AdminMood = lazy(() =>
  import("./pages/AdminMood").then((m) => ({ default: m.AdminMood }))
);
const AdminCate = lazy(() =>
  import("./pages/AdminCate").then((m) => ({ default: m.AdminCate }))
);
const AdminLink = lazy(() =>
  import("./pages/AdminLink").then((m) => ({ default: m.AdminLink }))
);
const AdminRemind = lazy(() =>
  import("./pages/AdminRemind").then((m) => ({ default: m.AdminRemind }))
);
const AdminUser = lazy(() =>
  import("./pages/AdminUser").then((m) => ({ default: m.AdminUser }))
);
const PostArticle = lazy(() =>
  import("./pages/PostArticle").then((m) => ({ default: m.PostArticle }))
);
const PostUser = lazy(() =>
  import("./pages/PostUser").then((m) => ({ default: m.PostUser }))
);
function useTitleTemplate(title?: string) {
  const location = useLocation();
  useEffect(() => {
    const base = "無處告別";
    document.title = title ? `${title} - ${base}` : base;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname, title]);
}

function TitleWrapper() {
  useTitleTemplate();
  return (
    <Suspense fallback={<div style={{ padding: 16 }}>页面加载中...</div>}>
      <Outlet />
    </Suspense>
  );
}

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<TitleWrapper />}>
      <Route path="/" element={<Layout />}>
        <Route index element={<ArticleList />} />
        <Route path="search" element={<ArticleList />} />
        <Route path="about" element={<About />} />
        <Route path="date/:year/:month" element={<ArticleList />} />
        <Route path="categroy/:domain" element={<ArticleList />} />
        <Route path="article/:id" element={<ArticleDetail />} />
      </Route>
      <Route path="/login" element={<Login />} />
      <Route path="/admin" element={<AdminLayout />}>
        <Route path="index" element={<AdminIndex />} />
        <Route path="articles" element={<AdminArticle />} />
        <Route path="post/article" element={<PostArticle />} />
        <Route path="comments" element={<AdminComment />} />
        <Route path="moods" element={<AdminMood />} />
        <Route path="cates" element={<AdminCate />} />
        <Route path="links" element={<AdminLink />} />
        <Route path="remind" element={<AdminRemind />} />
        <Route path="users" element={<AdminUser />} />
        <Route path="post/user" element={<PostUser />} />
      </Route>
    </Route>
  )
);

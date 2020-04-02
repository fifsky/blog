import Vue from 'vue'
import Router from 'vue-router'
import {AdminLayout, Layout} from "../components"
import {
  About,
  AdminArticle,
  AdminCate,
  AdminComment,
  AdminIndex,
  AdminLink,
  AdminMood,
  AdminRemind, AdminUser,
  ArticleDetail,
  ArticleList,
  Login, PostUser
} from "../pages"
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

Vue.use(Router)

const routes = [
  {
    path: '/',
    component: Layout,
    children: [
      {
        path: '',
        component: ArticleList,
      },
      {
        path: 'search',
        name: 'search',
        component: ArticleList,
      },
      {
        path: 'about',
        component: About,
      },
      {
        path: "date/:year/:month",
        component: ArticleList,
      },
      {
        path: "categroy/:domain",
        component: ArticleList,
      },
      {
        path: "article/:id",
        component: ArticleDetail,
      }
    ],
  },
  {
    path: "/login",
    component: Login
  },
  {
    path: "/admin",
    component: AdminLayout,
    children: [
      {
        path: 'index',
        component: AdminIndex
      },
      {
        path: 'articles',
        component: AdminArticle
      },
      {
        path: 'post/article',
        component: () => import("../pages/PostArticle")
      },
      {
        path: 'comments',
        component: AdminComment
      },
      {
        path: 'moods',
        component: AdminMood
      },
      {
        path: 'cates',
        component: AdminCate
      },
      {
        path: 'links',
        component: AdminLink
      },
      {
        path: 'remind',
        component: AdminRemind
      },
      {
        path: 'users',
        component: AdminUser
      },
      {
        path: 'post/user',
        component: PostUser
      }
    ]
  }
]

const router = new Router({
  routes,
  mode: 'history',
  scrollBehavior(to, from, savedPosition) {
    return {
      x: 0,
      y: 0
    }
  }
})


router.beforeEach((to, from, next) => {
  NProgress.start();
  next()
});

router.afterEach((to, from) => {
  NProgress.done();
});

export default router
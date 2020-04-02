<template>
    <div id="container">
        <CHeader/>
        <div class="admin" v-if="isLogin">
            <div class="tabs">
                <ul>
                    <li><router-link to="/admin/index" :class="{active:isPage('/admin/index')}">设置</router-link></li>
                    <li><router-link to="/admin/articles" :class="{active:isPage('/admin/articles','/admin/post/article')}">文章</router-link></li>
                    <li><router-link to="/admin/comments" :class="{active:isPage('/admin/comments')}">评论</router-link></li>
                    <li><router-link to="/admin/moods" :class="{active:isPage('/admin/moods')}">心情</router-link></li>
                    <li><router-link to="/admin/cates" :class="{active:isPage('/admin/cates')}">分类</router-link></li>
                    <li><router-link to="/admin/links" :class="{active:isPage('/admin/links')}">链接</router-link></li>
                    <li><router-link to="/admin/remind" :class="{active:isPage('/admin/remind')}">提醒</router-link></li>
                    <li><router-link to="/admin/users" :class="{active:isPage('/admin/users')}">用户</router-link></li>
                </ul>
            </div>
            <div id="content">
                <router-view/>
            </div>
        </div>
        <CFooter/>
    </div>
</template>

<script>

  import {mapActions,mapGetters} from "vuex"
  import {CFooter, CHeader} from './index'
  import {getAccessToken, sync} from "../utils";

  export default {
    name: 'AdminLayout',
    components: {
      CHeader,
      CFooter,
    },
    computed:{
      ...mapGetters(["isLogin"])
    },
    methods: {
      ...mapActions(["currentUserAction"]),
      isPage: function () {
        let is = false
        for (let i = 0; i < arguments.length; i++) {
          if(this.$route.path === arguments[i]){
            is = true
            break
          }
        }
        return is
      },
    },
    mounted() {
      if(getAccessToken()) {
        sync(async () => {
          try {
            await this.currentUserAction()
          } catch (e) {
            this.$router.push("/login")
          }
        })
      }
    },
    beforeRouteEnter(to, form, next) {
      if (!getAccessToken()) {
          next("/login")
      } else {
          next()
      }
    }
  }
</script>
<style>

</style>
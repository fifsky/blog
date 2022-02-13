<template>
    <div id="container">
        <CHeader/>
        <div id="main">
            <Mood/>
            <div class="tabs">
                <ul>
                    <li>
                        <router-link to="/about" :class="{active:isPage('/about')}">关于我</router-link>
                    </li>
                    <li>
                        <router-link to="/" :class="{active:isPage('/')}">所有文章</router-link>
                    </li>
                </ul>
            </div>
            <div id="content">
                <transition name="fade" appear>
                <router-view :key="key"/>
                </transition>
            </div>
        </div>
        <Sidebar/>
        <CFooter/>
    </div>
</template>

<script>
  import {sync,getAccessToken} from "../utils";
  import {mapActions} from 'vuex'
  import {CFooter, CHeader, Mood, Sidebar} from './index'

  export default {
    name: 'Layout',
    components: {
      CHeader,
      CFooter,
      Sidebar,
      Mood
    },
    computed: {
      key() {
        return this.$route.fullPath
      }
    },
    methods: {
      ...mapActions(['currentUserAction']),
      isPage: function (path) {
        return this.$route.path === path
      },
    },
    mounted() {
      if(getAccessToken()){
        sync(async () => {
          try {
            await this.currentUserAction()
          } catch (e) {
            localStorage.removeItem("access_token")
          }
        })
      }
    }
  }
</script>
<style>
    .hljs-keyword, .hljs-request, .hljs-status, .hljs-subst, .hljs-winutils, .nginx .hljs-title, .hljs-id, .hljs-title, .scss .hljs-preprocessor, .hljs-class .hljs-title, .hljs-type, .tex .hljs-command, .vhdl .hljs-literal, .hljs-cdata, .hljs-doctype, .hljs-pi, .hljs-pragma, .hljs-preprocessor, .hljs-shebang {
        font-weight: 700
    }
    .fade-enter-active, .fade-leave-active {
        transition: opacity .3s;
    }
    .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
        opacity: 0;
    }
</style>

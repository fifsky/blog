<template>
    <div id="header" class="clearfix">
        <h1>
            <router-link to="/"><img alt=首页 src="../assets/images/logo.png"><span>首页</span></router-link>
        </h1>

        <div id="navigation">
            <ul>
                <li><router-link to="/">首页</router-link></li>
                <li><a href="https://github.com/fifsky" target="_blank">Github</a></li>
                <li><a href="http://blog.fifsky.com/">技术</a></li>
                <li><router-link to="/about">关于</router-link></li>
                <li><a href="https://caixudong.com">简历</a></li>
                <li><router-link to="/admin/index" v-if="isLogin">管理中心</router-link></li>
                <li><a v-if="isLogin" @click="logOut">退出</a></li>
                <li><router-link to="/login" v-if="!isLogin">登录</router-link></li>
            </ul>
        </div>
    </div>
</template>

<script>
  import {mapGetters,mapMutations} from "vuex"

  export default {
    name: "CHeader",
    computed:{
      ...mapGetters(["isLogin"])
    },
    methods:{
      ...mapMutations(["setUserInfo"]),
      logOut () {
        localStorage.removeItem('access_token')
        this.setUserInfo({})
        this.$router.push('/')
      },
    }
  }
</script>
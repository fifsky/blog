<template>
    <div id="container">
        <CHeader/>
        <div class="admin">
            <div id="content">
                <div id="sign-in">
                    <h2>博客管理登录</h2>
                    <div class="message"></div>
                    <form method="post" @submit.prevent="submit" class="vf lf">
                        <p>
                            <label class="label_input">用户名：</label>
                            <input type="text" class="input_text" v-model="formdata.user_name"/>
                        </p>
                        <p>
                            <label class="label_input">密码：</label>
                            <input type="password" name="user_pass" class="input_text" v-model="formdata.password"/>
                        </p>
                        <p>
                            <label for="auto_login" class="label_check">
                                <input type="checkbox" v-model="formdata.auto_login" class="input_check"
                                       name="auto_login"/> 下次自动登录
                            </label>
                        </p>
                        <p class="act">
                            <input type="submit" class="formbutton" value="登录"/>
                        </p>
                    </form>
                </div>
            </div>
        </div>
        <CFooter/>
    </div>
</template>

<script>
  import {sync} from "../utils";
  import {mapActions} from 'vuex'
  import {CFooter, CHeader} from "../components";

  export default {
    name: "Login",
    data() {
      return {
        formdata: {}
      }
    },
    components: {
      CHeader,
      CFooter,
    },
    methods: {
      ...mapActions(['loginAction']),
      submit() {
        sync(async () => {
          await this.loginAction(this.formdata)
          this.$router.push("/admin/index")
        })
      }
    }
  }
</script>

<style scoped>

</style>
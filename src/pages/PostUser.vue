<template>
    <div>
        <h2>{{ user.id ? '编辑':'新增' }}用户
            <router-link to="/admin/users"><i class="iconfont icon-undo" style="color: #444"></i>返回列表</router-link>
        </h2>
        <div class="message"></div>
        <form class="vf" method="post" autocomplete="off" @submit.prevent="submit">
            <p><label class="label_input">用户名 <span class="desc">(必填)</span></label>
                <input class="input_text" size="50" v-model="user.name"></p>
            <span class="hint">此用户名将作为用户登录时所用的名称，请不要与系统中现有的用户名重复。</span>
            <p>

            <p><label class="label_input">邮箱 <span class="desc">(必填)</span></label>
                <input class="input_text" size="50" v-model="user.email"></p>
            <span class="hint">电子邮箱地址将作为此用户的主要联系方式，请不要与系统中现有的电子邮箱地址重复。</span>
            <p>

            <p><label class="label_input">昵称</label>
                <input class="input_text" size="50" v-model="user.nick_name"></p>
            <span class="hint">用户昵称可以与用户名不同, 用于前台显示，如果你将此项留空，将默认使用用户名。</span>

            <p>
                <label class="label_input">密码 <span class="desc">(必填)</span></label>
                <input type="password" class="input_text" size="50" value="" v-model="user.password1">
            </p>
            <span class="hint">为用户分配一个密码。</span>

            <p><label class="label_input">确认密码 <span class="desc">(必填)</span></label>
                <input type="password" class="input_text" value="" size="50" v-model="user.password2"/>
                <span class="hint">请确认你的密码，与上面输入的密码保持一致。</span>
            </p>

            <p><label class="label_input">角色 <span class="desc">(必填)</span></label>
                <select name="type" v-model="user.type">
                    <option value="1">管理员</option>
                    <option value="2">编辑</option>
                </select>
                <span class="hint">管理员具有所有操作权限，编辑仅能包含文章、评论、心情的操作权限。</span>
            </p>

            <p class="act"><input class="formbutton" type="submit" value="保存"></p>
        </form>

    </div>
</template>

<script>
  import {sync} from "../utils";
  import {userGetApi, userCreateApi, userUpdateApi} from "../service";

  export default {
    name: "PostUser",
    data() {
      return {
        user: {
          type: 1
        }
      }
    },
    methods: {
      submit() {
        if(this.user.password1 !== this.user.password2){
          this.$message.error("两次输入的密码不一致")
          return
        }

        let {id, name, nick_name, password1, email, type} = this.user
        let data = {id, name, nick_name, password:password1, email, type:parseInt(type)}
        sync(async () => {
          let ret = await (id ? userUpdateApi : userCreateApi)(data)
          this.$router.push("/admin/users")
        })
      }
    },
    mounted() {
      sync(async () => {
        if (this.$route.query && this.$route.query.id) {
          this.user = await userGetApi({id: parseInt(this.$route.query.id)})
        }
      })
    }
  }
</script>

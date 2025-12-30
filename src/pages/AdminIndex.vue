<template>
    <div id="settings">
        <h2>站点设置</h2>
        <div class="message">保存成功</div>
        <form class="nf" method="post" autocomplete="off" @submit.prevent="submit">
            <p>
                <label class="label_input">站点名称</label>
                <input type="text" class="input_text" size="50" name="site_name" v-model="formdata.site_name"/>
                <span class="hint">站点的名称将显示在网页的标题处。</span>
            </p>
            <p>
                <label class="label_input">站点描述</label>
                <textarea name="site_desc" rows="3" cols="50" v-model="formdata.site_desc"></textarea>
                <span class="hint">站点描述将显示在网页代码的头部。</span>
            </p>
            <p>
                <label class="label_input">关键字</label>
                <input type="text" class="input_text" size="50" name="site_keyword" v-model="formdata.site_keyword"/>
                <span class="hint">请以半角逗号","分割多个关键字。</span>
            </p>

            <p>
                <label class="label_input">每页显示文章数</label>
                <input class="input_text" style="width: 50px" name="post_num" type="text" v-model="formdata.post_num"/>
            </p>

            <p class="act"><input class="formbutton" type="submit" value=保存></p>
        </form>
    </div>
</template>

<script>
  import {sync} from "../utils";
  import {settingApi, settingUpdateApi} from "../service";

  export default {
    name: "AdminIndex",
    data() {
      return {
        formdata: {}
      }
    },
    methods: {
      submit() {
        sync(async () => {
          await settingUpdateApi(this.formdata)
          this.$message.success("保存成功")
        })
      }
    },
    mounted() {
      sync(async () => {
        this.formdata = await settingApi()
      })
    }
  }
</script>

<style scoped>

</style>

<template>
    <div class="clearfix">
        <h2>管理链接</h2>

        <div class="col-left">
            <div class="operate clearfix">
                <BatchHandle/>
            </div>

            <table class="list">
                <tbody>
                <tr>
                    <th width="20">&nbsp;</th>
                    <th width="80">连接名</th>
                    <th>地址</th>
                    <th width="80">操作</th>
                </tr>
                <tr v-if="list.length === 0">
                    <td colspan="7" align="center">还没有链接！</td>
                </tr>
                <tr v-if="list.length > 0" v-for="(v,k) in list" :key="v.id">
                    <td><input type="checkbox" name="ids" :value="v.id"/></td>
                    <td><a :href="v.url" target="_blank">{{v.name}}</a></td>
                    <td><a :href="v.url">{{v.url}}</a></td>
                    <td><a href="javascript:void(0)" @click="editItem(v.id)">编辑</a><span class="line">|</span><a href="javascript:void(0)" @click="deleteItem(v.id)">删除</a>
                    </td>
                </tr>
                </tbody>
            </table>
            <div class="operate clearfix">
                <BatchHandle/>
            </div>
        </div>
        <div class="col-right" style="width: 250px; padding-top: 31px;">
            <form class="vf" method="post" autocomplete="off" @submit.prevent="submit">
                <p><label class="label_input">链接名称</label>
                    <input class="input_text" size="30" name="name"
                           v-model="item.name"></p>

                <p><label class="label_input">链接地址</label>
                    <input class="input_text" size="30" name="url"
                           v-model="item.url">
                    <span class="hint">例如：http://fifsky.com/</span>
                </p>
                <p>
                    <label class="label_input">链接描述</label>
                    <textarea name="desc" rows="5" cols="30" v-model="item.desc"></textarea>
                </p>
                <p class="act">
                    <button class="formbutton" type="submit">{{item.id ? '修改':'添加'}}</button><a v-show="item.id" class="ml10" href="javascript:void(0)" @click="cancel">取消</a>
                </p>
            </form>
        </div>
    </div>
</template>

<script>
  import {linkDeleteApi, linkListApi, linkPostApi} from "../service";
  import {BatchHandle} from "../components";
  import list from "../mixins/list"

  export default {
    name: "AdminLink",
    data() {
      return {
        listApi: linkListApi,
        postApi: linkPostApi,
        deleteApi:linkDeleteApi
      }
    },
    mixins: [list],
    components: {
      BatchHandle
    },
    methods: {
      submit() {
        let {id, name,url,desc} = this.item
        this.triggerSubmit({id, name,url,desc})
      },
    },
    mounted() {
      this.loadList()
    }
  }
</script>

<style scoped>

</style>
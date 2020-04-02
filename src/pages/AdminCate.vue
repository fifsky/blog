<template>
    <div class="clearfix">
        <h2>管理分类</h2>

        <div class="col-left">
            <div class="operate clearfix">
                <BatchHandle/>
            </div>

            <table class="list">
                <tbody>
                <tr>
                    <th width="20">&nbsp;</th>
                    <th>分类名</th>
                    <th width="60">缩略名</th>
                    <th width="50">文章数</th>
                    <th width="80">操作</th>
                </tr>
                <tr v-if="list.length === 0">
                    <td colspan="7" align="center">还没有分类！</td>
                </tr>
                <tr v-if="list.length > 0" v-for="(v,k) in list" :key="v.id">
                    <td><input type="checkbox" name="ids" :value="v.id"/></td>
                    <td>{{v.name}}</td>
                    <td>{{ v.domain }}</td>
                    <td class="art-num">{{v.num}}</td>
                    <td><a href="javascript:void(0)" @click="editItem(v.id)">编辑</a><span class="line">|</span><a href="javascript:void(0)"
                                                                                      @click="deleteItem(v.id)">删除</a>
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
                <p><label class="label_input">分类名称</label>
                    <input type="text" class="input_text" size="30" name="name" v-model="item.name"/></p>

                <p><label class="label_input">分类缩略名</label>
                    <input type="text" class="input_text" size="30" name="domain" v-model="item.domain"/>
                    <span class="hint">缩略名，使用字母开头([a-z][0-9]-)</span>
                </p>
                <p>
                    <label class="label_input">分类描述</label>
                    <textarea name="desc" rows="5" cols="30" v-model="item.desc"></textarea>
                    <span class="hint">描述将在分类meta中显示</span>
                </p>
                <p class="act">
                    <button class="formbutton" type="submit">{{item.id ? '修改':'添加'}}</button><a v-show="item.id" class="ml10" href="javascript:void(0)" @click="cancel">取消</a>
                </p>
            </form>
        </div>
    </div>
</template>

<script>
  import {cateDeleteApi, cateListApi, catePostApi} from "../service";
  import {BatchHandle} from "../components";
  import list from "../mixins/list"

  export default {
    name: "AdminCate",
    data() {
      return {
        listApi: cateListApi,
        postApi: catePostApi,
        deleteApi:cateDeleteApi
      }
    },
    mixins: [list],
    components: {
      BatchHandle
    },
    methods: {
      submit() {
        let {id, name,domain,desc} = this.item
        this.triggerSubmit({id, name,domain,desc})
      },
    },
    mounted() {
      this.loadList()
    }
  }
</script>
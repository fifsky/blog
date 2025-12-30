<template>
    <div class="clearfix">
        <h2>管理心情</h2>

        <div class="col-left">

            <div class="operate clearfix">
                <BatchHandle/>
            </div>

            <table class="list">
                <tbody>
                <tr>
                    <th width="20">&nbsp;</th>
                    <th width="80">作者</th>
                    <th>心情</th>
                    <th width="90">日期</th>
                    <th width="80">操作</th>
                </tr>
                <tr v-if="list.length === 0">
                    <td colspan="7" align="center">还没有心情！</td>
                </tr>
                <tr v-if="list.length > 0" v-for="(v,k) in list" :key="v.id">
                    <td><input type="checkbox" name="ids" :value="v.id"/></td>
                    <td>{{v.user.name}}</td>
                    <td>{{ v.content }}</td>
                    <td>{{ v.created_at | formatDate('YYYY-MM-DD') }}</td>
                    <td><a href="javascript:void(0)" @click="editItem(v.id)">编辑</a><span class="line">|</span><a
                            href="javascript:void(0)"
                            @click="deleteItem(v.id)">删除</a>
                    </td>
                </tr>
                </tbody>
            </table>
            <div class="operate clearfix">
                <BatchHandle/>
                <Paginate
                        v-model="page"
                        :page-count="pageTotal"
                        :click-handler="changePage"
                        :prev-text="'<上一页'"
                        :next-text="'下一页>'"
                        :container-class="'paginator'">
                </Paginate>
            </div>
        </div>
        <div class="col-right" style="width: 250px; padding-top: 31px;">
            <form class="vf" method="post" autocomplete="off" @submit.prevent="submit">
                <p>
                    <label class="label_input">发表心情</label>
                    <textarea name="content" rows="5" cols="30" v-model="item.content"></textarea>
                </p>
                <p class="act">
                    <button class="formbutton" type="submit">{{item.id ? '修改':'添加'}}</button>
                    <a v-show="item.id" class="ml10" href="javascript:void(0)" @click="cancel">取消</a>
                </p>
            </form>
        </div>
    </div>
</template>

<script>
  import {moodDeleteApi, moodListApi, moodCreateApi, moodUpdateApi} from "../service";
  import Paginate from 'vuejs-paginate'
  import {BatchHandle} from "../components";
  import list from "../mixins/list"

  export default {
    name: "AdminMood",
    data() {
      return {
        createApi: moodCreateApi,
        updateApi: moodUpdateApi,
        deleteApi: moodDeleteApi,
        listApi: moodListApi
      }
    },
    mixins: [list],
    components: {
      Paginate,
      BatchHandle
    },
    methods: {
      submit() {
        let {id, content} = this.item
        this.triggerSubmit({id, content})
      },
    },
    mounted() {
      this.loadList()
    }
  }
</script>

<style scoped>

</style>

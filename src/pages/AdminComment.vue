<template>
    <div>
        <h2>管理评论</h2>

        <div class="operate clearfix">
            <BatchHandle/>
        </div>

        <table class="list">
            <tbody>
            <tr>
                <th width="20">&nbsp;</th>
                <th width="150">文章</th>
                <th width="60">昵称</th>
                <th>评论</th>
                <th width="80">IP</th>
                <th width="130">日期</th>
                <th width="80">操作</th>
            </tr>
            <tr v-if="list.length === 0">
                <td colspan="7" align="center">还没有评论！</td>
            </tr>
            <tr v-if="list.length > 0" v-for="(v,k) in list" :key="v.id">
                <td><input type="checkbox" name="ids" :value="v.id"/></td>
                <td><a :href="(v.type === 2 ? v.url : '/article' + v.id)+'#comments'"
                                           target="_blank">{{ v.article_title }}</a></td>
                <td>{{v.name}}</td>
                <td>{{ v.content }}</td>
                <td>{{ v.ip }}</td>
                <td>{{ v.created_at | formatDate('YYYY-MM-DD HH:mm') }}</td>
                <td><a href="javascript:void(0)" @click="deleteItem(v.id)">删除</a></td>
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
</template>

<script>
  import {commentAdminListApi,commentDeleteApi} from "../service";
  import Paginate from 'vuejs-paginate'
  import {BatchHandle} from "../components";
  import list from "../mixins/list"

  export default {
    name: "AdminComment",
    data() {
      return {
        listApi:commentAdminListApi,
        deleteApi:commentDeleteApi,
      }
    },
    mixins: [list],
    components: {
      Paginate,
      BatchHandle
    },
    mounted() {
      this.loadList()
    }
  }
</script>

<style scoped>

</style>
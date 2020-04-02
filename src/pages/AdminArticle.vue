<template>
    <div id="articles">
        <h2>管理文章
            <router-link to="/admin/post/article" class="add"><i class="iconfont icon-edit" style="color: #444"></i>写文章
            </router-link>
        </h2>

        <div class="operate clearfix">
            <BatchHandle/>
            <div class="fr">
                <form id="list_form" class="nf" action="" method="get" autocomplete="off">
                    <select name="post_status">
                        <option value="" selected>显示所有文章</option>
                        <option value="1">页面</option>
                        <option value="2">回收站</option>
                    </select>
                    <select name="cate_id">
                        <option value="" selected>显示所有分类</option>
                    </select>
                </form>
            </div>
        </div>

        <table class="list">
            <tbody>
            <tr>
                <th width="20">&nbsp;</th>
                <th width="20"><i class="iconfont icon-comment fs-12"></i></th>
                <th>标题</th>
                <th width="60">作者</th>
                <th width="80">分类</th>
                <th width="80">类型</th>
                <th width="90">日期</th>
                <th width="80">操作</th>
            </tr>
            <tr v-if="list.length === 0">
                <td colspan="7" align="center">还没有文章，来
                    <router-link to="/admin/post/article">创建一篇</router-link>
                    文章吧！
                </td>
            </tr>
            <tr v-if="list.length > 0" v-for="(v,k) in list" :key="v.id">
                <td><input type="checkbox" name="ids" :value="v.id"/></td>
                <td class="comment-num"><a :href="(v.type === 2 ? v.url : '/article' + v.id)+'#comments'"
                                           target="_blank">{{ v.comment_num }}</a></td>
                <td><a :href="(v.type === 2 ? v.url : '/article/' + v.id)" target="_blank">{{ v.title }}</a>
                </td>
                <td>{{ v.user.nick_name }}</td>
                <td><a :href="'/category/'+v.cate.domain" target="_blank">{{ v.cate.name}}</a></td>
                <td>{{ v.type === 1 ? '文章' : '页面' }}</td>
                <td>{{ v.updated_at | formatDate('YYYY-MM-DD') }}</td>
                <td>
                    <template v-if="v.user_id === userInfo.id">
                        <router-link :to="'/admin/post/article?id='+v.id">编辑</router-link><span class="line">|</span><a href="javascript:void(0)" @click="deleteItem(v.id)">删除</a>
                    </template>
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
</template>

<script>
  import {articleDeleteApi, articleListApi} from "../service";
  import Paginate from 'vuejs-paginate'
  import {mapState} from "vuex"
  import {BatchHandle} from "../components";
  import list from "../mixins/list"

  export default {
    name: "AdminArticle",
    data() {
      return {
        listApi: articleListApi,
        deleteApi: articleDeleteApi,
      }
    },
    mixins: [list],
    components: {
      Paginate,
      BatchHandle
    },
    computed: {
      ...mapState(["userInfo"])
    },
    mounted() {
      this.loadList()
    }
  }
</script>

<style scoped>

</style>
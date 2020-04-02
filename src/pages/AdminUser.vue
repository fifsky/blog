<template>
    <div>
        <h2>管理用户
            <router-link to="/admin/post/user" class="add"><i class="iconfont icon-add" style="color: #444"></i>新增用户
            </router-link>
        </h2>

        <div class="operate clearfix">
            <BatchHandle/>
        </div>

        <table class="list">
            <tbody>
            <tr>
                <th width="20">&nbsp;</th>
                <th width="80">用户名</th>
                <th width="80">昵称</th>
                <th>邮箱</th>
                <th width="60">角色</th>
                <th width="60">状态</th>
                <th width="90">操作</th>
            </tr>
            <tr v-if="list.length === 0">
                <td colspan="7" align="center">还没有用户！</td>
            </tr>
            <tr v-if="list.length > 0" v-for="(v,k) in list" :key="v.id">
                <td><input type="checkbox" name="ids" :value="v.id"/></td>
                <td>{{ v.name }}</td>
                <td>{{ v.nick_name }}</td>
                <td>{{ v.email }}</td>
                <td>{{ v.type === 1 ? '管理员' :'编辑' }}</td>
                <td>{{ v.status === 1 ? '启用' :'停用' }}</td>
                <td><router-link :to="'/admin/post/user?id='+v.id">编辑</router-link><span class="line">|</span><a href="javascript:void(0)" @click="deleteItem(v.id)">{{ v.status === 1 ? '停用' : '启用'}}</a></td>
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
  import {userListApi, userStatusApi} from "../service";
  import Paginate from 'vuejs-paginate'
  import {BatchHandle} from "../components";
  import list from "../mixins/list"

  export default {
    name: "AdminUser",
    data() {
      return {
        listApi: userListApi,
        deleteApi: userStatusApi,
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
<template>
    <div>
        <div class="articles" v-for="(v,k) in list" :key="k">
            <CArticle :article="v"></CArticle>
            <div class="post-meta">
<!--                <span v-if="v.comment_num === 0">暂无评论</span><span-->
<!--                    v-if="v.comment_num > 0">评论({{v.comment_num}})</span><span class="line">|</span>-->
<!--                <router-link-->
<!--                        :to="'/article/'+v.id">查看更多&gt;&gt;-->
<!--                </router-link>-->
            </div>
        </div>
        <Paginate
                v-model="page"
                :page-count="pageTotal"
                :click-handler="changePage"
                :prev-text="'<上一页'"
                :next-text="'下一页>'"
                :container-class="'paginator'">
        </Paginate>
    </div>
</template>

<script>
  import {sync} from "../utils";
  import {articleListApi} from "../service";
  import {mapMutations, mapState} from "vuex"
  import {CArticle} from "../components";
  import Paginate from 'vuejs-paginate'

  export default {
    name: 'ArticleList',
    components: {
      CArticle,
      Paginate
    },
    data() {
      return {
        list: [],
        pageTotal: 0,
        page: 1
      }
    },
    methods: {
      ...mapMutations(['setKeyword']),
      loadList() {
        if (this.$route.query.page) {
          this.page = parseInt(this.$route.query.page)
        }else{
          this.page = 1
        }
        let data = {...this.$route.params, ...this.$route.query}
        data.page = this.page
        data.type = 1
        sync(async () => {
          let ret = await articleListApi(data)
          this.list = ret.list
          this.pageTotal = ret.page_total
        })
      },
      changePage(pageNum) {
        this.page = pageNum
        let q = {...this.$route.query}
        q.page = pageNum
        this.$router.push({path:this.$router.path,query:q})
      }
    },
    computed: {
      ...mapState(['keyword'])
    },
    mounted(){
      if (this.$route.path !== '/search') {
        this.setKeyword('')
      }
      this.loadList()
    }
  }
</script>

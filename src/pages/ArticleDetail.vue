<template>
    <div v-if="article.id">
        <div class="article-single">
            <CArticle :article="article"></CArticle>
            <div class="post-navi">
                <div class="prev"><strong>上一篇：</strong><span v-if="data.prev.id"><router-link
                        :to="'/article/'+data.prev.id">{{data.prev.title}}</router-link></span><span
                        v-if="!data.prev.id">嘿，这已经是最新的文章啦</span></div>
                <div class="next"><strong>下一篇：</strong><span v-if="data.next.id"><router-link
                        :to="'/article/'+data.next.id">{{data.next.title}}</router-link></span><span
                        v-if="!data.next.id">嘿，这已经是最后的文章啦</span></div>
            </div>
        </div>
<!--        <Comment :postId="article.id"></Comment>-->
    </div>
</template>

<script>
  import {CArticle, Comment} from "../components";
  import {sync} from "../utils";
  import {articleDetailApi, prevnextArticleApi} from "../service";

  export default {
    name: 'ArticleDetail',
    components: {
      CArticle,
      Comment
    },
    data() {
      return {
        article:{},
        data: {
          prev:{},
          next:{}
        },
      }
    },
    metaInfo() {
      return {
        title: this.article.title,
      }
    },
    methods: {
      load(data) {
        if(data.id){
          data.id = parseInt(data.id)
        }
        sync(async () => {
          this.article = await articleDetailApi(data)
          this.data = await prevnextArticleApi(data)
        })
      }
    },
    mounted() {
      this.load(this.$route.params)
    },
  }
</script>
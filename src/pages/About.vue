<template>
    <div v-if="article.id">
        <div class="article-single">
            <CArticle :article="article"></CArticle>
        </div>
<!--        <Comment :postId="article.id"></Comment>-->
    </div>
</template>

<script>
  import {CArticle, Comment} from "../components";
  import {sync} from "../utils";
  import {articleDetailApi} from "../service";

  export default {
    name: 'About',
    components: {
      CArticle,
      Comment
    },
    data() {
      return {
        article: {},
      }
    },
    metaInfo: {
      title: "关于我"
    },
    methods: {
      load() {
        sync(async () => {
          this.article = await articleDetailApi({url: "about"})
        })
      }
    },
    mounted() {
      this.load()
    }
  }
</script>
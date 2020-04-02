<template>
    <div class="article" v-if="article">
        <div class="entry-title">
            <img class="avatar" src="../assets/images/avatar.jpg" alt="">
            <h2>
                <router-link :to="'/article/'+article.id" v-html="markHigh(article.title,$route.query.keyword)"></router-link>
            </h2>
            <div class="entry-meta">
                by&nbsp;{{article.user.nick_name}}&nbsp;&nbsp;/&nbsp;&nbsp;<a rel="category tag"
                                                                         title="查看 默认分类 中的全部文章"
                                                                         :href="'/categroy/'+article.cate.domain">{{article.cate.name}}</a>&nbsp;&nbsp;/&nbsp;&nbsp;{{article.created_at
                | formatDate('YYYY-MM-DD HH:mm')}}
            </div>
        </div>
        <div class="entry" v-html="article.content" v-highlight>
        </div>
    </div>
</template>

<script>
    export default {
        name: "CArticle",
        props: ["article"],
        methods: {
            markHigh(content,keyword) {
                if(keyword){
                    return content.replace(keyword, '<mark>' + keyword + '</mark>')
                }else{
                    return content
                }
            }
        }
    }
</script>
<template>
    <div id="sidebar">
        <div class="sect" id="search">
            <form id="searchpanel" method="get" @submit.prevent="submit">
                <p>
                    <input class="input_text" type="text" name="keyword" @change="changeKeyword" :value="keyword"/>&nbsp;<input
                        class="formbutton" type="submit" value="搜索"/>
                </p>
            </form>
        </div>

        <Calendar/>
        <SidebarList title="文章分类" api="cateAllApi"></SidebarList>
        <SidebarList title="历史存档" api="archiveApi"></SidebarList>
<!--        <SidebarList title="最新评论" api="newCommentApi"></SidebarList>-->
        <SidebarList title="我关注的" api="linkAllApi"></SidebarList>

        <div class="rss">
            <i class="iconfont icon-feed" style="color: orangered"></i><a href="https://fifsky.com/feed.xml" target="_blank">订阅我的消息</a>
        </div>
    </div>
</template>

<script>
    import SidebarList from "./SidebarList";
    import Calendar from "./Calendar";
    import {mapState,mapMutations} from "vuex"

    export default {
        name: "Sidebar",
        components: {
            SidebarList,
            Calendar
        },
        computed:{
            ...mapState(['keyword'])
        },
        methods: {
            ...mapMutations(['setKeyword']),
            submit: function () {
                this.$router.push({name: 'search', query: {keyword: this.keyword}})
            },
            changeKeyword(e){
                this.setKeyword(e.target.value)
            }
        },
        created(){
            this.setKeyword(this.$route.query.keyword)
        }
    }
</script>
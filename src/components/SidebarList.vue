<template>
    <div class="sect">
        <h2>{{title}}</h2>
        <ul class="tlist">
            <li v-for="(v,k) in items" :key="k">
                <a v-if="v.url && v.url.substr(0, 4) === 'http'" target="_blank" :href="v.url">{{v.content}}</a>
                <router-link :to="v.url" v-else>{{v.content}}</router-link>
            </li>
        </ul>
    </div>
</template>

<script>
  import {cateAllApi,newCommentApi,archiveApi,linkAllApi} from "../service"
  import {sync} from "../utils";

  const apiMap = {
    cateAllApi,newCommentApi,archiveApi,linkAllApi
  }

  export default {
    name: "SidebarList",
    props: ["title", "api"],
    data() {
      return {
        items: {}
      }
    },
    mounted() {
      sync(async () => {
        this.items = await apiMap[this.api]()
      })
    }
  }
</script>
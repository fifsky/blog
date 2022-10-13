<template>
    <div class="comment" id="comments">
    </div>
</template>

<script>
  export default {
    name: "Comment",
    props: ["postId"],
    mounted() {
      new Artalk({
        useBackendConf: true,
        el:        '#comments',              // 绑定元素的 Selector
        server:    'https://comment-api.fifsky.com',  // 后端地址
        site:      'FIFSKY',           // 你的站点名
      })
    },
  }
</script>

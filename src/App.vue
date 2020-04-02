<template>
  <div id="app">
    <router-view/>
  </div>
</template>
<script>
export default {
  name: 'app',
  metaInfo: {
    // all titles will be injected into this template
    titleTemplate: (titleChunk) => {
      // If undefined or blank then we don't need the hyphen
      return titleChunk ? `${titleChunk} - 無處告別` : '無處告別';
    },
    meta: [
      { name: 'keywords', content: 'fifsky,rita,生活,博客,豆豆' },
      { name: 'description', content: '生活在别处' }
    ]
  }
}
</script>
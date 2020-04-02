<template>
    <div class="fl">
        <a href="javascript:void(0);" class="all-selected">全选</a><span class="line">|</span><a
            href="javascript:void(0);" class="inverse-selected">反选</a>&nbsp;
        <select name="batch_operation">
            <option value="" selected>批量操作</option>
            <option value="1">置顶</option>
            <option value="1">删除</option>
        </select>
    </div>
</template>

<script>
  export default {
    name: "BatchHandle"
  }
</script>

<style scoped>

</style>
<template>
    <div id="info">
        <div id="avatar"><img title="莫一哲" alt="莫一哲" src="../assets/images/faceicon.jpg"></div>
        <div id="latest">
            <p class="current active" v-if="moods[index]">
                {{moods[index].content}}
                <span class="stamp">
                <span class="method">{{moods[index].created_at | humanTime }} by {{moods[index].user.nick_name}}</span>
            </span>
            </p>
        </div>
        <div class="handle">
            <i class="iconfont icon-left" title="上一条" @click="prev"></i>
            <i class="iconfont icon-right" title="下一条" @click="next"></i>
        </div>
    </div>
</template>

<script>
  import {sync} from "../utils";
  import {moodListApi} from "../service";

  export default {
    name: 'Mood',
    data() {
      return {
        moods: [],
        index: 0
      }
    },
    methods: {
      prev() {
        let i = this.index - 1
        if (i >= 0) {
          this.index = i
        }
      },
      next() {
        let i = this.index + 1
        if (i < this.moods.length) {
          this.index = i
        }
      }
    },
    mounted() {
      sync(async () => {
        let ret = await moodListApi({page:1})
        this.moods = ret.list
        if (this.moods.length > 0) {
          this.current = this.moods[0]
        }
      })
    }
  }
</script>
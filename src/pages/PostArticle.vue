<template>
  <div id="articles">
    <h2>
      {{ article.id ? "编辑" : "撰写" }}文章
      <router-link to="/admin/articles"
        ><i class="iconfont icon-undo" style="color: #444"></i
        >返回列表</router-link
      >
    </h2>
    <div class="message"></div>
    <form class="vf" method="post" autocomplete="off" @submit.prevent="submit">
      <div class="clearfix">
        <div class="col-left">
          <p>
            <label class="label_input">标题</label>
            <input
              type="text"
              class="input_text"
              maxlength="200"
              size="50"
              name="title"
              v-model="article.title"
            />
          </p>
          <p>
            <label class="label_input">分类</label>

            <select
              name="cate_id"
              v-model="article.cate_id"
              v-if="cates.length"
            >
              <option v-for="(v, k) in cates" :key="v.id" :value="v.id">{{
                v.name
              }}</option>
            </select>
          </p>
          <p v-show="article.type == 2">
            <label class="label_input">缩略名</label>
            <input
              type="text"
              class="input_text"
              maxlength="200"
              size="50"
              name="url"
              v-model="article.url"
            />
            <span class="hint"
              >页面的URL名称，如红色部分http://domain.com/<span
                style="color: red;"
                >about</span
              ></span
            >
          </p>
        </div>
        <div class="col-right">
          <p>
            <label class="label_input">类型</label>
            <input
              class="input_check"
              name="type"
              type="radio"
              v-model="article.type"
              value="1"
            />
            文章
            <input
              class="input_check"
              name="type"
              type="radio"
              v-model="article.type"
              value="2"
            />
            页面
            <br />
            <!--<select name="power">-->
            <!--<option value="1">公开</option>-->
            <!--<option value="2">私密</option>-->
            <!--</select>-->
            <!--<input style="display: none" type="text" class="input_text" name="post_password"-->
            <!--value="" size="20"/>-->
          </p>
        </div>
      </div>
      <div id="editor"></div>
      <p class="act">
        <input class="formbutton" type="submit" value="发布" /><a
          id="_save_draft"
          href="javascript:void(0)"
          class="ml10"
          >保存草稿</a
        >
      </p>
    </form>
  </div>
</template>

<script>
import { getAccessToken, getApiUrl, sync } from "../utils";
import {
  articleDetailApi,
  articleCreateApi,
  articleUpdateApi,
  cateListApi,
} from "../service";
import { mapState } from "vuex";
import "highlight.js/styles/atelier-lakeside-light.css";
import WangEditor from "wangeditor";

export default {
  name: "PostArticle",
  data() {
    return {
      article: {
        type: 1,
      },
      cates: [],
      editor: "", //保存simditor对象
    };
  },
  computed: {
    ...mapState(["userInfo"]),
  },
  methods: {
    submit() {
      let { id, cate_id, title, content, type, url } = this.article;

      if (title.length === 0) {
        this.$message.success("标题不能为空");
      }

      if (cate_id < 1) {
        this.$message.success("请选择分类");
      }

      let data = { id, cate_id, title, content, type, url };
      sync(async () => {
        let ret = await (id ? articleUpdateApi : articleCreateApi)(data);
        this.$router.push("/admin/articles");
      });
    },
    createEditor() {
      const editor = new WangEditor("#editor");
      editor.config.uploadImgMaxSize = 3 * 1024 * 1024;
      editor.config.uploadImgMaxLength = 5;
      editor.config.uploadImgServer = "https://api.fifsky.com/api/admin/upload";
      editor.config.uploadFileName = "uploadFile";
      editor.config.uploadImgHeaders = {
        "Access-Token": getAccessToken(),
      };
      editor.config.uploadImgHooks = {
        error: function(xhr, editor) {
          // 图片上传出错时触发
          // xhr 是 XMLHttpRequst 对象，editor 是编辑器对象
          console.log(xhr);
        },
      };
      editor.config.onchange = (newHtml) => {
        this.article.content = newHtml;
      };
      editor.create();
      this.editor = editor;
    },
  },
  mounted() {
    this.createEditor();
    sync(async () => {
      if (this.$route.query && this.$route.query.id) {
        this.article = await articleDetailApi({
          id: parseInt(this.$route.query.id),
        });
        this.editor.txt.html(this.article.content);
      }
      let ret = await cateListApi();
      this.cates = ret.list;
      this.article.cate_id = this.article.cate_id || this.cates[0].id;
    });
  },
  beforeDestroy() {
    // 销毁编辑器
    this.editor.destroy();
    this.editor = null;
  },
};
</script>

<style>
#editor pre code {
  background-color: #ebf8ff;
  color: #516d7b;
  overflow: visible;
}

#editor .w-e-toolbar .w-e-menu {
  height: 38px !important;
  width: 38px !important;
}
#editor .w-e-text p {
  font-size: 14px !important;
}
</style>

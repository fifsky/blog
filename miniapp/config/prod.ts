import type { UserConfigExport } from "@tarojs/cli";
export default {
  // 关键点：用 env 显式设置 NODE_ENV，保证在 config/index.ts、业务代码里判断环境时一致
  env: {
    NODE_ENV: JSON.stringify("production"),
  },
  defineConstants: {
    // 关键点：defineConstants 会直接替换代码中的标识符，这里必须是合法 JS 字面量（用 JSON.stringify 最稳）
    API_BASE_URL: JSON.stringify("https://api.fifsky.com"),
  },
  mini: {},
  h5: {
    /**
     * WebpackChain 插件配置
     * @docs https://github.com/neutrinojs/webpack-chain
     */
    // webpackChain (chain) {
    //   /**
    //    * 如果 h5 端编译后体积过大，可以使用 webpack-bundle-analyzer 插件对打包体积进行分析。
    //    * @docs https://github.com/webpack-contrib/webpack-bundle-analyzer
    //    */
    //   chain.plugin('analyzer')
    //     .use(require('webpack-bundle-analyzer').BundleAnalyzerPlugin, [])
    //   /**
    //    * 如果 h5 端首屏加载时间过长，可以使用 prerender-spa-plugin 插件预加载首页。
    //    * @docs https://github.com/chrisvfritz/prerender-spa-plugin
    //    */
    //   const path = require('path')
    //   const Prerender = require('prerender-spa-plugin')
    //   const staticDir = path.join(__dirname, '..', 'dist')
    //   chain
    //     .plugin('prerender')
    //     .use(new Prerender({
    //       staticDir,
    //       routes: [ '/pages/index/index' ],
    //       postProcess: (context) => ({ ...context, outputPath: path.join(staticDir, 'index.html') })
    //     }))
    // }
  },
} satisfies UserConfigExport<"vite">;

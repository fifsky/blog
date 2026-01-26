import { defineConfig, type UserConfigExport } from "@tarojs/cli";
import type { Plugin } from "vite";
import devConfig from "./dev";
import prodConfig from "./prod";

// https://taro-docs.jd.com/docs/next/config#defineconfig-辅助函数
export default defineConfig<"vite">(async (merge, { command, mode }) => {
  void command;
  const baseConfig: UserConfigExport<"vite"> = {
    projectName: "miniapp",
    date: "2026-1-26",
    designWidth: 750,
    deviceRatio: {
      640: 2.34 / 2,
      750: 1,
      375: 2,
      828: 1.81 / 2,
    },
    sourceRoot: "src",
    outputRoot: "dist",
    plugins: [],
    defineConstants: {},
    copy: {
      patterns: [],
      options: {},
    },
    framework: "react",
    compiler: {
      type: "vite",
      vitePlugins: [] as Plugin[],
    },
    mini: {
      postcss: {
        pxtransform: {
          enable: true,
          config: {},
        },
        cssModules: {
          enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
          config: {
            namingPattern: "module", // 转换模式，取值为 global/module
            generateScopedName: "[name]__[local]___[hash:base64:5]",
          },
        },
      },
    },
    h5: {
      publicPath: "/",
      staticDirectory: "static",

      miniCssExtractPluginOption: {
        ignoreOrder: true,
        filename: "css/[name].[hash].css",
        chunkFilename: "css/[name].[chunkhash].css",
      },
      postcss: {
        autoprefixer: {
          enable: true,
          config: {},
        },
        cssModules: {
          enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
          config: {
            namingPattern: "module", // 转换模式，取值为 global/module
            generateScopedName: "[name]__[local]___[hash:base64:5]",
          },
        },
      },
    },
    rn: {
      appName: "taroDemo",
      postcss: {
        cssModules: {
          enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
        },
      },
    },
  };

  // 关键点：Taro 官方推荐用 mode 区分环境（development / production）。
  // - mode 来自 taro cli 的 --mode 参数；--watch 默认等价于 --mode development
  // - 注意：在“选择 dev/prod 配置”这个阶段，NODE_ENV 可能还没被 dev/prod.ts 的 env 注入
  //   所以这里不能用 NODE_ENV 作为主判定，否则会出现本地 dev 仍走 prod 的情况
  const hasWatchFlag = process.argv.includes("--watch") || process.argv.includes("-w");

  const shouldUseDevConfig =
    mode === "development" || hasWatchFlag || process.env.NODE_ENV === "development";

  if (shouldUseDevConfig) {
    return merge({}, baseConfig, devConfig);
  }

  // 生产构建配置（默认开启压缩混淆等，API 指向线上等）
  return merge({}, baseConfig, prodConfig);
});

import type { UserConfigExport } from "@tarojs/cli";
export default {
  // 关键点：用 env 显式设置 NODE_ENV，保证在 config/index.ts、业务代码里判断环境时一致
  env: {
    NODE_ENV: JSON.stringify("development"),
  },
  defineConstants: {
    // 关键点：defineConstants 会直接替换代码中的标识符，这里必须是合法 JS 字面量（用 JSON.stringify 最稳）
    API_BASE_URL: JSON.stringify("http://192.168.1.85:8080"),
  },
  mini: {},
  h5: {},
} satisfies UserConfigExport<"vite">;

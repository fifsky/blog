import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { resolve } from "path";

export default defineConfig(() => {
  const isProd = process.env.NODE_ENV === "production";
  return {
    base: isProd ? "https://static.fifsky.com/" : "/",
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": resolve(__dirname, "src"),
      },
    },
    server: {
      port: 5173,
      proxy: {
        "/api": {
          target: "http://127.0.0.1:8080",
          changeOrigin: true,
        },
        "/feed.xml": {
          target: "http://127.0.0.1:8080",
          changeOrigin: true,
        },
      },
    },
    build: {
      outDir: "dist",
      assetsDir: "assets",
      rollupOptions: {
        output: {
          manualChunks: {
            "utils-vendor": ["dayjs", "prismjs", "@wangeditor/editor"],
          },
        },
      },
    },
  };
});

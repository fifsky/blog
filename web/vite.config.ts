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
        "/blog": {
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
            "vendor-react": ["react", "react-dom", "react-router"],
            "vendor-ui": ["motion", "lucide-react"],
            "vendor-radix": [
              "@radix-ui/react-dialog",
              "@radix-ui/react-select",
              "@radix-ui/react-checkbox",
              "@radix-ui/react-label",
              "@radix-ui/react-radio-group",
              "@radix-ui/react-separator",
              "@radix-ui/react-slot",
              "@radix-ui/react-switch",
              "@radix-ui/react-tooltip",
              "@radix-ui/react-alert-dialog"
            ],
            "vendor-editor": [
              "bytemd",
              "@bytemd/react",
              "@bytemd/plugin-breaks",
              "@bytemd/plugin-gfm",
              "@bytemd/plugin-medium-zoom",
              "highlight.js",
              "rehype-highlight"
            ],
            "vendor-form": ["react-hook-form", "zod", "@hookform/resolvers"],
            "vendor-utils": ["dayjs", "zustand", "sonner", "clsx", "tailwind-merge"],
          },
        },
      },
    },
  };
});

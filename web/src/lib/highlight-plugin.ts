// Custom highlight plugin with limited languages for smaller bundle size
// Only includes: js, go, php, shell, css, html, python, java, yml, json
import type { BytemdPlugin } from "bytemd";
import rehypeHighlight from "rehype-highlight";
import hljs from "highlight.js/lib/core";

// Import only the languages we need
import javascript from "highlight.js/lib/languages/javascript";
import typescript from "highlight.js/lib/languages/typescript";
import go from "highlight.js/lib/languages/go";
import php from "highlight.js/lib/languages/php";
import bash from "highlight.js/lib/languages/bash";
import css from "highlight.js/lib/languages/css";
import xml from "highlight.js/lib/languages/xml"; // includes html
import python from "highlight.js/lib/languages/python";
import java from "highlight.js/lib/languages/java";
import yaml from "highlight.js/lib/languages/yaml";
import json from "highlight.js/lib/languages/json";
import sql from "highlight.js/lib/languages/sql";
import markdown from "highlight.js/lib/languages/markdown";

// Register languages
hljs.registerLanguage("javascript", javascript);
hljs.registerLanguage("js", javascript);
hljs.registerLanguage("typescript", typescript);
hljs.registerLanguage("ts", typescript);
hljs.registerLanguage("go", go);
hljs.registerLanguage("golang", go);
hljs.registerLanguage("php", php);
hljs.registerLanguage("bash", bash);
hljs.registerLanguage("shell", bash);
hljs.registerLanguage("sh", bash);
hljs.registerLanguage("css", css);
hljs.registerLanguage("xml", xml);
hljs.registerLanguage("html", xml);
hljs.registerLanguage("python", python);
hljs.registerLanguage("py", python);
hljs.registerLanguage("java", java);
hljs.registerLanguage("yaml", yaml);
hljs.registerLanguage("yml", yaml);
hljs.registerLanguage("json", json);
hljs.registerLanguage("sql", sql);
hljs.registerLanguage("markdown", markdown);
hljs.registerLanguage("md", markdown);

export function highlightPlugin(): BytemdPlugin {
  return {
    rehype: (processor) =>
      // Use rehype-highlight with our custom subset hljs instance
      processor.use(rehypeHighlight, {
        // Provide our custom hljs instance with only the languages we registered
        subset: [
          "javascript",
          "js",
          "typescript",
          "ts",
          "go",
          "golang",
          "php",
          "bash",
          "shell",
          "sh",
          "css",
          "xml",
          "html",
          "python",
          "py",
          "java",
          "yaml",
          "yml",
          "json",
          "sql",
          "markdown",
          "md",
        ],
        ignoreMissing: true,
      }),
  };
}

export default highlightPlugin;

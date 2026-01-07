import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import hljs from "highlight.js/lib/core";
import javascript from "highlight.js/lib/languages/javascript";
import typescript from "highlight.js/lib/languages/typescript";
import php from "highlight.js/lib/languages/php";
import go from "highlight.js/lib/languages/go";
import python from "highlight.js/lib/languages/python";
import nginx from "highlight.js/lib/languages/nginx";
import sql from "highlight.js/lib/languages/sql";
import lua from "highlight.js/lib/languages/lua";
import bash from "highlight.js/lib/languages/bash";
import css from "highlight.js/lib/languages/css";
import java from "highlight.js/lib/languages/java";
import xml from "highlight.js/lib/languages/xml";
import "highlight.js/styles/atom-one-light.css";
import "nprogress/nprogress.css";
import "./index.css";

dayjs.extend(isBetween);

hljs.registerLanguage("javascript", javascript);
hljs.registerLanguage("typescript", typescript);
hljs.registerLanguage("php", php);
hljs.registerLanguage("go", go);
hljs.registerLanguage("python", python);
hljs.registerLanguage("nginx", nginx);
hljs.registerLanguage("sql", sql);
hljs.registerLanguage("lua", lua);
hljs.registerLanguage("bash", bash);
hljs.registerLanguage("css", css);
hljs.registerLanguage("java", java);
hljs.registerLanguage("xml", xml);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

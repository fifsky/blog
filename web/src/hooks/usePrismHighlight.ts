import { useEffect } from "react";
import Prism from "prismjs";
import "@/assets/prismjs.css";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import "prismjs/components/prism-typescript";
import "prismjs/components/prism-markup-templating";
import "prismjs/components/prism-php";
import "prismjs/components/prism-go";
import "prismjs/components/prism-python";
import "prismjs/components/prism-nginx";
import "prismjs/components/prism-sql";
import "prismjs/components/prism-lua";
import "prismjs/components/prism-bash";
import "prismjs/components/prism-css";
import "prismjs/components/prism-java";
import "prismjs/components/prism-markup";
import "prismjs/plugins/line-numbers/prism-line-numbers.js";
import "prismjs/plugins/line-numbers/prism-line-numbers.css";

export function usePrismHighlight(dependencies: any[] = []) {
  useEffect(() => {
    if (!dependencies) {
      return;
    }
    const highlightCode = () => {
      document.querySelectorAll("pre code").forEach((block) => {
        if (!block.classList.contains("token") && !(block as HTMLElement).dataset.highlighted) {
          (block as HTMLElement).dataset.highlighted = "true";
          const pre = block.parentElement;
          pre?.classList.add("line-numbers");
          Prism.highlightElement(block as HTMLElement);
        }
      });
    };
    setTimeout(highlightCode, 200);
  }, dependencies);
}

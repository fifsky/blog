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
import { useEffect, useRef } from "react";

const CodeHighlight = ({ htmlContent }: { htmlContent: string }) => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    container.innerHTML = htmlContent;

    // 3. 获取所有代码块并高亮
    const codeBlocks = container.querySelectorAll("pre code");

    codeBlocks.forEach((block) => {
      const element = block as HTMLElement;
      if (!element.dataset.highlighted) {
        Prism.highlightElement(element);
        element.dataset.highlighted = "true";
      }
    });
  }, [htmlContent]);

  return <div className="article" ref={containerRef} />;
};

export default CodeHighlight;

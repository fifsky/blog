import { useEffect, useState } from "react";
import { CArticle } from "@/components/CArticle";
import { Comment } from "@/components/Comment";
import { articleDetailApi } from "@/service";

export default function About() {
  const [article, setArticle] = useState<any>({});
  useEffect(() => {
    (async () => {
      const a = await articleDetailApi({ url: "about" });
      setArticle(a);
      document.title = "关于我 - 無處告別";
    })();
  }, []);
  if (!article.id) return null;
  return (
    <div>
      <div className="mb-[10px]">
        <CArticle article={article} />
      </div>
      <Comment />
    </div>
  );
}

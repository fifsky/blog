import { useEffect, useState } from "react";
import { CArticle } from "@/components/CArticle";
import { Comment } from "@/components/Comment";
import { articleDetailApi } from "@/service";
import { ArticleItem } from "@/types/openapi";

export default function About() {
  const [article, setArticle] = useState<ArticleItem>();
  useEffect(() => {
    (async () => {
      const a = await articleDetailApi({ url: "about" });
      setArticle(a);
    })();
  }, []);
  if (!article?.id) return null;
  return (
    <>
      <title>关于我 - 無處告別</title>
      <div>
        <div className="mb-[10px]">
          <CArticle article={article} />
        </div>
        <Comment />
      </div>
    </>
  );
}

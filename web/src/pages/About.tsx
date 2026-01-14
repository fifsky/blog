import { useEffect, useState } from "react";
import { CArticle } from "@/components/CArticle";
import { Comment } from "@/components/Comment";
import { articleDetailApi, settingApi } from "@/service";
import { ArticleItem, Options } from "@/types/openapi";

export default function About() {
  const [article, setArticle] = useState<ArticleItem>();
  const [settings, setSettings] = useState<Options>();
  useEffect(() => {
    (async () => {
      const [a, s] = await Promise.all([articleDetailApi({ url: "about" }), settingApi()]);
      setArticle(a);
      setSettings(s);
    })();
  }, []);
  if (!article?.id) return null;
  const siteName = settings?.kv?.site_name || "無處告別";
  const pageTitle = `关于我 - ${siteName}`;
  return (
    <>
      <title>{pageTitle}</title>
      <meta name="description" content={settings?.kv?.site_desc || ""} />
      <meta name="keywords" content={settings?.kv?.site_keyword || ""} />
      <div>
        <div className="mb-[10px]">
          <CArticle article={article} />
        </div>
        <Comment />
      </div>
    </>
  );
}

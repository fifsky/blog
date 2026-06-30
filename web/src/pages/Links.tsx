import { useEffect, useState } from "react";
import { Link2 } from "lucide-react";
import { Empty } from "@/components/Empty";
import { PageTransition } from "@/components/PageTransition";
import { linkAllApi } from "@/service";
import type { LinkMenuItem } from "@/types/openapi";

// 从 URL 中提取域名，用于获取 favicon
function getDomain(url: string): string {
  try {
    // 补齐协议前缀，兼容不带 http(s) 的 URL
    const normalized = /^https?:\/\//i.test(url) ? url : `https://${url}`;
    return new URL(normalized).hostname;
  } catch {
    return url;
  }
}

// 构建 favicon 地址
function getFaviconUrl(url: string): string {
  const domain = getDomain(url);
  return `https://icon.bqb.cool/?url=${domain}`;
}

// 首字母占位（favicon 加载失败时回退）
function getInitial(name: string): string {
  return (name || "?").charAt(0).toUpperCase();
}

export default function Links() {
  const [list, setList] = useState<LinkMenuItem[]>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const ret = await linkAllApi();
        if (!cancelled) {
          setList(ret.list || []);
        }
      } catch {
        if (!cancelled) setList([]);
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  if (!list) return null;

  return (
    <div>
      <title>友情链接 - 無處告別</title>

      {/* 页面标题 */}
      <div className="mb-6 pb-4 border-b border-[#e5e7eb]">
        <h2 className="text-xl font-bold text-[#1f2937] flex items-center gap-2">
          <Link2 className="text-[#0066cc]" size={22} />
          友情链接
        </h2>
        <p className="mt-1.5 text-sm text-[#6b7280]">
          这里有双向的链接，也有我单方面收藏的优秀网站，如果你不希望出现在这里请给我留言。
        </p>
      </div>

      {list.length === 0 ? (
        <Empty icon={<Link2 size={24} />} title="暂无友链" content="当前还没有可展示的友情链接" />
      ) : (
        <PageTransition loading={loading}>
          <div className="grid grid-cols-3 gap-3">
            {list.map((item, idx) => {
              const domain = getDomain(item.url);
              const faviconUrl = getFaviconUrl(item.url);
              return (
                <a
                  key={`${domain}-${idx}`}
                  href={item.url}
                  target="_blank"
                  rel="noreferrer"
                  className="group flex items-center gap-3 p-3 rounded-lg border border-[#89d5ef] bg-white transition-all duration-200 hover:border-[#89d5ef] hover:-translate-y-0.5 hover:shadow-[0_4px_12px_rgba(0,102,204,0.12)]"
                >
                  {/* favicon logo */}
                  <div className="shrink-0 w-10 h-10 rounded-md overflow-hidden flex items-center justify-center border border-[#e5e7eb]">
                    <img
                      src={faviconUrl}
                      alt={item.content}
                      className="w-6 h-6 object-contain"
                      onError={(e) => {
                        // 加载失败回退为首字母占位
                        const img = e.currentTarget;
                        const fallback = document.createElement("div");
                        fallback.className =
                          "w-full h-full flex items-center justify-center text-sm font-bold text-white bg-[#89d5ef]";
                        fallback.textContent = getInitial(item.content);
                        img.parentElement?.replaceChild(fallback, img);
                      }}
                    />
                  </div>
                  {/* 友链信息 */}
                  <div className="min-w-0 flex-1">
                    <div className="text-sm font-bold text-[#1f2937] truncate group-hover:text-[#0066cc] transition-colors">
                      {item.content}
                    </div>
                    <div className="text-xs text-[#9ca3af] truncate">{domain}</div>
                  </div>
                </a>
              );
            })}
          </div>
        </PageTransition>
      )}
    </div>
  );
}

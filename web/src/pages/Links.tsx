import { useEffect, useState } from "react";
import { Link2, Send } from "lucide-react";
import { toast } from "sonner";
import { Empty } from "@/components/Empty";
import { PageTransition } from "@/components/PageTransition";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { linkAllApi, linkSubmitApi } from "@/service";
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

  // 提交表单状态
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ name: "", url: "", desc: "" });

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

  const handleSubmit = async () => {
    if (!form.name.trim() || !form.url.trim()) {
      toast.error("请填写链接名称和链接地址");
      return;
    }
    setSubmitting(true);
    try {
      await linkSubmitApi({
        name: form.name.trim(),
        url: form.url.trim(),
        desc: form.desc.trim(),
      });
      toast.success("提交成功，博主会尽快审核～～");
      setForm({ name: "", url: "", desc: "" });
    } catch {
      toast.error("提交失败，请稍后重试");
    } finally {
      setSubmitting(false);
    }
  };

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

      {/* 提交友情链接表单 */}
      <div className="mt-10 pt-6 border-t border-[#e5e7eb]">
        <h3 className="text-base font-bold text-[#1f2937] flex items-center gap-2 mb-4">
          <Send size={16} className="text-[#0066cc]" />
          申请友链
        </h3>
        <div className="max-w-lg space-y-3">
          <div>
            <Input
              placeholder="链接名称 *"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              disabled={submitting}
            />
          </div>
          <div>
            <Input
              placeholder="链接地址 * (如 https://example.com)"
              value={form.url}
              onChange={(e) => setForm({ ...form, url: e.target.value })}
              disabled={submitting}
            />
          </div>
          <div>
            <Textarea
              placeholder="链接描述（选填）"
              value={form.desc}
              onChange={(e) => setForm({ ...form, desc: e.target.value })}
              rows={3}
              disabled={submitting}
            />
          </div>
          <Button onClick={handleSubmit} loading={submitting} size="sm">
            提交申请
          </Button>
        </div>
      </div>
    </div>
  );
}

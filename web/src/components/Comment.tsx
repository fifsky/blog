import { useEffect, useState } from "react";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import "dayjs/locale/zh-cn";
import { MessageSquare } from "lucide-react";
import { Empty } from "@/components/Empty";
import { commentListApi } from "@/service";
import type { CommentItem } from "@/types/openapi";
import { CommentForm } from "@/components/CommentForm";
import { useStore } from "@/store/context";

// 评论时间使用相对时间展示（如"3小时前"）
dayjs.extend(relativeTime);
dayjs.locale("zh-cn");

// 渲染过的评论树结构：主评论 + 其下的回复列表（平铺）
type ReplyNode = CommentItem;
type RootNode = CommentItem & { replies: ReplyNode[] };

export function Comment({ postId }: { postId: number }) {
  const [roots, setRoots] = useState<RootNode[]>([]);
  const [loading, setLoading] = useState(true);
  // 当前展开的内联回复框目标：commentId 表示在哪条评论下渲染回复框
  const [replyTargetId, setReplyTargetId] = useState<number | null>(null);
  const userInfo = useStore((s) => s.userInfo);
  // 管理员登录时使用管理员身份评论：自动填充昵称/邮箱/网址（仅 host）
  const adminInfo =
    userInfo.id && userInfo.name && userInfo.email
      ? {
          name: userInfo.name,
          email: userInfo.email,
          host: typeof window !== "undefined" ? window.location.origin : "",
        }
      : undefined;

  // 加载评论列表并组装两级结构
  const loadComments = async () => {
    setLoading(true);
    try {
      const ret = await commentListApi({ post_id: postId });
      setRoots(buildTree(ret.list || []));
    } catch {
      // 错误由 request 拦截器处理
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadComments();
  }, [postId]);

  // 将平铺的评论按 pid 分组为主评论 + 回复列表
  const buildTree = (list: CommentItem[]): RootNode[] => {
    const rootMap = new Map<number, RootNode>();
    const replies: CommentItem[] = [];
    // 先收集主评论，保留原始顺序
    list.forEach((item) => {
      if (item.pid === 0) {
        rootMap.set(item.id, { ...item, replies: [] });
      } else {
        replies.push(item);
      }
    });
    // 回复挂到对应主评论下，按时间正序
    replies
      .sort((a, b) => a.created_at.localeCompare(b.created_at))
      .forEach((r) => {
        const root = rootMap.get(r.pid);
        if (root) {
          root.replies.push(r);
        }
      });
    // 主评论按时间倒序（最新的在最前）
    return Array.from(rootMap.values()).sort((a, b) =>
      b.created_at.localeCompare(a.created_at),
    );
  };

  // 昵称渲染：有网址则渲染为链接
  const renderName = (item: CommentItem) => {
    if (item.website) {
      return (
        <a
          href={item.website}
          target="_blank"
          rel="noreferrer"
          className="text-[#0066cc] hover:underline"
        >
          {item.name}
        </a>
      );
    }
    return <span className="font-medium text-[#1f2937]">{item.name}</span>;
  };

  const totalComments = roots.reduce((sum, r) => sum + 1 + r.replies.length, 0);

  return (
    <div className="mt-8">
      <div className="mb-5 flex items-center gap-2">
        <MessageSquare className="text-[#0066cc]" size={20} />
        <h2 className="text-base font-bold text-[#1f2937]">
          评论 {totalComments > 0 && <span className="text-[#9ca3af]">({totalComments})</span>}
        </h2>
      </div>

      {/* 顶部主评论框 */}
      <div className="mb-8">
        <CommentForm
          postId={postId}
          adminInfo={adminInfo}
          onSubmit={loadComments}
          submitText="发表评论"
        />
      </div>

      {/* 评论列表 */}
      {loading ? (
        <div className="py-10 text-center text-sm text-[#9ca3af]">加载中...</div>
      ) : roots.length === 0 ? (
        <Empty
          icon={<MessageSquare size={24} />}
          title="还没有评论"
          content="快来抢占第一条评论吧"
        />
      ) : (
        <div className="space-y-5">
          {roots.map((root) => (
            <div
              key={root.id}
              className="flex gap-3 border-b border-dashed border-gray-300 pb-5 last:border-b-0 last:pb-0"
            >
              {/* 头像 */}
              <img
                src={root.avatar}
                alt={root.name}
                className="w-10 h-10 rounded-full border border-[#e5e7eb] shrink-0 object-cover"
                loading="lazy"
              />
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 text-sm">
                  {renderName(root)}
                  <span className="text-xs text-[#9ca3af]">
                    {dayjs(root.created_at).fromNow()}
                  </span>
                </div>
                <div
                  className="mt-1 text-sm text-[#374151] break-words whitespace-pre-wrap"
                  dangerouslySetInnerHTML={{ __html: root.content }}
                />
                <button
                  type="button"
                  onClick={() => setReplyTargetId(replyTargetId === root.id ? null : root.id)}
                  className="mt-1 text-xs text-[#0066cc] hover:underline"
                >
                  回复
                </button>

                {/* 主评论下方的内联回复框 */}
                {replyTargetId === root.id && (
                  <div className="mt-3">
                    <CommentForm
                      postId={postId}
                      adminInfo={adminInfo}
                      pid={root.id}
                      replyName=""
                      submitText="确认回复"
                      onCancel={() => setReplyTargetId(null)}
                      onSubmit={async () => {
                        setReplyTargetId(null);
                        await loadComments();
                      }}
                    />
                  </div>
                )}

                {/* 该主评论下的回复列表（平铺） */}
                {root.replies.length > 0 && (
                  <div className="mt-3 space-y-3 pl-3 border-l-2 border-[#eef6fb]">
                    {root.replies.map((reply) => (
                      <div key={reply.id} className="flex gap-2">
                        <img
                          src={reply.avatar}
                          alt={reply.name}
                          className="w-8 h-8 rounded-full border border-[#e5e7eb] shrink-0 object-cover"
                          loading="lazy"
                        />
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 text-sm flex-wrap">
                            {renderName(reply)}
                            {reply.reply_name && (
                              <span className="text-xs text-[#9ca3af]">
                                回复 <span className="text-[#0066cc]">@{reply.reply_name}</span>
                              </span>
                            )}
                            <span className="text-xs text-[#9ca3af]">
                              {dayjs(reply.created_at).fromNow()}
                            </span>
                          </div>
                          <div
                            className="mt-1 text-sm text-[#374151] break-words whitespace-pre-wrap"
                            dangerouslySetInnerHTML={{ __html: reply.content }}
                          />
                          <button
                            type="button"
                            onClick={() =>
                              setReplyTargetId(replyTargetId === reply.id ? null : reply.id)
                            }
                            className="mt-1 text-xs text-[#0066cc] hover:underline"
                          >
                            回复
                          </button>

                          {/* 回复下方的内联回复框（pid 仍指向顶层主评论，replyName 填被回复人） */}
                          {replyTargetId === reply.id && (
                            <div className="mt-3">
                              <CommentForm
                                postId={postId}
                                adminInfo={adminInfo}
                                pid={root.id}
                                replyName={reply.name}
                                submitText="确认回复"
                                onCancel={() => setReplyTargetId(null)}
                                onSubmit={async () => {
                                  setReplyTargetId(null);
                                  await loadComments();
                                }}
                              />
                            </div>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

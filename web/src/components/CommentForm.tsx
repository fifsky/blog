import { useRef, useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, Smile } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Field, FieldLabel, FieldError, FieldGroup } from "@/components/ui/field";
import { commentCreateApi } from "@/service";

// 常用 emoji 表情集合
const EMOJI_LIST = [
  "😀", "😄", "😁", "😆", "😅", "😂", "🤣", "😊",
  "😇", "🙂", "😉", "😌", "😍", "🥰", "😘", "😋",
  "😜", "🤪", "🤨", "🧐", "🤓", "😎", "🥳", "😢",
  "😭", "😤", "😠", "🤬", "🙄", "😏", "😱", "🤯",
  "👍", "👎", "👏", "🙌", "🙏", "💪", "🤝", "✌️",
  "❤️", "💔", "🔥", "✨", "🎉", "💯", "🌹", "☕",
];

// 表单校验：内容必填，邮箱选填但需合法格式，网址选填
// 昵称在访客模式下必填，由 handleSubmit 手动校验（管理员模式跳过）
const formSchema = z.object({
  name: z.string().max(50, "昵称最多50个字符").optional(),
  email: z
    .string()
    .max(255, "邮箱过长")
    .optional()
    .refine((v) => !v || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v), "邮箱格式不正确"),
  website: z
    .string()
    .max(255, "网址过长")
    .optional()
    .refine((v) => !v || /^https?:\/\/.+/i.test(v), "网址需以 http(s):// 开头"),
  content: z.string().min(1, "请输入评论内容").max(1000, "评论内容最多1000个字符"),
});

type FormValues = z.infer<typeof formSchema>;

interface CommentFormProps {
  postId: number;
  /** 顶层主评论ID，主评论为0 */
  pid?: number;
  /** 被回复人昵称，回复的回复时填入 */
  replyName?: string;
  /** 提交按钮文案 */
  submitText: string;
  /** 提交成功后的回调（通常用于刷新评论列表） */
  onSubmit: () => void | Promise<void>;
  /** 取消回调（仅内联回复框显示取消按钮时使用） */
  onCancel?: () => void;
  /**
   * 管理员信息：传入时隐藏昵称/邮箱/网址输入框，并使用管理员身份提交。
   * - name 作为评论昵称
   * - email 作为邮箱（用于头像）
   * - host 作为网址（仅 host 部分）
   */
  adminInfo?: { name: string; email: string; host: string };
}

export function CommentForm({
  postId,
  pid = 0,
  replyName = "",
  submitText,
  onSubmit,
  onCancel,
  adminInfo,
}: CommentFormProps) {
  const [submitting, setSubmitting] = useState(false);
  const [showEmoji, setShowEmoji] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  // 记录最近一次光标位置，用于在失焦后仍能定位插入 emoji
  const cursorPosRef = useRef<number | null>(null);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { name: "", email: "", website: "", content: "" },
    mode: "onChange",
  });

  // 记录内容输入框失焦时的光标位置
  const handleContentBlur = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    cursorPosRef.current = e.target.selectionStart;
  };

  // 在内容光标位置插入 emoji
  const insertEmoji = (emoji: string) => {
    const textarea = textareaRef.current;
    const currentValue = form.getValues("content") || "";
    // 优先使用当前光标位置，否则追加到末尾
    const pos = textarea?.selectionStart ?? cursorPosRef.current ?? currentValue.length;
    const newValue = currentValue.slice(0, pos) + emoji + currentValue.slice(pos);
    form.setValue("content", newValue, { shouldValidate: true });
    // 插入后关闭 emoji 面板
    setShowEmoji(false);
    // 恢复光标到插入的 emoji 之后并聚焦
    const newPos = pos + emoji.length;
    cursorPosRef.current = newPos;
    setTimeout(() => {
      textarea?.focus();
      textarea?.setSelectionRange(newPos, newPos);
    }, 0);
  };

  const handleSubmit = async (data: FormValues) => {
    // 访客模式下昵称必填（管理员模式使用管理员信息，跳过校验）
    if (!adminInfo && !data.name) {
      form.setError("name", { type: "manual", message: "请输入昵称" });
      return;
    }
    setSubmitting(true);
    try {
      const name = adminInfo ? adminInfo.name : (data.name as string);
      await commentCreateApi({
        post_id: postId,
        name,
        email: adminInfo ? adminInfo.email : data.email || "",
        website: adminInfo ? adminInfo.host : data.website || "",
        content: data.content,
        pid,
        reply_name: replyName,
      });
      // 重置内容，非管理员时保留昵称/邮箱/网址以便连续评论
      form.reset({
        name: adminInfo ? "" : name,
        email: adminInfo ? "" : data.email || "",
        website: adminInfo ? "" : data.website || "",
        content: "",
      });
      await onSubmit();
    } catch {
      // 错误由 request 拦截器处理（dialog.message 弹窗）
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form
      method="post"
      onSubmit={form.handleSubmit(handleSubmit)}
      className="p-4 rounded-lg border border-[#e5e7eb] bg-[#fafafa]"
    >
      <FieldGroup className="gap-3">
        {/* 管理员身份时隐藏昵称/邮箱/网址，自动使用管理员信息 */}
        {!adminInfo && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          <Controller
            name="name"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <div className="flex items-center gap-2">
                  <FieldLabel htmlFor={field.name} className="text-xs text-[#6b7280] shrink-0">
                    昵称 *
                  </FieldLabel>
                  <Input
                    {...field}
                    id={field.name}
                    placeholder="昵称"
                    className="h-8 bg-white flex-1 min-w-0"
                  />
                </div>
                {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
              </Field>
            )}
          />
          <Controller
            name="email"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <div className="flex items-center gap-2">
                  <FieldLabel htmlFor={field.name} className="text-xs text-[#6b7280] shrink-0">
                    邮箱
                  </FieldLabel>
                  <Input
                    {...field}
                    id={field.name}
                    type="email"
                    placeholder="email@example.com"
                    className="h-8 bg-white flex-1 min-w-0"
                  />
                </div>
                {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
              </Field>
            )}
          />
          <Controller
            name="website"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <div className="flex items-center gap-2">
                  <FieldLabel htmlFor={field.name} className="text-xs text-[#6b7280] shrink-0">
                    网址
                  </FieldLabel>
                  <Input
                    {...field}
                    id={field.name}
                    placeholder="https://example.com"
                    className="h-8 bg-white flex-1 min-w-0"
                  />
                </div>
                {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
              </Field>
            )}
          />
        </div>
        )}
        <Controller
          name="content"
          control={form.control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <Textarea
                {...field}
                ref={(el) => {
                  field.ref(el);
                  textareaRef.current = el;
                }}
                onBlur={(e) => {
                  field.onBlur();
                  handleContentBlur(e);
                }}
                id={field.name}
                rows={3}
                placeholder="说点什么吧..."
                className="resize-y bg-white"
              />
              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />
        <div className="relative flex items-center justify-between">
          {/* emoji 选择器 */}
          <div className="relative">
            <Button
              type="button"
              variant="ghost"
              size="icon-sm"
              onClick={() => setShowEmoji((v) => !v)}
              title="插入表情"
              className="text-[#6b7280]"
            >
              <Smile size={16} />
            </Button>
            {showEmoji && (
              <>
                {/* 点击遮罩关闭面板 */}
                <div
                  className="fixed inset-0 z-10"
                  onClick={() => setShowEmoji(false)}
                />
                <div className="absolute bottom-full left-0 z-20 mb-1 w-[280px] rounded-lg border border-[#e5e7eb] bg-white p-2 shadow-lg">
                  <div className="grid grid-cols-8 gap-1">
                    {EMOJI_LIST.map((emoji) => (
                      <button
                        key={emoji}
                        type="button"
                        onClick={() => insertEmoji(emoji)}
                        className="grid h-8 w-8 place-items-center rounded text-lg hover:bg-[#eef6fb]"
                      >
                        {emoji}
                      </button>
                    ))}
                  </div>
                </div>
              </>
            )}
          </div>
          <div className="flex items-center gap-2">
            {onCancel && (
              <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
                取消
              </Button>
            )}
            <Button type="submit" size="sm" disabled={submitting}>
              {submitting && <Loader2 className="animate-spin" size={14} />}
              {submitText}
            </Button>
          </div>
        </div>
      </FieldGroup>
    </form>
  );
}

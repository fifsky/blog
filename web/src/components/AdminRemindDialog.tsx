import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select";
import { Field, FieldLabel, FieldError, FieldGroup, FieldContent } from "./ui/field";
import { Textarea } from "./ui/textarea";
import { Button } from "./ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "./ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./ui/tabs";
import { remindType, monthFormat, weekFormat, numFormat, cronToForm } from "@/utils/remind_date";
import { RemindItem } from "@/types/openapi";
import { useEffect, useState, useRef } from "react";

interface AdminRemindDialogProps {
  isOpen: boolean;
  onClose: () => void;
  item: RemindItem | undefined;
  onSubmit: (values: z.infer<typeof formSchema>) => Promise<void>;
  onSmartSubmit?: (content: string) => Promise<void>;
}

const formSchema = z.object({
  type: z.number(),
  month: z.number().optional(),
  week: z.number().optional(),
  day: z.number().optional(),
  hour: z.number().optional(),
  minute: z.number().optional(),
  content: z.string().min(1, "请输入提醒内容"),
});

export function AdminRemindDialog({
  isOpen,
  onClose,
  item,
  onSubmit,
  onSmartSubmit,
}: AdminRemindDialogProps) {
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState("smart");
  const [smartContent, setSmartContent] = useState("");
  const smartTextareaRef = useRef<HTMLTextAreaElement>(null);
  const manualTextareaRef = useRef<HTMLTextAreaElement>(null);

  const initialForm = item?.cron
    ? cronToForm(item.cron)
    : { type: 0, month: 1, week: 1, day: 1, hour: 0, minute: 0 };

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      type: initialForm.type,
      month: initialForm.month,
      week: initialForm.week,
      day: initialForm.day,
      hour: initialForm.hour,
      minute: initialForm.minute,
      content: item?.content ?? "",
    },
    mode: "onChange",
  });

  // 当 item 属性变化时，更新表单值
  useEffect(() => {
    if (isOpen) {
      if (item) {
        setActiveTab("manual");
      } else {
        setActiveTab("smart");
      }
      setSmartContent("");
    }

    if (item?.cron) {
      const v = cronToForm(item.cron);
      form.reset({
        type: v.type,
        month: v.month,
        week: v.week,
        day: v.day,
        hour: v.hour,
        minute: v.minute,
        content: item?.content ?? "",
      });
    } else {
      form.reset({
        type: 0,
        month: 1,
        week: 1,
        day: 1,
        hour: 0,
        minute: 0,
        content: "",
      });
    }
  }, [item, form]);

  // 当 activeTab 变化时，自动让对应的输入框获得焦点
  useEffect(() => {
    // 使用 setTimeout 确保 DOM 已经渲染完毕
    const timer = setTimeout(() => {
      if (activeTab === "smart") {
        smartTextareaRef.current?.focus();
      } else if (activeTab === "manual") {
        manualTextareaRef.current?.focus();
      }
    }, 100);
    return () => clearTimeout(timer);
  }, [activeTab, isOpen]);

  const intRemindType = Number(form.watch("type"));

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    setLoading(true);
    try {
      await onSubmit(values);
      onClose();
      form.reset();
    } finally {
      setLoading(false);
    }
  };

  const handleSmartSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!smartContent.trim() || !onSmartSubmit) return;
    setLoading(true);
    try {
      await onSmartSubmit(smartContent);
      onClose();
      setSmartContent("");
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    onClose();
    form.reset();
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{item?.id ? "编辑提醒" : "新增提醒"}</DialogTitle>
        </DialogHeader>
        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          {!item?.id && (
            <TabsList className="grid w-full grid-cols-2 mb-4">
              <TabsTrigger value="smart">智能提醒</TabsTrigger>
              <TabsTrigger value="manual">手动录入</TabsTrigger>
            </TabsList>
          )}
          <TabsContent value="smart">
            <form onSubmit={handleSmartSubmit} className="w-full px-1">
              <FieldGroup>
                <Field orientation="vertical">
                  <FieldContent>
                    <Textarea
                      ref={smartTextareaRef}
                      value={smartContent}
                      onChange={(e) => setSmartContent(e.target.value)}
                      placeholder="请描述你的提醒，例如：每天早上9点提醒我喝水"
                      rows={5}
                    />
                  </FieldContent>
                </Field>
                <Field orientation="horizontal">
                  <Button
                    type="submit"
                    size="sm"
                    loading={loading}
                    disabled={!smartContent.trim() || loading}
                  >
                    智能生成并添加
                  </Button>
                </Field>
              </FieldGroup>
            </form>
          </TabsContent>
          <TabsContent value="manual">
            <form
              className="w-full px-1"
              method="post"
              autoComplete="off"
              onSubmit={form.handleSubmit(handleSubmit)}
            >
              <FieldGroup>
                <Controller
                  name="type"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field orientation="vertical" data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="type">提醒类别</FieldLabel>
                      <FieldContent>
                        <Select
                          value={String(field.value)}
                          onValueChange={(v) => field.onChange(Number(v))}
                        >
                          <SelectTrigger size="sm">
                            <SelectValue placeholder="请选择提醒类别" />
                          </SelectTrigger>
                          <SelectContent>
                            {Object.entries(remindType).map(([k, v]) => (
                              <SelectItem key={k} value={k}>
                                {v}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </FieldContent>
                    </Field>
                  )}
                />
                <Field orientation="vertical">
                  <FieldLabel>提醒时间</FieldLabel>
                  <FieldContent>
                    <div className="flex flex-row gap-2">
                      {[0, 6].includes(intRemindType) && (
                        <Controller
                          name="month"
                          control={form.control}
                          render={({ field }) => (
                            <Select
                              value={String(field.value)}
                              onValueChange={(v) => field.onChange(Number(v))}
                            >
                              <SelectTrigger size="sm">
                                <SelectValue placeholder="月" />
                              </SelectTrigger>
                              <SelectContent>
                                {Array.from({ length: 12 }, (_, i) => i + 1).map((m) => (
                                  <SelectItem key={m} value={String(m)}>
                                    {monthFormat[m]}月
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                      )}
                      {[4].includes(intRemindType) && (
                        <Controller
                          name="week"
                          control={form.control}
                          render={({ field }) => (
                            <Select
                              value={String(field.value)}
                              onValueChange={(v) => field.onChange(Number(v))}
                            >
                              <SelectTrigger size="sm">
                                <SelectValue placeholder="周" />
                              </SelectTrigger>
                              <SelectContent>
                                {Array.from({ length: 7 }, (_, i) => i + 1).map((d) => (
                                  <SelectItem key={d} value={String(d)}>
                                    周{weekFormat[d]}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                      )}
                      {[0, 5, 6].includes(intRemindType) && (
                        <Controller
                          name="day"
                          control={form.control}
                          render={({ field }) => (
                            <Select
                              value={String(field.value)}
                              onValueChange={(v) => field.onChange(Number(v))}
                            >
                              <SelectTrigger size="sm">
                                <SelectValue placeholder="日" />
                              </SelectTrigger>
                              <SelectContent>
                                {Array.from({ length: 31 }, (_, i) => i + 1).map((d) => (
                                  <SelectItem key={d} value={String(d)}>
                                    {numFormat(d)}日
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                      )}
                      {[0, 3, 4, 5, 6].includes(intRemindType) && (
                        <Controller
                          name="hour"
                          control={form.control}
                          render={({ field }) => (
                            <Select
                              value={String(field.value)}
                              onValueChange={(v) => field.onChange(Number(v))}
                            >
                              <SelectTrigger size="sm">
                                <SelectValue placeholder="时" />
                              </SelectTrigger>
                              <SelectContent>
                                {Array.from({ length: 24 }, (_, i) => i).map((d) => (
                                  <SelectItem key={d} value={String(d)}>
                                    {numFormat(d)}时
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                      )}
                      {[0, 2, 3, 4, 5, 6].includes(intRemindType) && (
                        <Controller
                          name="minute"
                          control={form.control}
                          render={({ field }) => (
                            <Select
                              value={String(field.value)}
                              onValueChange={(v) => field.onChange(Number(v))}
                            >
                              <SelectTrigger size="sm">
                                <SelectValue placeholder="分" />
                              </SelectTrigger>
                              <SelectContent>
                                {Array.from({ length: 60 }, (_, i) => i).map((d) => (
                                  <SelectItem key={d} value={String(d)}>
                                    {numFormat(d)}分
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                      )}
                    </div>
                  </FieldContent>
                </Field>
                <Controller
                  name="content"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field orientation="vertical" data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor={field.name}>提醒内容</FieldLabel>
                      <FieldContent>
                        <Textarea
                          {...field}
                          id={field.name}
                          rows={5}
                          ref={(e) => {
                            field.ref(e);
                            // @ts-ignore
                            manualTextareaRef.current = e;
                          }}
                        />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </FieldContent>
                    </Field>
                  )}
                />
                <Field orientation="horizontal">
                  <Button type="submit" size="sm" loading={loading} disabled={loading}>
                    {item?.id ? "修改" : "添加"}
                  </Button>
                  {item?.id && (
                    <Button
                      size={"sm"}
                      variant="link"
                      onClick={(e) => {
                        e.preventDefault();
                        handleCancel();
                      }}
                      disabled={loading}
                    >
                      取消
                    </Button>
                  )}
                </Field>
              </FieldGroup>
            </form>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}

export default AdminRemindDialog;

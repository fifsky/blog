
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import {
  Field,
  FieldLabel,
  FieldError,
  FieldGroup,
  FieldContent,
} from "./ui/field";
import { Textarea } from "./ui/textarea";
import { Button } from "./ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "./ui/dialog";
import {
  remindType,
  monthFormat,
  weekFormat,
  numFormat,
} from "@/utils/remind_date";
import { RemindItem } from "@/types/openapi";
import { useEffect } from "react";

interface AdminRemindDialogProps {
  isOpen: boolean;
  onClose: () => void;
  item: RemindItem | undefined;
  onSubmit: (values: z.infer<typeof formSchema>) => Promise<void>;
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
}: AdminRemindDialogProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      type: item?.type ?? 0,
      month: item?.month ?? 1,
      week: item?.week ?? 1,
      day: item?.day ?? 1,
      hour: item?.hour ?? 0,
      minute: item?.minute ?? 0,
      content: item?.content ?? "",
    },
    mode: "onChange",
  });

  // 当 item 属性变化时，更新表单值
  useEffect(() => {
    if (item) {
      form.reset({
        type: item?.type ?? 0,
        month: item?.month ?? 1,
        week: item?.week ?? 1,
        day: item?.day ?? 1,
        hour: item?.hour ?? 0,
        minute: item?.minute ?? 0,
        content: item?.content ?? "",
      });
    }
  }, [item, form]);

  const intRemindType = Number(form.watch("type"));

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    await onSubmit(values);
    onClose();
    form.reset();
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
                <Field
                  orientation="vertical"
                  data-invalid={fieldState.invalid}
                >
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
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
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
                            {Array.from({ length: 12 }, (_, i) => i + 1).map(
                              (m) => (
                                <SelectItem key={m} value={String(m)}>
                                  {monthFormat[m]}月
                                </SelectItem>
                              )
                            )}
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
                            {Array.from({ length: 7 }, (_, i) => i + 1).map(
                              (d) => (
                                <SelectItem key={d} value={String(d)}>
                                  周{weekFormat[d]}
                                </SelectItem>
                              )
                            )}
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
                            {Array.from({ length: 31 }, (_, i) => i + 1).map(
                              (d) => (
                                <SelectItem key={d} value={String(d)}>
                                  {numFormat(d)}日
                                </SelectItem>
                              )
                            )}
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
                            {Array.from({ length: 24 }, (_, i) => i).map(
                              (d) => (
                                <SelectItem key={d} value={String(d)}>
                                  {numFormat(d)}时
                                </SelectItem>
                              )
                            )}
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
                            {Array.from({ length: 60 }, (_, i) => i).map(
                              (d) => (
                                <SelectItem key={d} value={String(d)}>
                                  {numFormat(d)}分
                                </SelectItem>
                              )
                            )}
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
                <Field
                  orientation="vertical"
                  data-invalid={fieldState.invalid}
                >
                  <FieldLabel htmlFor={field.name}>提醒内容</FieldLabel>
                  <FieldContent>
                    <Textarea {...field} id={field.name} rows={5} />
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
                  </FieldContent>
                </Field>
              )}
            />
            <Field orientation="horizontal">
              <Button type="submit" size="sm">
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
                >
                  取消
                </Button>
              )}
            </Field>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  );
}

export default AdminRemindDialog;
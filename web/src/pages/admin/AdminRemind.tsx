import { useEffect, useState } from "react";
import {
  remindDeleteApi,
  remindListApi,
  remindCreateApi,
  remindUpdateApi,
} from "@/service";
import dayjs from "dayjs";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Field,
  FieldLabel,
  FieldError,
  FieldGroup,
  FieldContent,
} from "@/components/ui/field";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";

export default function AdminRemind() {
  const [list, setList] = useState<any[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [item, setItem] = useState<any>({});
  const formSchema = z.object({
    type: z.number(),
    month: z.number().optional(),
    week: z.number().optional(),
    day: z.number().optional(),
    hour: z.number().optional(),
    minute: z.number().optional(),
    content: z.string().min(1, "请输入提醒内容"),
  });
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      type: 0,
      month: 1,
      week: 1,
      day: 1,
      hour: 0,
      minute: 0,
      content: "",
    },
    mode: "onChange",
  });
  const remindType: Record<number, string> = {
    0: "固定",
    1: "每分钟",
    2: "每小时",
    3: "每天",
    4: "每周",
    5: "每月",
    6: "每年",
  };
  const monthFormat: Record<number, string> = {
    1: "01",
    2: "02",
    3: "03",
    4: "04",
    5: "05",
    6: "06",
    7: "07",
    8: "08",
    9: "09",
    10: "10",
    11: "11",
    12: "12",
  };
  const weekFormat: Record<number, string> = {
    1: "一",
    2: "二",
    3: "三",
    4: "四",
    5: "五",
    6: "六",
    7: "日",
  };
  const intRemindType = Number(form.watch("type"));
  const numFormat = (n: number) => (n < 10 ? "0" + n : String(n));
  const remindTimeFormat = (v: any) => {
    let str = "";
    switch (v.type) {
      case 0:
        str =
          dayjs(v.created_at).year() +
          "年" +
          monthFormat[v.month] +
          "月" +
          numFormat(v.day) +
          "日 " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
      case 3:
        str = numFormat(v.hour) + "时" + numFormat(v.minute) + "分";
        break;
      case 4:
        str =
          "周" +
          weekFormat[v.week] +
          " " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
      case 5:
        str =
          numFormat(v.day) +
          "日 " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
      case 6:
        str =
          monthFormat[v.month] +
          "月" +
          numFormat(v.day) +
          "日 " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
    }
    return str;
  };
  const loadList = async () => {
    const ret = await remindListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.pageTotal || 0);
  };
  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it || {});
    form.reset({
      type: it?.type ?? 0,
      month: it?.month ?? 1,
      week: it?.week ?? 1,
      day: it?.day ?? 1,
      hour: it?.hour ?? 0,
      minute: it?.minute ?? 0,
      content: it?.content ?? "",
    });
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await remindDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => {
    setItem({});
    form.reset({
      type: 0,
      month: 1,
      week: 1,
      day: 1,
      hour: 0,
      minute: 0,
      content: "",
    });
  };
  const submit = async (values: z.infer<typeof formSchema>) => {
    const { id } = item || {};
    const data = {
      id,
      type: Number(values.type),
      content: values.content,
      month: values.month,
      week: values.week,
      day: values.day,
      hour: values.hour,
      minute: values.minute,
    };
    if (id) await remindUpdateApi(data);
    else await remindCreateApi(data);
    cancel();
    loadList();
  };
  useEffect(() => {
    loadList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);
  return (
    <div>
      <h2>管理提醒</h2>
      <div className="flex justify-between">
        <div className="w-[700px]">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          <table className="list">
            <tbody>
              <tr>
                <th style={{ width: 20 }}>&nbsp;</th>
                <th style={{ width: 80 }}>提醒类别</th>
                <th style={{ width: 180 }}>时间</th>
                <th>内容</th>
                <th style={{ width: 90 }}>操作</th>
              </tr>
              {list.length === 0 && (
                <tr>
                  <td colSpan={7} align="center">
                    还没有提醒！
                  </td>
                </tr>
              )}
              {list.length > 0 &&
                list.map((v) => (
                  <tr key={v.id}>
                    <td>
                      <input type="checkbox" name="ids" value={v.id} />
                    </td>
                    <td>{remindType[v.type]}</td>
                    <td>{remindTimeFormat(v)}</td>
                    <td>{v.content}</td>
                    <td>
                      <a
                        href="#"
                        onClick={(e) => {
                          e.preventDefault();
                          editItem(v.id);
                        }}
                      >
                        编辑
                      </a>
                      <span className="px-1.5 text-[#ccc]">|</span>
                      <a
                        href="#"
                        onClick={(e) => {
                          e.preventDefault();
                          deleteItem(v.id);
                        }}
                      >
                        删除
                      </a>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
          <div className="my-[10px] flex items-center justify-between">
            <BatchHandle />
            <Paginate page={page} pageTotal={pageTotal} onChange={setPage} />
          </div>
        </div>
        <div className="w-[250px]" style={{ paddingTop: 31 }}>
          <form
            className="w-full px-1"
            method="post"
            autoComplete="off"
            onSubmit={form.handleSubmit(submit)}
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
                  <div className="flex flex-col gap-2">
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
                  {item.id ? "修改" : "添加"}
                </Button>
                {item.id && (
                  <Button
                    size={"sm"}
                    variant="link"
                    onClick={(e) => {
                      e.preventDefault();
                      cancel();
                    }}
                  >
                    取消
                  </Button>
                )}
              </Field>
            </FieldGroup>
          </form>
        </div>
      </div>
    </div>
  );
}

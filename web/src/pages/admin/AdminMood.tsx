import { useEffect, useState } from "react";
import {
  moodDeleteApi,
  moodListApi,
  moodCreateApi,
  moodUpdateApi,
} from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Field,
  FieldLabel,
  FieldError,
  FieldGroup,
  FieldContent,
} from "@/components/ui/field";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";

export default function AdminMood() {
  const [list, setList] = useState<any[]>([]);
  const [item, setItem] = useState<any>({});
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const formSchema = z.object({
    content: z.string().min(1, "请输入心情内容"),
  });
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { content: "" },
    mode: "onChange",
  });
  const loadList = async () => {
    const ret = await moodListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it || {});
    form.reset({ content: it?.content || "" });
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await moodDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => {
    setItem({});
    form.reset({ content: "" });
  };
  const submit = async (values: z.infer<typeof formSchema>) => {
    const { id } = item || {};
    if (id) await moodUpdateApi({ id, content: values.content });
    else await moodCreateApi({ content: values.content });
    cancel();
    loadList();
  };
  useEffect(() => {
    loadList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);
  return (
    <div>
      <h2>管理心情</h2>
      <div className="flex justify-between">
        <div className="w-[700px]">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          <table className="list">
            <tbody>
              <tr>
                <th style={{ width: 20 }}>&nbsp;</th>
                <th style={{ width: 80 }}>作者</th>
                <th>心情</th>
                <th style={{ width: 180 }}>日期</th>
                <th style={{ width: 90 }}>操作</th>
              </tr>
              {list.length === 0 && (
                <tr>
                  <td colSpan={7} align="center">
                    还没有心情！
                  </td>
                </tr>
              )}
              {list.length > 0 &&
                list.map((v) => (
                  <tr key={v.id}>
                    <td>
                      <input type="checkbox" name="ids" value={v.id} />
                    </td>
                    <td>{v.user.name}</td>
                    <td>{v.content}</td>
                    <td>{v.created_at}</td>
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
          <div className="my-2.5 flex items-center justify-between">
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
                name="content"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field
                    orientation="vertical"
                    data-invalid={fieldState.invalid}
                  >
                    <FieldLabel htmlFor={field.name}>发表心情</FieldLabel>
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

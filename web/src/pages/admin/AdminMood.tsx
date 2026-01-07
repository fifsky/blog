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
import { CTable, Column } from "@/components/CTable";
import { MoodItem } from "@/types/openapi";

export default function AdminMood() {
  const [list, setList] = useState<MoodItem[]>([]);
  const [item, setItem] = useState<MoodItem | undefined>();
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
    setItem(it);
    form.reset({ content: it?.content || "" });
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await moodDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => {
    setItem(undefined);
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
  }, [page]);

  // 定义表格列配置
  const columns: Column<MoodItem>[] = [
    {
      title: <div style={{ width: 20 }}>&nbsp;</div>,
      key: "id",
      render: (_, record) => (
        <input type="checkbox" name="ids" value={record.id} />
      )
    },
    {
      title: <div style={{ width: 80 }}>作者</div>,
      key: "user.name"
    },
    {
      title: "心情",
      key: "content"
    },
    {
      title: <div style={{ width: 180 }}>日期</div>,
      key: "created_at"
    },
    {
      title: <div style={{ width: 90 }}>操作</div>,
      key: "id",
      render: (_,record) => (
        <>
          <Button 
            variant={"link"}
            className="p-0 m-0 h-auto text-[13px]"
            onClick={(e) => {
              e.preventDefault();
              editItem(record.id);
            }}
          >
            编辑
          </Button>
          <span className="px-1.5 text-[#ccc]">|</span>
          <Button 
            variant={"link"}
            className="p-0 m-0 h-auto text-[13px]"
            onClick={(e) => {
              e.preventDefault();
              deleteItem(record.id);
            }}
          >
            删除
          </Button>
        </>
      )
    }
  ];

  return (
    <div>
      <h2 className="border-b border-b-[#cccccc] text-base">管理心情</h2>
      <div className="flex justify-between">
        <div className="w-[700px]">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          {/* 使用自定义表格组件 */}
          <CTable data={list} columns={columns} />
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
                  {item?.id ? "修改" : "添加"}
                </Button>
                {item?.id && (
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

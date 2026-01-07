import { useEffect, useState } from "react";
import {
  cateDeleteApi,
  cateListApi,
  cateCreateApi,
  cateUpdateApi,
} from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Field,
  FieldLabel,
  FieldError,
  FieldGroup,
  FieldDescription,
  FieldContent,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { CTable, Column } from "@/components/CTable";

export default function AdminCate() {
  const [list, setList] = useState<any[]>([]);
  const [item, setItem] = useState<any>({});
  const formSchema = z.object({
    name: z.string().min(1, "请输入分类名称"),
    domain: z
      .string()
      .regex(/^[a-z][a-z0-9-]*$/, "缩略名需字母开头，包含小写字母、数字或-"),
    desc: z.string().optional(),
  });
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { name: "", domain: "", desc: "" },
    mode: "onChange",
  });
  const loadList = async () => {
    const ret = await cateListApi({});
    setList(ret.list || []);
  };
  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it || {});
    form.reset({
      name: it?.name || "",
      domain: it?.domain || "",
      desc: it?.desc || "",
    });
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await cateDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => setItem({});
  const submit = async (values: z.infer<typeof formSchema>) => {
    const { id } = item || {};
    if (id) await cateUpdateApi({ id, ...values });
    else await cateCreateApi(values);
    setItem({});
    form.reset({ name: "", domain: "", desc: "" });
    loadList();
  };
  useEffect(() => {
    loadList();
  }, []);

  // 定义表格列配置
  const columns: Column<any>[] = [
    {
      title: <div style={{ width: 20 }}>&nbsp;</div>,
      key: "id",
      render: (_, record) => (
        <input type="checkbox" name="ids" value={record.id} />
      )
    },
    {
      title: "分类名",
      key: "name"
    },
    {
      title: <div style={{ width: 90 }}>缩略名</div>,
      key: "domain"
    },
    {
      title: <div style={{ width: 60 }}>文章数</div>,
      key: "num",
      render: (value) => (
        <Badge variant="secondary">{value}</Badge>
      )
    },
    {
      title: <div style={{ width: 90 }}>操作</div>,
      key: "id",
      render: (_, record) => (
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
      <h2 className="border-b border-b-[#cccccc] text-base">管理分类</h2>
      <div className="flex justify-between">
        <div className="w-[700px]">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          {/* 使用自定义表格组件 */}
          <CTable data={list} columns={columns} />
          <div className="my-[10px] flex items-center justify-between">
            <BatchHandle />
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
                name="name"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field
                    orientation="vertical"
                    data-invalid={fieldState.invalid}
                  >
                    <FieldLabel htmlFor={field.name}>分类名称</FieldLabel>
                    <FieldContent>
                      <Input {...field} id={field.name} />
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </FieldContent>
                  </Field>
                )}
              />
              <Controller
                name="domain"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field
                    orientation="vertical"
                    data-invalid={fieldState.invalid}
                  >
                    <FieldLabel htmlFor={field.name}>分类缩略名</FieldLabel>
                    <FieldContent>
                      <Input {...field} id={field.name} />
                      <FieldDescription>
                        缩略名，使用字母开头([a-z][0-9]-)
                      </FieldDescription>
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </FieldContent>
                  </Field>
                )}
              />
              <Controller
                name="desc"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field
                    orientation="vertical"
                    data-invalid={fieldState.invalid}
                  >
                    <FieldLabel htmlFor={field.name}>分类描述</FieldLabel>
                    <FieldContent>
                      <Textarea {...field} id={field.name} rows={5} />
                      <FieldDescription>
                        描述将在分类meta中显示
                      </FieldDescription>
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

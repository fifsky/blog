import { useEffect, useState } from "react";
import {
  linkDeleteApi,
  linkListApi,
  linkCreateApi,
  linkUpdateApi,
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

export default function AdminLink() {
  const [list, setList] = useState<any[]>([]);
  const [item, setItem] = useState<any>({});
  const formSchema = z.object({
    name: z.string().min(1, "请输入链接名称"),
    url: z.string().url("请输入正确的链接地址"),
    desc: z.string().optional(),
  });
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { name: "", url: "", desc: "" },
    mode: "onChange",
  });
  const loadList = async () => {
    const ret = await linkListApi({});
    setList(ret.list || []);
  };
  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it || {});
    form.reset({
      name: it?.name || "",
      url: it?.url || "",
      desc: it?.desc || "",
    });
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await linkDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => {
    setItem({});
    form.reset({ name: "", url: "", desc: "" });
  };
  const submit = async (values: z.infer<typeof formSchema>) => {
    const { id } = item || {};
    if (id) await linkUpdateApi({ id, ...values });
    else await linkCreateApi(values);
    cancel();
    loadList();
  };
  useEffect(() => {
    loadList();
  }, []);
  return (
    <div>
      <h2>管理链接</h2>
      <div className="flex justify-between">
        <div className="w-[700px]">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          <table className="list">
            <tbody>
              <tr>
                <th style={{ width: 20 }}>&nbsp;</th>
                <th style={{ width: 120 }}>连接名</th>
                <th>地址</th>
                <th style={{ width: 90 }}>操作</th>
              </tr>
              {list.length === 0 && (
                <tr>
                  <td colSpan={7} align="center">
                    还没有链接！
                  </td>
                </tr>
              )}
              {list.length > 0 &&
                list.map((v) => (
                  <tr key={v.id}>
                    <td>
                      <input type="checkbox" name="ids" value={v.id} />
                    </td>
                    <td>
                      <a href={v.url} target="_blank" rel="noreferrer">
                        {v.name}
                      </a>
                    </td>
                    <td>
                      <a href={v.url}>{v.url}</a>
                    </td>
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
                    <FieldLabel htmlFor={field.name}>链接名称</FieldLabel>
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
                name="url"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field
                    orientation="vertical"
                    data-invalid={fieldState.invalid}
                  >
                    <FieldLabel htmlFor={field.name}>链接地址</FieldLabel>
                    <FieldContent>
                      <Input {...field} id={field.name} />
                      <FieldDescription>
                        例如：http://fifsky.com/
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
                    <FieldLabel htmlFor={field.name}>链接描述</FieldLabel>
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

import React from "react";
import { settingApi, settingUpdateApi } from "@/service";
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
import { Alert, AlertTitle } from "@/components/ui/alert";
import { CheckCircle2Icon } from "lucide-react";
import { useAsyncEffect } from "@/hooks";

export default function AdminIndex() {
  const [loading, setLoading] = React.useState(false);

  const formSchema = z.object({
    site_name: z.string().min(1, "请输入站点名称"),
    site_desc: z.string().optional(),
    site_keyword: z.string().optional(),
    post_num: z.string().regex(/^\d+$/, "请输入数字"),
    map_regions: z.string().optional(),
    map_footprints: z.string().optional(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      site_name: "",
      site_desc: "",
      site_keyword: "",
      post_num: "",
      map_regions: "",
      map_footprints: "",
    },
    mode: "onChange",
  });

  const [showMessage, setShowMessage] = React.useState(false);

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setLoading(true);
    try {
      await settingUpdateApi({ kv: values });
      setShowMessage(true);
      setTimeout(() => setShowMessage(false), 3000);
    } finally {
      setLoading(false);
    }
  };
  useAsyncEffect(async () => {
    const data = await settingApi();
    form.reset({
      site_name: data.kv?.site_name || "",
      site_desc: data.kv?.site_desc || "",
      site_keyword: data.kv?.site_keyword || "",
      post_num: data.kv?.post_num || "",
      map_regions: data.kv?.map_regions || "",
      map_footprints: data.kv?.map_footprints || "",
    });
  }, [form]);
  return (
    <div>
      <title>站点设置 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">站点设置</h2>
      {showMessage && (
        <Alert variant="success" className="mt-2">
          <CheckCircle2Icon />
          <AlertTitle>保存成功</AlertTitle>
        </Alert>
      )}
      <div className="max-w-xl mx-auto mt-3">
        <form method="post" autoComplete="off" onSubmit={form.handleSubmit(onSubmit)}>
          <FieldGroup>
            <Controller
              name="site_name"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>站点名称</FieldLabel>
                  <FieldContent>
                    <Input {...field} id={field.name} />
                    <FieldDescription>站点的名称将显示在网页的标题处。</FieldDescription>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />
            <Controller
              name="site_desc"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>站点描述</FieldLabel>
                  <FieldContent>
                    <Textarea {...field} id={field.name} rows={3} />
                    <FieldDescription>站点描述将显示在网页代码的头部。</FieldDescription>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />
            <Controller
              name="site_keyword"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>关键字</FieldLabel>
                  <FieldContent>
                    <Input {...field} id={field.name} />
                    <FieldDescription>请以半角逗号","分割多个关键字。</FieldDescription>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />
            <Controller
              name="post_num"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>每页显示文章数</FieldLabel>
                  <FieldContent>
                    <Input {...field} id={field.name} style={{ width: 80 }} />
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />
            <Controller
              name="map_regions"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>点亮省份 (JSON数组)</FieldLabel>
                  <FieldContent>
                    <Textarea
                      className="h-20"
                      {...field}
                      id={field.name}
                      placeholder='例如: ["北京市", "上海市"]'
                    />
                    <FieldDescription>
                      输入JSON数组格式的省份名称列表，用于点亮地图背景。
                    </FieldDescription>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />
            <Controller
              name="map_footprints"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>足迹数据 (JSON数组)</FieldLabel>
                  <FieldContent>
                    <Textarea
                      className="h-20"
                      {...field}
                      id={field.name}
                      placeholder='例如: [{"name": "北京", "value": [116.40, 39.90]}]'
                    />
                    <FieldDescription>
                      输入JSON数组格式的足迹数据，包含名称和经纬度。
                    </FieldDescription>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />
            <Field orientation="responsive">
              <Button type="submit" size="sm" loading={loading}>
                保存
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </div>
    </div>
  );
}

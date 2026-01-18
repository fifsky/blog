import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { adminArticleDetailApi, articleCreateApi, articleUpdateApi, cateListApi } from "@/service";
import { useLocation, useNavigate, Link } from "react-router";
import { getApiUrl, getAccessToken } from "@/utils/common";
import { Editor } from "@bytemd/react";
import gfm from "@bytemd/plugin-gfm";
import mediumZoom from "@bytemd/plugin-medium-zoom";

import highlight from "@bytemd/plugin-highlight";
import { CateListItem } from "@/types/openapi";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Field, FieldContent, FieldDescription } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { cn } from "@/lib/utils";
import { useAsyncEffect } from "@/hooks";

const articleSchema = z.object({
  id: z.number().optional(),
  title: z.string().min(1, "标题不能为空"),
  cate_id: z.number().min(1, "请选择分类"),
  url: z.string().optional(),
  content: z.string().optional(),
  type: z.union([z.literal(1), z.literal(2)]),
  status: z.number().optional(),
});

type ArticleFormValues = z.infer<typeof articleSchema>;

// ByteMD 插件配置
const plugins = [gfm(), highlight(), mediumZoom()];

// 图片上传函数
const uploadImages = async (
  files: File[],
): Promise<{ url: string; alt: string; title: string }[]> => {
  const results: { url: string; alt: string; title: string }[] = [];

  for (const file of files) {
    const formData = new FormData();
    formData.append("uploadFile", file);

    const response = await fetch(getApiUrl("/blog/admin/upload"), {
      method: "POST",
      headers: {
        "Access-Token": getAccessToken(),
      },
      body: formData,
    });

    const data = await response.json();
    if (data.url) {
      results.push({
        url: data.url,
        alt: file.name,
        title: file.name,
      });
    }
  }

  return results;
};

export default function PostArticle() {
  const [cates, setCates] = useState<CateListItem[]>([]);
  const [loading, setLoading] = useState(false);

  const location = useLocation();
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search);

  const form = useForm<ArticleFormValues>({
    resolver: zodResolver(articleSchema),
    defaultValues: {
      type: 1,
      title: "",
      cate_id: 0,
      url: "",
      content: "",
    },
  });

  const isEditing = !!form.watch("id");
  const articleType = form.watch("type");
  const articleStatus = form.watch("status");
  const contentValue = form.watch("content");

  // 提交表单逻辑
  const submit = async (values: ArticleFormValues) => {
    setLoading(true);
    try {
      const { id, cate_id, title, content, type, url, status } = values;
      // 构造通用请求载荷，status 默认为 1 (发布)
      const payload = {
        cate_id,
        title,
        content: content || "",
        type,
        url: url || "",
        status: status || 1,
      };

      if (id) {
        // 编辑状态：调用更新接口，需要传入 id
        await articleUpdateApi({ ...payload, id });
      } else {
        // 新建状态：调用创建接口
        await articleCreateApi(payload);
      }
      navigate("/admin/articles");
    } finally {
      setLoading(false);
    }
  };

  // 发布文章：设置为发布状态并提交
  const handlePublish = () => {
    form.setValue("status", 1);
    form.handleSubmit(submit)();
  };

  // 保存草稿：设置为草稿状态并提交
  const handleSaveDraft = () => {
    form.setValue("status", 3); // 3 代表草稿
    form.handleSubmit(submit)();
  };

  useAsyncEffect(async () => {
    // 先加载分类列表
    const ret = await cateListApi({});
    const categories = ret.list || [];
    setCates(categories);

    // 然后处理文章详情
    if (params.get("id")) {
      const a = await adminArticleDetailApi({ id: parseInt(params.get("id")!) });
      form.reset({
        id: a.id,
        title: a.title || "",
        cate_id: a.cate_id || 0,
        url: a.url || "",
        content: a.content || "",
        type: a.type === 1 || a.type === 2 ? a.type : 1,
        status: a.status,
      });
    }
  }, []);

  // 在分类列表加载完成后设置默认分类
  useEffect(() => {
    // 只在非编辑模式下设置默认分类
    if (!isEditing && cates?.[0]?.id) {
      form.setValue("cate_id", cates[0].id);
    }
  }, [cates, isEditing, form]);

  return (
    <div className="max-w-5xl mx-auto">
      <title>{isEditing ? "编辑文章" : "撰写文章"}</title>
      <h2 className="border-b border-b-[#cccccc] text-base">
        {isEditing ? "编辑" : "撰写"}文章
        <Link to="/admin/articles" className="ml-3 text-[14px]">
          <i className="iconfont icon-undo" style={{ color: "#444" }}></i>
          返回列表
        </Link>
      </h2>
      <div className="mt-3">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(submit)} className="space-y-6">
            {/* 两行两列布局：第一行标题和类型，第二行分类和缩略名 */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* 第一行：标题 */}
              <FormField
                control={form.control}
                name="title"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      标题 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <FormControl>
                          <Input {...field} placeholder="请输入标题" maxLength={200} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                  </Field>
                )}
              />

              {/* 第一行：类型 */}
              <FormField
                control={form.control}
                name="type"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      类型 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <FormControl>
                          <RadioGroup
                            onValueChange={(value) => field.onChange(parseInt(value) as 1 | 2)}
                            defaultValue={field.value.toString()}
                            className="flex space-x-6"
                          >
                            <div className="flex items-center space-x-2">
                              <RadioGroupItem value="1" id="type-1" />
                              <label htmlFor="type-1" className="cursor-pointer">
                                文章
                              </label>
                            </div>
                            <div className="flex items-center space-x-2">
                              <RadioGroupItem value="2" id="type-2" />
                              <label htmlFor="type-2" className="cursor-pointer">
                                页面
                              </label>
                            </div>
                          </RadioGroup>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                  </Field>
                )}
              />

              {/* 第二行：分类 */}
              <FormField
                control={form.control}
                name="cate_id"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      分类 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <Select
                          onValueChange={(value) => field.onChange(parseInt(value))}
                          value={field.value.toString()}
                        >
                          <FormControl>
                            <SelectTrigger size={"sm"}>
                              <SelectValue placeholder="请选择分类" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {cates.map((v) => (
                              <SelectItem key={v.id} value={v.id.toString()}>
                                {v.name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                  </Field>
                )}
              />

              {/* 第二行：缩略名 */}
              {articleType === 2 && (
                <FormField
                  control={form.control}
                  name="url"
                  render={({ field }) => (
                    <Field>
                      <FormLabel>缩略名</FormLabel>
                      <FieldContent>
                        <FormItem>
                          <FormControl>
                            <Input
                              {...field}
                              placeholder="请输入缩略名"
                              maxLength={200}
                              className={cn("w-52")}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      </FieldContent>
                      <FieldDescription>
                        页面的URL名称，如http://domain.com/
                        <span style={{ color: "red" }}>about</span>
                      </FieldDescription>
                    </Field>
                  )}
                />
              )}
            </div>

            <FormField
              control={form.control}
              name="content"
              render={({ field }) => (
                <FormItem className="grid grid-cols-1">
                  <FormControl>
                    <div className="bytemd-editor-wrapper">
                      <Editor
                        value={contentValue || ""}
                        plugins={plugins}
                        onChange={(v) => field.onChange(v)}
                        uploadImages={uploadImages}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Field orientation="horizontal">
              <Button type="button" size={"sm"} loading={loading} onClick={handlePublish}>
                发布
              </Button>
              {(isEditing && articleStatus !== 1) || !isEditing ? (
                <Button
                  type="button"
                  size={"sm"}
                  variant="outline"
                  disabled={loading}
                  onClick={handleSaveDraft}
                >
                  保存草稿
                </Button>
              ) : null}
            </Field>
          </form>
        </Form>
      </div>
    </div>
  );
}

import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  articleDetailApi,
  articleCreateApi,
  articleUpdateApi,
  cateListApi,
} from "@/service";
import { useLocation, useNavigate, Link } from "react-router";
import "@wangeditor/editor/dist/css/style.css";
import { Editor, Toolbar } from "@wangeditor/editor-for-react";
import type {
  IDomEditor,
  IEditorConfig,
  IToolbarConfig,
} from "@wangeditor/editor";
import { getApiUrl, getAccessToken } from "@/utils/common";

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

const articleSchema = z.object({
  id: z.number().optional(),
  title: z.string().min(1, "标题不能为空"),
  cate_id: z.number().min(1, "请选择分类"),
  url: z.string().optional(),
  content: z.string().optional(),
  type: z.union([z.literal(1), z.literal(2)]),
});

type ArticleFormValues = z.infer<typeof articleSchema>;

export default function PostArticle() {
  const [cates, setCates] = useState<any[]>([]);
  const [editor, setEditor] = useState<IDomEditor | null>(null);

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

  const toolbarConfig: Partial<IToolbarConfig> = {
    excludeKeys: ["uploadVideo", "fontFamily", "lineHeight", "group-indent"],
  };

  const editorConfig: Partial<IEditorConfig> = {
    placeholder: "请输入内容...",
    MENU_CONF: {
      uploadImage: {
        server: getApiUrl("/api/admin/upload"),
        fieldName: "uploadFile",
        headers: { "Access-Token": getAccessToken() },
        withCredentials: false,
        maxFileSize: 10 * 1024 * 1024,
        allowedFileTypes: ["image/*"],
      },
    },
  };

  const submit = async (values: ArticleFormValues) => {
    const { id, cate_id, title, content, type, url } = values;
    if (id) {
      // 编辑状态：id 是必填的，content 可以是 undefined
      await articleUpdateApi({
        id,
        cate_id,
        title,
        content: content || "",
        type,
        url: url || "",
      });
    } else {
      // 新建状态：id 不需要，content 是必填的
      await articleCreateApi({
        cate_id,
        title,
        content: content || "",
        type,
        url: url || "",
      });
    }
    navigate("/admin/articles");
  };

  useEffect(() => {
    (async () => {
      // 先加载分类列表
      const ret = await cateListApi({});
      const categories = ret.list || [];
      setCates(categories);

      // 然后处理文章详情
      if (params.get("id")) {
        const a = await articleDetailApi({ id: parseInt(params.get("id")!) });
        form.reset({
          id: a.id,
          title: a.title || "",
          cate_id: a.cate_id || 0,
          url: a.url || "",
          content: a.content || "",
          type: a.type === 1 || a.type === 2 ? a.type : 1,
        });
      }
    })();
  }, []);

  // 在分类列表加载完成后设置默认分类
  useEffect(() => {
    // 只在非编辑模式下设置默认分类
    if (!isEditing && cates?.[0]?.id) {
      form.setValue("cate_id", cates[0].id);
    }
  }, [cates, isEditing, form]);

  useEffect(() => {
    return () => {
      if (editor) editor.destroy();
    };
  }, [editor]);

  return (
    <div id="articles" className="max-w-5xl mx-auto mt-3">
      <h2>
        {isEditing ? "编辑" : "撰写"}文章
        <Link to="/admin/articles">
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
                          <Input
                            {...field}
                            placeholder="请输入标题"
                            maxLength={200}
                          />
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
                            onValueChange={(value) =>
                              field.onChange(parseInt(value) as 1 | 2)
                            }
                            defaultValue={field.value.toString()}
                            className="flex space-x-6"
                          >
                            <div className="flex items-center space-x-2">
                              <RadioGroupItem value="1" id="type-1" />
                              <label
                                htmlFor="type-1"
                                className="cursor-pointer"
                              >
                                文章
                              </label>
                            </div>
                            <div className="flex items-center space-x-2">
                              <RadioGroupItem value="2" id="type-2" />
                              <label
                                htmlFor="type-2"
                                className="cursor-pointer"
                              >
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
                          onValueChange={(value) =>
                            field.onChange(parseInt(value))
                          }
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
                <Field>
                  <FieldContent>
                    <FormItem>
                      <FormControl>
                        <div className="border border-border rounded-md">
                          <Toolbar
                            editor={editor}
                            defaultConfig={toolbarConfig}
                            mode="default"
                            style={{ borderBottom: "1px solid #ddd" }}
                          />
                          <Editor
                            style={{ height: 500, overflowY: "hidden" }}
                            defaultConfig={editorConfig}
                            value={field.value || ""}
                            onCreated={(ed: IDomEditor) => setEditor(ed)}
                            onChange={(ed: IDomEditor) =>
                              field.onChange(ed.getHtml())
                            }
                            mode="default"
                          />
                        </div>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  </FieldContent>
                </Field>
              )}
            />

            <Field orientation="horizontal">
              <Button type="submit" size={"sm"}>
                发布
              </Button>
              <Button type="button" size={"sm"} variant="outline">
                保存草稿
              </Button>
            </Field>
          </form>
        </Form>
      </div>
    </div>
  );
}

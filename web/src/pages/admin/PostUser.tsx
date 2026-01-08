import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { userGetApi, userCreateApi, userUpdateApi } from "@/service";
import { useLocation, useNavigate, Link } from "react-router";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Field, FieldGroup, FieldContent, FieldDescription } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const userSchema = z
  .object({
    id: z.number().optional(),
    name: z.string().min(1, "用户名不能为空"),
    email: z.string().email("请输入有效的邮箱地址"),
    nick_name: z.string().optional(),
    password1: z.string().min(6, "密码至少6个字符"),
    password2: z.string().min(6, "确认密码至少6个字符"),
    type: z.union([z.literal(1), z.literal(2)]),
  })
  .refine((data) => data.password1 === data.password2, {
    message: "两次输入的密码不一致",
    path: ["password2"],
  });

type UserFormValues = z.infer<typeof userSchema>;

export default function PostUser() {
  const location = useLocation();
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search);
  const [loading, setLoading] = useState(false);

  const form = useForm<UserFormValues>({
    resolver: zodResolver(userSchema),
    defaultValues: {
      type: 1,
      name: "",
      email: "",
      nick_name: "",
      password1: "",
      password2: "",
    },
  });

  const submit = async (values: UserFormValues) => {
    setLoading(true);
    try {
      const { id, name, nick_name, password1, email, type } = values;
      const data = {
        id: id || 0,
        name,
        nick_name: nick_name || "",
        password: password1,
        email,
        type,
      };
      if (id) await userUpdateApi(data);
      else await userCreateApi(data);
      navigate("/admin/users");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    (async () => {
      if (params.get("id")) {
        const u = await userGetApi({ id: parseInt(params.get("id")!) });
        form.reset({
          id: u.id,
          name: u.name || "",
          email: u.email || "",
          nick_name: u.nick_name || "",
          type: u.type === 1 || u.type === 2 ? u.type : 1,
          password1: "",
          password2: "",
        });
      }
    })();
  }, []);

  const isEditing = !!form.watch("id");

  return (
    <div>
      <h2 className="border-b border-b-[#cccccc] text-base">
        {isEditing ? "编辑" : "新增"}用户
        <Link to="/admin/users" className="ml-3 text-[14px]">
          <i className="iconfont icon-undo" style={{ color: "#444" }}></i>
          返回列表
        </Link>
      </h2>
      <div className="max-w-xl mx-auto mt-3">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(submit)} className="space-y-6">
            <FieldGroup>
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      用户名 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <FormControl>
                          <Input {...field} placeholder="请输入用户名" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                    <FieldDescription>
                      此用户名将作为用户登录时所用的名称，请不要与系统中现有的用户名重复。
                    </FieldDescription>
                  </Field>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      邮箱 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <FormControl>
                          <Input {...field} placeholder="请输入邮箱" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                    <FieldDescription>
                      电子邮箱地址将作为此用户的主要联系方式，请不要与系统中现有的电子邮箱地址重复。
                    </FieldDescription>
                  </Field>
                )}
              />

              <FormField
                control={form.control}
                name="nick_name"
                render={({ field }) => (
                  <Field>
                    <FormLabel>昵称</FormLabel>
                    <FieldContent>
                      <FormItem>
                        <FormControl>
                          <Input {...field} placeholder="请输入昵称" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                    <FieldDescription>
                      用户昵称可以与用户名不同，用于前台显示，如果你将此项留空，将默认使用用户名。
                    </FieldDescription>
                  </Field>
                )}
              />

              <FormField
                control={form.control}
                name="password1"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      密码 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <FormControl>
                          <Input {...field} type="password" placeholder="请输入密码" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                    <FieldDescription>为用户分配一个密码。</FieldDescription>
                  </Field>
                )}
              />

              <FormField
                control={form.control}
                name="password2"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      确认密码 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <FormControl>
                          <Input {...field} type="password" placeholder="请再次输入密码" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                    <FieldDescription>请确认你的密码，与上面输入的密码保持一致。</FieldDescription>
                  </Field>
                )}
              />

              <FormField
                control={form.control}
                name="type"
                render={({ field }) => (
                  <Field>
                    <FormLabel>
                      角色 <span className="text-destructive">*</span>
                    </FormLabel>
                    <FieldContent>
                      <FormItem>
                        <Select
                          onValueChange={(value) => field.onChange(parseInt(value) as 1 | 2)}
                          defaultValue={field.value.toString()}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="请选择角色" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            <SelectItem value="1">管理员</SelectItem>
                            <SelectItem value="2">编辑</SelectItem>
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    </FieldContent>
                    <FieldDescription>
                      管理员具有所有操作权限，编辑仅能包含文章、评论、心情的操作权限。
                    </FieldDescription>
                  </Field>
                )}
              />
              <Field orientation="horizontal">
                <Button type="submit" loading={loading}>
                  保存
                </Button>
              </Field>
            </FieldGroup>
          </form>
        </Form>
      </div>
    </div>
  );
}

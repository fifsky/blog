import { CHeader } from "@/components/CHeader";
import { CFooter } from "@/components/CFooter";
import { useStore } from "@/store/context";
import { useNavigate } from "react-router";
import { LoginRequest } from "@/types/openapi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Field,
  FieldLabel,
  FieldError,
  FieldGroup,
} from "@/components/ui/field";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

export default function Login() {
  // 表单校验规则：用户名与密码为必填
  const formSchema = z.object({
    user_name: z.string().min(1, "请输入用户名"),
    password: z.string().min(1, "请输入密码"),
  });

  // 初始化 react-hook-form，集成 zodResolver 进行校验
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      user_name: "",
      password: "",
    },
    mode: "onChange",
  });

  const loginAction = useStore((s) => s.loginAction);
  const navigate = useNavigate();
  // 提交处理：通过 react-hook-form 的 handleSubmit 获取已校验的数据
  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    await loginAction(data as LoginRequest);
    navigate("/admin/index");
  };
  return (
    <div id="container">
      <CHeader />
      <div className="admin">
        <div className="p-5 border border-[#89d5ef] bg-white">
          <div className="px-[30px]">
            <h2>博客管理登录</h2>
            <form
              method="post"
              onSubmit={form.handleSubmit(onSubmit)}
              className="w-[300px] mx-auto my-[30px]"
            >
              {/* 使用 FieldGroup 包裹所有字段，统一布局与间距 */}
              <FieldGroup>
                {/* 用户名字段 */}
                <Controller
                  name="user_name"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      {/* 使用 FieldLabel 作为 label */}
                      <FieldLabel htmlFor={field.name}>用户名</FieldLabel>
                      {/* 输入框，绑定 RHF 的 field */}
                      <Input
                        {...field}
                        id={field.name}
                        type="text"
                        aria-invalid={fieldState.invalid}
                        autoComplete="username"
                      />
                      {/* 错误展示 */}
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </Field>
                  )}
                />
                {/* 密码字段 */}
                <Controller
                  name="password"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor={field.name}>密码</FieldLabel>
                      <Input
                        {...field}
                        id={field.name}
                        type="password"
                        aria-invalid={fieldState.invalid}
                        autoComplete="current-password"
                      />
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </Field>
                  )}
                />
                <Field>
                  <Button type="submit" size={"sm"}>
                    登 录
                  </Button>
                </Field>
              </FieldGroup>
            </form>
          </div>
        </div>
      </div>
      <CFooter />
    </div>
  );
}

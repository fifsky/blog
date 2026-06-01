import { CHeader } from "@/components/CHeader";
import { CFooter } from "@/components/CFooter";
import { useStore } from "@/store/context";
import { useNavigate } from "react-router";
import { LoginRequest, Setting } from "@/types/openapi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldLabel, FieldError, FieldGroup } from "@/components/ui/field";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { settingApi } from "@/service";
import { useAsyncEffect } from "@/hooks";

export default function Login() {
  const [requireTotp, setRequireTotp] = useState(false);

  // 表单校验规则：用户名与密码为必填
  const formSchema = z.object({
    user_name: z.string().min(1, "请输入用户名"),
    password: z.string().min(1, "请输入密码"),
    totp_code: z.string().optional(),
  });

  // 初始化 react-hook-form，集成 zodResolver 进行校验
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      user_name: "",
      password: "",
      totp_code: "",
    },
    mode: "onChange",
  });
  const [loading, setLoading] = useState(false);
  const [settings, setSettings] = useState<Setting>();
  const loginAction = useStore((s) => s.loginAction);
  const navigate = useNavigate();

  useAsyncEffect(async () => {
    const s = await settingApi();
    setSettings(s);
  }, []);

  // 提交处理：通过 react-hook-form 的 handleSubmit 获取已校验的数据
  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    if (requireTotp && !data.totp_code) {
      dialog.message("请输入2FA验证码");
      return;
    }
    setLoading(true);
    try {
      const res = await loginAction(data as LoginRequest);
      if (res && res.require_totp) {
        setRequireTotp(true);
        setLoading(false);
        return;
      }
      navigate("/admin/index");
    } catch (e) {
      // 错误处理由 request 拦截器完成
    } finally {
      setLoading(false);
    }
  };

  const siteName = settings?.site_name || "無處告別";
  const pageTitle = `登录 - ${siteName}`;
  return (
    <div className="w-[1024px] mt-4 mx-auto min-h-[500px]">
      <title>{pageTitle}</title>
      <meta name="description" content={settings?.site_desc || ""} />
      <meta name="keywords" content={settings?.site_keyword || ""} />
      <CHeader />
      <div className="p-5 border border-[#89d5ef] bg-white relative overflow-hidden">
        <div className="px-[30px] relative z-10">
          <h2 className="border-b border-b-[#cccccc] text-base">博客管理登录</h2>
          <form
            method="post"
            onSubmit={form.handleSubmit(onSubmit)}
            className="w-[300px] mx-auto my-[30px]"
          >
            {/* 使用 FieldGroup 包裹所有字段，统一布局与间距 */}
            <FieldGroup>
              {/* 用户名字段 */}
              {!requireTotp && (
                <>
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
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
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
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </>
              )}
              {requireTotp && (
                <Controller
                  name="totp_code"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor={field.name}>2FA 验证码</FieldLabel>
                      <Input
                        {...field}
                        id={field.name}
                        type="text"
                        maxLength={6}
                        placeholder="请输入 6 位验证码"
                        aria-invalid={fieldState.invalid}
                        autoComplete="one-time-code"
                      />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
              )}
              <Field>
                <Button type="submit" size={"sm"} loading={loading}>
                  登 录
                </Button>
              </Field>
            </FieldGroup>
          </form>
        </div>
      </div>
      <CFooter />
    </div>
  );
}

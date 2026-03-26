import { useState } from "react";
import { useNavigate } from "react-router";
import { CHeader } from "@/components/CHeader";
import { CFooter } from "@/components/CFooter";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldLabel, FieldGroup } from "@/components/ui/field";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { settingApi } from "@/service";
import { aiLoginInitApi } from "@/service/aiLogin";
import { AuthChat } from "@/components/AuthChat";
import { useStore } from "@/store/context";
import { useAsyncEffect } from "@/hooks";

const passwordFormSchema = z.object({
  user_name: z.string().min(1, "请输入用户名"),
  password: z.string().min(1, "请输入密码"),
});

const aiFormSchema = z.object({
  user_name: z.string().min(1, "请输入用户名"),
});

type TabType = "password" | "ai";

export default function Login() {
  const [activeTab, setActiveTab] = useState<TabType>("password");

  const passwordForm = useForm<z.infer<typeof passwordFormSchema>>({
    resolver: zodResolver(passwordFormSchema),
    defaultValues: { user_name: "", password: "" },
    mode: "onChange",
  });

  const aiForm = useForm<z.infer<typeof aiFormSchema>>({
    resolver: zodResolver(aiFormSchema),
    defaultValues: { user_name: "" },
    mode: "onChange",
  });

  const [passwordLoading, setPasswordLoading] = useState(false);
  const [aiLoading, setAILoading] = useState(false);
  const [session, setSession] = useState<{ sessionId: string; welcomeMessage: string } | null>(
    null,
  );
  const [error, setError] = useState("");
  const [settings, setSettings] = useState<any>();

  const loginAction = useStore((s) => s.loginAction);
  const navigate = useNavigate();

  useAsyncEffect(async () => {
    const s = await settingApi();
    setSettings(s);
  }, []);

  const handlePasswordLogin = async (data: z.infer<typeof passwordFormSchema>) => {
    setPasswordLoading(true);
    try {
      await loginAction(data);
      navigate("/admin/index");
    } catch {
      setError("用户名或密码错误");
    } finally {
      setPasswordLoading(false);
    }
  };

  const handleAIInit = async (data: z.infer<typeof aiFormSchema>) => {
    setAILoading(true);
    setError("");
    try {
      const response = await aiLoginInitApi({ user_name: data.user_name });
      setSession({
        sessionId: response.session_id,
        welcomeMessage: response.welcome_message,
      });
    } catch (e: any) {
      setError(e.message || "初始化失败");
    } finally {
      setAILoading(false);
    }
  };

  const handleAISuccess = (token: string) => {
    localStorage.setItem("access_token", token);
    navigate("/admin/index");
  };

  const handleAIFailed = () => {
    setSession(null);
    setError("验证失败，请重试");
  };

  const siteName = settings?.kv?.site_name || "無處告別";

  const tabClass = (tab: TabType) =>
    `px-4 py-2 cursor-pointer border-b-2 transition-colors ${
      activeTab === tab
        ? "border-blue-500 text-blue-600 font-medium"
        : "border-transparent text-gray-500 hover:text-gray-700"
    }`;

  return (
    <div className="w-[1024px] mt-4 mx-auto min-h-[500px]">
      <title>登录 - {siteName}</title>
      <CHeader />
      <div className="p-5 border border-[#89d5ef] bg-white relative overflow-hidden">
        <div className="px-[30px] relative z-10">
          <h2 className="border-b border-b-[#cccccc] text-base">博客管理登录</h2>

          <div className="w-[500px] mx-auto my-[30px]">
            <div className="flex border-b mb-6">
              <div className={tabClass("password")} onClick={() => setActiveTab("password")}>
                密码登录
              </div>
              <div className={tabClass("ai")} onClick={() => setActiveTab("ai")}>
                AI验证登录
              </div>
            </div>

            {activeTab === "password" ? (
              <form onSubmit={passwordForm.handleSubmit(handlePasswordLogin)}>
                <FieldGroup>
                  <Controller
                    name="user_name"
                    control={passwordForm.control}
                    render={({ field, fieldState }) => (
                      <Field data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor={field.name}>用户名</FieldLabel>
                        <Input {...field} id={field.name} type="text" autoComplete="username" />
                        {fieldState.invalid && (
                          <div className="text-red-500 text-sm">{fieldState.error?.message}</div>
                        )}
                      </Field>
                    )}
                  />
                  <Controller
                    name="password"
                    control={passwordForm.control}
                    render={({ field, fieldState }) => (
                      <Field data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor={field.name}>密码</FieldLabel>
                        <Input
                          {...field}
                          id={field.name}
                          type="password"
                          autoComplete="current-password"
                        />
                        {fieldState.invalid && (
                          <div className="text-red-500 text-sm">{fieldState.error?.message}</div>
                        )}
                      </Field>
                    )}
                  />
                  {error && <div className="text-red-500 text-sm">{error}</div>}
                  <Field>
                    <Button type="submit" size="sm" loading={passwordLoading}>
                      登 录
                    </Button>
                  </Field>
                </FieldGroup>
              </form>
            ) : (
              <>
                {!session ? (
                  <form onSubmit={aiForm.handleSubmit(handleAIInit)}>
                    <FieldGroup>
                      <Controller
                        name="user_name"
                        control={aiForm.control}
                        render={({ field, fieldState }) => (
                          <Field data-invalid={fieldState.invalid}>
                            <FieldLabel htmlFor={field.name}>用户名</FieldLabel>
                            <Input
                              {...field}
                              id={field.name}
                              type="text"
                              placeholder="请输入用户名"
                            />
                            {fieldState.invalid && (
                              <div className="text-red-500 text-sm">
                                {fieldState.error?.message}
                              </div>
                            )}
                          </Field>
                        )}
                      />
                      {error && <div className="text-red-500 text-sm">{error}</div>}
                      <Field>
                        <Button type="submit" size="sm" loading={aiLoading}>
                          开始验证
                        </Button>
                      </Field>
                    </FieldGroup>
                  </form>
                ) : (
                  <div>
                    <AuthChat
                      sessionId={session.sessionId}
                      welcomeMessage={session.welcomeMessage}
                      onSuccess={handleAISuccess}
                      onFailed={handleAIFailed}
                    />
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>
      <CFooter />
    </div>
  );
}

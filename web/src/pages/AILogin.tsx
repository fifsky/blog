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
import { aiLoginInitApi } from "@/service/aiLogin";
import { AuthChat } from "@/components/AuthChat";
import { useAsyncEffect } from "@/hooks";
import { settingApi } from "@/service";

const formSchema = z.object({
  user_name: z.string().min(1, "请输入用户名"),
});

export default function AILogin() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { user_name: "" },
    mode: "onChange",
  });

  const [loading, setLoading] = useState(false);
  const [session, setSession] = useState<{ sessionId: string; welcomeMessage: string } | null>(
    null
  );
  const [error, setError] = useState("");
  const [settings, setSettings] = useState<any>();

  const navigate = useNavigate();

  useAsyncEffect(async () => {
    const s = await settingApi();
    setSettings(s);
  }, []);

  const handleInit = async (data: z.infer<typeof formSchema>) => {
    setLoading(true);
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
      setLoading(false);
    }
  };

  const handleSuccess = (token: string) => {
    localStorage.setItem("access_token", token);
    navigate("/admin/index");
  };

  const handleFailed = () => {
    setSession(null);
    setError("验证失败，请重试");
  };

  const siteName = settings?.kv?.site_name || "無處告別";

  return (
    <div className="w-[1024px] mt-4 mx-auto min-h-[500px]">
      <title>AI验证登录 - {siteName}</title>
      <CHeader />
      <div className="p-5 border border-[#89d5ef] bg-white relative overflow-hidden">
        <div className="px-[30px] relative z-10">
          <h2 className="border-b border-b-[#cccccc] text-base">AI身份验证登录</h2>

          <div className="w-[500px] mx-auto my-[30px]">
            {!session ? (
              <form onSubmit={form.handleSubmit(handleInit)}>
                <FieldGroup>
                  <Controller
                    name="user_name"
                    control={form.control}
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
                          <div className="text-red-500 text-sm">{fieldState.error?.message}</div>
                        )}
                      </Field>
                    )}
                  />
                  {error && <div className="text-red-500 text-sm">{error}</div>}
                  <Field>
                    <Button type="submit" size="sm" loading={loading}>
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
                  onSuccess={handleSuccess}
                  onFailed={handleFailed}
                />
              </div>
            )}
          </div>
        </div>
      </div>
      <CFooter />
    </div>
  );
}

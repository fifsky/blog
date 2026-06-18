import { useCallback, useEffect, useMemo, useState } from "react";
import type { ReactNode } from "react";
import {
  ActivityIcon,
  AlertTriangleIcon,
  BotIcon,
  CheckCircle2Icon,
  ClipboardIcon,
  LoaderCircleIcon,
  LogInIcon,
  QrCodeIcon,
  RefreshCwIcon,
  UnplugIcon,
  WifiIcon,
  WifiOffIcon,
} from "lucide-react";
import { QRCodeSVG } from "qrcode.react";
import {
  clawBotCheckLoginApi,
  clawBotDisconnectApi,
  clawBotStartLoginApi,
  clawBotStatusApi,
} from "@/service";
import {
  Alert,
  AlertContent,
  AlertDescription,
  AlertIcon,
  AlertTitle,
} from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import type {
  ClawBotCheckLoginResponse,
  ClawBotLoginSession,
  ClawBotStatusResponse,
} from "@/types/openapi";
import { dialog } from "@/utils/dialog";

const POLL_INTERVAL = 2500;
const FINAL_LOGIN_STATUS = new Set([
  "connected",
  "success",
  "expired",
  "failed",
  "cancelled",
  "canceled",
  "timeout",
]);

const statusText: Record<string, string> = {
  wait: "等待扫码",
  pending: "等待扫码",
  waiting: "等待扫码",
  scaned: "已扫码",
  scanned: "已扫码",
  confirmed: "确认中",
  connected: "已连接",
  success: "已连接",
  expired: "已过期",
  failed: "登录失败",
  cancelled: "已取消",
  canceled: "已取消",
  timeout: "已超时",
  running: "运行中",
  stopped: "已停止",
  idle: "空闲",
  error: "异常",
};

function normalizeStatus(value?: string) {
  return (value || "").trim().toLowerCase();
}

function getStatusText(value?: string) {
  const normalized = normalizeStatus(value);
  return statusText[normalized] || value || "未知";
}

function isFinalLoginStatus(value?: string) {
  const normalized = normalizeStatus(value);
  return normalized ? FINAL_LOGIN_STATUS.has(normalized) : false;
}

function displayValue(value?: string) {
  const text = value?.trim();
  return text || "-";
}

function nextStatusFromCheck(
  current: ClawBotStatusResponse | null,
  result: ClawBotCheckLoginResponse,
): ClawBotStatusResponse | null {
  if (!result.connected) return current;

  return {
    connected: true,
    account: result.account || current?.account,
    monitoring: current?.monitoring || false,
    monitor_status: current?.monitor_status || "",
    last_event_at: current?.last_event_at || "",
    last_error: "",
  };
}

function StatusPanel({
  icon,
  label,
  value,
  description,
  active,
}: {
  icon: ReactNode;
  label: string;
  value: string;
  description?: string;
  active?: boolean;
}) {
  return (
    <div className="flex min-h-24 flex-col gap-3 border-t border-border px-4 py-3 first:rounded-t-md first:border-t-0 last:rounded-b-md md:border-t-0 md:border-l md:first:border-l-0 md:first:rounded-l-md md:first:rounded-tr-none md:last:rounded-r-md md:last:rounded-bl-none">
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span
            className={cn(
              "flex size-7 items-center justify-center rounded-md border",
              active && "bg-muted",
            )}
          >
            {icon}
          </span>
          <span>{label}</span>
        </div>
        <Badge variant={active ? "default" : "secondary"}>{value}</Badge>
      </div>
      {description && <p className="text-xs leading-5 text-muted-foreground">{description}</p>}
    </div>
  );
}

function AccountRow({ label, value }: { label: string; value?: string }) {
  return (
    <div className="grid gap-1 border-b border-border px-4 py-3 last:border-b-0 md:grid-cols-[112px_1fr] md:gap-4">
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className="break-all font-mono text-[13px] leading-5 text-foreground">
        {displayValue(value)}
      </span>
    </div>
  );
}

function StatusSkeleton() {
  return (
    <div className="mt-3 flex flex-col gap-4">
      <div className="rounded-md border border-border">
        <div className="flex items-center justify-between gap-3 border-b border-border px-4 py-3">
          <Skeleton className="h-5 w-32" />
          <Skeleton className="h-8 w-28" />
        </div>
        <div className="grid md:grid-cols-3">
          <Skeleton className="m-4 h-16" />
          <Skeleton className="m-4 h-16" />
          <Skeleton className="m-4 h-16" />
        </div>
      </div>
      <Skeleton className="h-40 w-full" />
    </div>
  );
}

function QrContent({ value, onCopy }: { value: string; onCopy: () => void }) {
  return (
    <div className="flex flex-col gap-3 rounded-md border border-border bg-muted/30 p-5">
      <div className="flex justify-center">
        <div className="rounded bg-white p-3">
          <QRCodeSVG value={value} size={224} level="M" marginSize={2} />
        </div>
      </div>
      <div className="flex flex-wrap items-center justify-center gap-2">
        <Button type="button" variant="outline" size="sm" onClick={onCopy}>
          <ClipboardIcon data-icon="inline-start" />
          复制二维码内容
        </Button>
      </div>
      <details className="rounded-md border border-border bg-background px-3 py-2">
        <summary className="cursor-pointer text-xs font-medium text-muted-foreground">
          查看 QR 原始内容
        </summary>
        <pre className="mt-2 max-h-32 overflow-auto whitespace-pre-wrap break-all font-mono text-xs leading-5">
          {value}
        </pre>
      </details>
    </div>
  );
}

function LoginPanel({
  session,
  checking,
  onCopy,
}: {
  session?: ClawBotLoginSession;
  checking: boolean;
  onCopy: () => void;
}) {
  if (!session) {
    return (
      <div className="flex min-h-52 flex-col items-center justify-center gap-3 rounded-md border border-dashed border-border bg-muted/20 px-6 py-8 text-center">
        <QrCodeIcon className="size-9 text-muted-foreground" />
        <div className="flex flex-col gap-1">
          <p className="text-sm font-medium">等待创建登录会话</p>
          <p className="text-xs text-muted-foreground">点击开始登录后，这里会显示二维码内容。</p>
        </div>
      </div>
    );
  }

  const qrContent = session.qr_content || session.qr_code;

  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Badge variant="outline">{getStatusText(session.status)}</Badge>
          {checking && (
            <span className="inline-flex items-center gap-1 text-xs text-muted-foreground">
              <LoaderCircleIcon className="size-3.5 animate-spin" />
              正在检查
            </span>
          )}
        </div>
        <span className="text-xs text-muted-foreground">
          有效期至 {displayValue(session.expires_at)}
        </span>
      </div>
      {qrContent ? (
        <QrContent value={qrContent} onCopy={onCopy} />
      ) : (
        <Alert variant="warning" appearance="light">
          <AlertIcon>
            <AlertTriangleIcon />
          </AlertIcon>
          <AlertContent>
            <AlertTitle>未返回二维码内容</AlertTitle>
            <AlertDescription>请重新发起登录，或检查 ClawBot 登录服务状态。</AlertDescription>
          </AlertContent>
        </Alert>
      )}
    </div>
  );
}

export default function AdminClawBot() {
  const [status, setStatus] = useState<ClawBotStatusResponse | null>(null);
  const [loginSession, setLoginSession] = useState<ClawBotLoginSession>();
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [starting, setStarting] = useState(false);
  const [checking, setChecking] = useState(false);
  const [disconnecting, setDisconnecting] = useState(false);

  const loadStatus = useCallback(async (silent = false) => {
    if (silent) {
      const data = await clawBotStatusApi(() => undefined);
      setStatus(data);
      if (data.connected) setLoginSession(undefined);
      return data;
    }

    setRefreshing(true);
    try {
      const data = await clawBotStatusApi();
      setStatus(data);
      if (data.connected) setLoginSession(undefined);
      return data;
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  const handleStartLogin = async () => {
    setStarting(true);
    try {
      const session = await clawBotStartLoginApi();
      setLoginSession(session);
      dialog.message("登录会话已创建，请使用微信扫码确认");
    } finally {
      setStarting(false);
    }
  };

  const handleCopyQrContent = async () => {
    const content = loginSession?.qr_content || loginSession?.qr_code || "";
    if (!content) return;

    try {
      await navigator.clipboard.writeText(content);
      dialog.message("二维码内容已复制");
    } catch {
      dialog.message("复制失败，请手动复制");
    }
  };

  const handleDisconnect = () => {
    dialog.confirm("确认要断开微信 ClawBot 连接？", {
      confirmText: "断开",
      onOk: async () => {
        setDisconnecting(true);
        try {
          const data = await clawBotDisconnectApi();
          setStatus(data);
          setLoginSession(undefined);
          dialog.message("已断开微信 ClawBot 连接");
        } finally {
          setDisconnecting(false);
        }
      },
    });
  };

  useEffect(() => {
    void loadStatus();
  }, [loadStatus]);

  useEffect(() => {
    const sessionKey = loginSession?.session_key;
    if (!sessionKey || isFinalLoginStatus(loginSession?.status)) return;

    let timer = 0;
    let cancelled = false;

    // 登录会话未结束时持续轮询，组件卸载或状态收敛后自动停止。
    const poll = async () => {
      let shouldContinue = true;
      setChecking(true);
      try {
        const result = await clawBotCheckLoginApi({ session_key: sessionKey }, () => undefined);
        if (cancelled) return;

        if (result.connected) {
          setStatus((current) => nextStatusFromCheck(current, result));
          setLoginSession(undefined);
          await loadStatus(true);
          dialog.message("微信 ClawBot 已连接");
          shouldContinue = false;
          return;
        }

        const nextSession = result.session || { ...loginSession, status: result.status };
        setLoginSession(nextSession);
        shouldContinue = !isFinalLoginStatus(nextSession.status);
        if (!shouldContinue) {
          dialog.message(`登录会话${getStatusText(nextSession.status)}`);
        }
      } catch {
        shouldContinue = false;
      } finally {
        if (!cancelled) {
          setChecking(false);
          if (shouldContinue) {
            timer = window.setTimeout(poll, POLL_INTERVAL);
          }
        }
      }
    };

    void poll();

    return () => {
      cancelled = true;
      window.clearTimeout(timer);
    };
  }, [loadStatus, loginSession?.session_key, loginSession?.status]);

  const connected = !!status?.connected;
  const account = status?.account;
  const monitorDescription = useMemo(() => {
    const parts = [
      status?.monitor_status ? `状态：${getStatusText(status.monitor_status)}` : "",
      status?.last_event_at ? `最近事件：${status.last_event_at}` : "",
    ].filter(Boolean);
    return parts.join("，") || "暂无监控事件";
  }, [status?.last_event_at, status?.monitor_status]);

  return (
    <div>
      <title>微信 ClawBot - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">微信 ClawBot</h2>
      {loading ? (
        <StatusSkeleton />
      ) : (
        <div className="mt-3 flex flex-col gap-4">
          <section className="overflow-hidden rounded-md border border-border bg-background">
            <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-4 py-3">
              <div className="flex items-center gap-2">
                <BotIcon className="size-4 text-muted-foreground" />
                <div className="flex flex-col gap-0.5">
                  <span className="text-sm font-medium">连接概览</span>
                  <span className="text-xs text-muted-foreground">微信账号连接与消息监控状态</span>
                </div>
              </div>
              <div className="flex flex-wrap items-center gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  loading={refreshing}
                  onClick={() => loadStatus()}
                >
                  <RefreshCwIcon data-icon="inline-start" />
                  刷新
                </Button>
                {connected ? (
                  <Button
                    type="button"
                    variant="destructive"
                    size="sm"
                    loading={disconnecting}
                    onClick={handleDisconnect}
                  >
                    <UnplugIcon data-icon="inline-start" />
                    断开
                  </Button>
                ) : (
                  <Button type="button" size="sm" loading={starting} onClick={handleStartLogin}>
                    <LogInIcon data-icon="inline-start" />
                    开始登录
                  </Button>
                )}
              </div>
            </div>
            <div className="grid md:grid-cols-3">
              <StatusPanel
                icon={
                  connected ? <WifiIcon className="size-4" /> : <WifiOffIcon className="size-4" />
                }
                label="连接状态"
                value={connected ? "已连接" : "未连接"}
                description={connected ? "ClawBot 已保存微信登录凭据。" : "需要扫码完成微信登录。"}
                active={connected}
              />
              <StatusPanel
                icon={<ActivityIcon className="size-4" />}
                label="监控状态"
                value={status?.monitoring ? "监控中" : "未监控"}
                description={monitorDescription}
                active={!!status?.monitoring}
              />
              <StatusPanel
                icon={
                  status?.last_error ? (
                    <AlertTriangleIcon className="size-4" />
                  ) : (
                    <CheckCircle2Icon className="size-4" />
                  )
                }
                label="最近错误"
                value={status?.last_error ? "有异常" : "正常"}
                description={status?.last_error || "暂无错误记录"}
                active={!status?.last_error}
              />
            </div>
          </section>

          <section className="rounded-md border border-border bg-background">
            <div className="flex items-center justify-between gap-3 px-4 py-3">
              <div className="flex flex-col gap-0.5">
                <span className="text-sm font-medium">账号信息</span>
                <span className="text-xs text-muted-foreground">当前保存的微信 ClawBot 账号</span>
              </div>
              <Badge variant={account ? "outline" : "secondary"}>
                {account ? "已保存" : "暂无账号"}
              </Badge>
            </div>
            <Separator />
            {account ? (
              <div>
                <AccountRow label="Account ID" value={account.account_id} />
                <AccountRow label="User ID" value={account.user_id} />
                <AccountRow label="Base URL" value={account.base_url} />
                <AccountRow label="保存时间" value={account.saved_at} />
              </div>
            ) : (
              <div className="px-4 py-4">
                <Alert variant="info" appearance="light">
                  <AlertIcon>
                    <QrCodeIcon />
                  </AlertIcon>
                  <AlertContent>
                    <AlertTitle>尚未连接微信账号</AlertTitle>
                    <AlertDescription>
                      创建登录会话后扫码，连接成功后账号信息会显示在这里。
                    </AlertDescription>
                  </AlertContent>
                </Alert>
              </div>
            )}
          </section>

          {!connected && (
            <section className="rounded-md border border-border bg-background">
              <div className="flex flex-wrap items-center justify-between gap-3 px-4 py-3">
                <div className="flex flex-col gap-0.5">
                  <span className="text-sm font-medium">登录二维码</span>
                  <span className="text-xs text-muted-foreground">
                    使用二维码组件渲染登录内容，并保留原始值用于排查。
                  </span>
                </div>
                {loginSession && (
                  <Badge variant="outline">Session {loginSession.session_key}</Badge>
                )}
              </div>
              <Separator />
              <div className="p-4">
                <LoginPanel
                  session={loginSession}
                  checking={checking}
                  onCopy={handleCopyQrContent}
                />
              </div>
            </section>
          )}
        </div>
      )}
    </div>
  );
}

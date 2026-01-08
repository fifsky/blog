import { type ReactNode } from "react";

interface EmptyProps {
  icon?: ReactNode;
  title?: string;
  content?: string;
  className?: string;
}

export function Empty({
  icon,
  title = "暂无数据",
  content = "当前没有可显示的内容",
  className = "",
}: EmptyProps) {
  return (
    <div className={`flex flex-col items-center justify-center py-16 px-4 ${className}`}>
      {icon && <div className="mb-4 text-foreground bg-muted p-4 rounded-full">{icon}</div>}
      <h3 className="text-lg font-semibold text-foreground mb-2">{title}</h3>
      <p className="text-sm/relaxed text-muted-foreground text-center max-w-md">{content}</p>
    </div>
  );
}

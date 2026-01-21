import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

export function SkeletonArticle({ className }: { className?: string }) {
  return (
    <div className={cn(className)}>
      <div className="flex items-center space-x-4">
        <Skeleton className="h-12 w-12 rounded-full" />
        <div className="space-y-2">
          <Skeleton className="h-4 w-[350px]" />
          <Skeleton className="h-4 w-[200px]" />
        </div>
      </div>
      <Skeleton className="h-4 w-full mt-4" />
      <Skeleton className="h-4 w-full mt-4" />
      {Array.from({ length: 3 }).map((_, j) => (
        <Skeleton key={j} className="h-4 w-full mt-4" />
      ))}
    </div>
  );
}

export function SkeletonArticleList() {
  return (
    <div>
      {Array.from({ length: 5 }).map((_, index) => (
        <div key={index}>
          <SkeletonArticle />
          <div className="border-t border-dashed border-t-[#dbdbdb] mt-5 pt-2.5 pb-2.5 text-right"></div>
        </div>
      ))}
    </div>
  );
}

export function SkeletonArchive() {
  return (
    <div className="h-200">
      <h2 className="text-xl font-bold mb-6">文章归档</h2>
      <div className="space-y-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i}>
            <Skeleton className="h-5 w-32 mb-3" />
            <div className="pl-4 border-l-2 border-gray-200 space-y-2">
              {Array.from({ length: 5 }).map((_, j) => (
                <Skeleton key={j} className="h-4 w-full" />
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

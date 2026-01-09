import { Skeleton } from "@/components/ui/skeleton";

export function SkeletonArticle() {
  return (
    <div>
      <div className="flex items-center space-x-4">
        <Skeleton className="h-12 w-12 rounded-full" />
        <div className="space-y-2">
          <Skeleton className="h-4 w-[350px]" />
          <Skeleton className="h-4 w-[200px]" />
        </div>
      </div>
      <Skeleton className="h-4 w-full mt-4" />
      <Skeleton className="h-4 w-full mt-4" />
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

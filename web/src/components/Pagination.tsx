import { cn } from "@/lib/utils";

interface PaginationProps {
  page: number;
  total: number;
  pageSize?: number;
  onChange: (p: number) => void;
  showPageCount?: number;
}

export function Pagination({
  page,
  total,
  pageSize = 10,
  onChange,
  showPageCount = 5,
}: PaginationProps) {
  const pageTotal = Math.max(1, Math.ceil(total / pageSize));

  const getPageNumbers = () => {
    if (pageTotal <= showPageCount + 2) {
      return Array.from({ length: pageTotal }, (_, i) => i + 1);
    }

    const half = Math.floor(showPageCount / 2);
    let start = Math.max(1, page - half);
    const end = Math.min(pageTotal, start + showPageCount - 1);

    if (end - start + 1 < showPageCount) {
      start = Math.max(1, end - showPageCount + 1);
    }

    const pages: (number | "...")[] = [];

    if (start > 1) {
      pages.push(1, "...");
    }

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    if (end < pageTotal) {
      pages.push("...", pageTotal);
    }

    return pages;
  };

  const pages = getPageNumbers();

  return (
    <ul className="list-none py-[10px] text-center">
      <li className="mx-[0.2em] inline">
        <a
          className={cn(
            "px-2 py-[2px] border border-[#ddd] no-underline select-none outline-none",
            page === 1
              ? "bg-[#f5f5f5] text-[#ccc] cursor-not-allowed"
              : "hover:border-[#06c] hover:bg-[#06c] hover:text-white cursor-pointer",
          )}
          onClick={page === 1 ? undefined : () => onChange(page - 1)}
        >
          上一页
        </a>
      </li>
      {pages.map((p, index) => (
        <li key={index} className="mx-[0.2em] inline">
          {p === "..." ? (
            <span className="px-2 py-[2px] border border-[#ddd] text-[#999]">...</span>
          ) : (
            <a
              className={cn(
                "px-2 py-[2px] border border-[#ddd] no-underline select-none outline-none",
                p === page
                  ? "bg-[#ddd] text-[#555] cursor-default"
                  : "hover:border-[#06c] hover:bg-[#06c] hover:text-white cursor-pointer",
              )}
              onClick={p === page ? undefined : () => onChange(p as number)}
            >
              {p}
            </a>
          )}
        </li>
      ))}
      <li className="mx-[0.2em] inline">
        <a
          className={cn(
            "px-2 py-[2px] border border-[#ddd] no-underline select-none outline-none",
            page === pageTotal
              ? "bg-[#f5f5f5] text-[#ccc] cursor-not-allowed"
              : "hover:border-[#06c] hover:bg-[#06c] hover:text-white cursor-pointer",
          )}
          onClick={page === pageTotal ? undefined : () => onChange(page + 1)}
        >
          下一页
        </a>
      </li>
    </ul>
  );
}

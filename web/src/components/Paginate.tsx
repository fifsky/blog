export function Paginate({
  page,
  pageTotal,
  onChange,
}: {
  page: number;
  pageTotal: number;
  onChange: (p: number) => void;
}) {
  const pages = Array.from({ length: pageTotal }, (_, i) => i + 1);
  return (
    <ul className="list-none py-[10px] text-center">
      {pages.map((p) => (
        <li key={p} className="mx-[0.2em] inline">
          <a
            className={
              "px-2 py-[2px] border no-underline select-none outline-none border-[#ddd]" +
              (p === page
                ? "border-[#ddd] bg-[#ddd] text-[#555] cursor-default"
                : "border-[#ddd] hover:border-[#06c] hover:bg-[#06c] hover:text-white cursor-pointer")
            }
            onClick={p === page ? undefined : () => onChange(p)}
          >
            {p}
          </a>
        </li>
      ))}
    </ul>
  );
}

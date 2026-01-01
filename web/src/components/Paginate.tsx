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
        <ul className="paginator">
            {pages.map((p) => (
                <li key={p} className={p === page ? "active" : ""}>
                    <a onClick={() => onChange(p)}>{p}</a>
                </li>
            ))}
        </ul>
    );
}

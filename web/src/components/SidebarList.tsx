import { useEffect, useState } from "react";
import { Link } from "react-router";
import { archiveApi, cateAllApi, linkAllApi, newCommentApi } from "@/service";

const apiMap = { cateAllApi, newCommentApi, archiveApi, linkAllApi };

export function SidebarList({ title, api }: { title: string; api: keyof typeof apiMap }) {
  const [items, setItems] = useState<{ url: string; content: string }[]>([]);
  useEffect(() => {
    (async () => {
      const data = await apiMap[api]();
      setItems(data.list || []);
    })();
  }, [api]);
  return (
    <div className="mb-5">
      <h2 className="ml-[-10px] mr-0 mb-0 pl-[10px] pr-2 py-1.5 font-normal bg-[url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAMUAAAAlCAMAAADFot1AAAAANlBMVEUAAAD////8/v/6/f73/P7z+/3w+f3r9/vt+Pzn9vvh8/rl9frc8vnU7vfG6fTM6/a55PK/5vPI1xmBAAAAAXRSTlMAQObYZgAAAW1JREFUWMPd1tuugyAQQFEug9ystv//sycVpOCIEk1Pgd2kjw0rzGgJIXZ64l5J06Y57YEbk0yUdWmf9KmlwSeiIMSjmI++I+Qxj+FXQ/KTWhuiBA6i+PKViXG2G70e0bO0Js6uYZBEnENMjsJyXVdMRhvjLh1RdGjXITYOeH9AnDJuC7DiKa1ZRnfEo2tTAhotIdLLKBDkp+nWRL2kLdm9s1m6en50+MsKGw+RF0hdtgu/GSG8F0p/1tlXeAGwie/1nfNjhdThcZSfoBuPIgz4gmKQTlA6QXWd3yvmQZ1dQDmA/RMAK8QgVcEK/H6FTxTHG1D3+b3iIUThBPEqAUEhTt9hUO35VwW0egGpAsoBrD6AU4wAUPsjqEDBk+O3MUFIYTj3gIYmaE/Bq/gXcU/ROMApLGtqAzIKzVjTAKeQrWzwsaLl0weFoqmgSQcZKG32CmJFBxFBO6gTBdAO6kTBaQd1omj9VeEUtIf+AGt9QDY3RRkOAAAAAElFTkSuQmCC')] bg-no-repeat bg-left-top">
        {title}
      </h2>
      <ul className="list-disc list-inside pl-[20px] whitespace-nowrap overflow-hidden text-ellipsis">
        {items.map((v, k) => (
          <li key={k} className="whitespace-nowrap overflow-hidden text-ellipsis">
            {v.url && v.url.startsWith("http") ? (
              <a target="_blank" href={v.url} rel="noreferrer" className="text-ellipsis">
                {v.content}
              </a>
            ) : (
              <Link to={v.url} className="text-ellipsis">
                {v.content}
              </Link>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}

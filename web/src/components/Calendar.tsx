import { useEffect, useState } from "react";
import dayjs from "dayjs";
import { Link } from "react-router";
import { articleCalendarApi } from "@/service";

export function Calendar() {
  const [currDay] = useState(dayjs().format("YYYY-MM-DD"));
  const [currMonth, setCurrMonth] = useState(dayjs().format("YYYY-MM"));
  const [items, setItems] = useState<string[][]>([]);
  const [articleDays, setArticleDays] = useState<number[]>([]);

  useEffect(() => {
    const fetchArticleDays = async () => {
      // 如果查询月份大于当前月份，不请求接口
      if (dayjs(currMonth).isAfter(dayjs(), "month")) {
        setArticleDays([]);
        return;
      }

      const [year, month] = currMonth.split("-").map(Number);
      try {
        const res = await articleCalendarApi({ year, month });
        setArticleDays(res.days || []);
      } catch (e) {
        console.error("Failed to fetch article calendar", e);
      }
    };
    fetchArticleDays();
  }, [currMonth]);

  const builder = () => {
    const days = dayjs(currMonth).daysInMonth();
    const startWeek = dayjs(currMonth + "-01").day();
    const dayItem: (string | number)[] = [];
    for (let i = 0; i < startWeek; i++) dayItem.push("");
    for (let i = 1; i <= days; i++) dayItem.push(i);
    if (dayItem.length < 35) {
      for (let i = 0; i < 35 - dayItem.length; i++) dayItem.push("");
    }
    const result: string[][] = [];
    for (let i = 0; i < dayItem.length; i += 7) result.push(dayItem.slice(i, i + 7).map(String));
    setItems(result);
  };
  const isCurrDay = (d: string) => currMonth + "-" + d === currDay;
  const hasArticle = (d: string) => d && articleDays.includes(Number(d));

  useEffect(builder, [currMonth]);
  return (
    <div className="mb-6">
      <div id="blog-calendar">
        <table className="w-[200px] text-[13px]">
          <tbody>
            <tr>
              <td className="text-center">
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    setCurrMonth(dayjs(currMonth).subtract(1, "year").format("YYYY-MM"));
                  }}
                >
                  &lt;&lt;
                </a>
              </td>
              <td className="text-center">
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    setCurrMonth(dayjs(currMonth).subtract(1, "month").format("YYYY-MM"));
                  }}
                >
                  &lt;
                </a>
              </td>
              <td className="text-center" colSpan={3}>
                {currMonth}
              </td>
              <td className="text-center">
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    setCurrMonth(dayjs(currMonth).add(1, "month").format("YYYY-MM"));
                  }}
                >
                  &gt;
                </a>
              </td>
              <td className="text-center">
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    setCurrMonth(dayjs(currMonth).add(1, "year").format("YYYY-MM"));
                  }}
                >
                  &gt;&gt;
                </a>
              </td>
            </tr>
            <tr>
              <td className="text-center text-[#d08c00]">日</td>
              <td className="text-center">一</td>
              <td className="text-center">二</td>
              <td className="text-center">三</td>
              <td className="text-center">四</td>
              <td className="text-center">五</td>
              <td className="text-center text-[#d08c00]">六</td>
            </tr>
            {items.map((days, index) => (
              <tr key={index}>
                {days.map((d, i) => (
                  <td
                    key={i}
                    className={`text-center ${isCurrDay(d) ? "bg-[#eeeeee] font-bold" : ""}`}
                  >
                    {hasArticle(d) ? (
                      <Link to={`/date/${currMonth.replace("-", "/")}/${d.padStart(2, "0")}`}>
                        {d}
                      </Link>
                    ) : (
                      d
                    )}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

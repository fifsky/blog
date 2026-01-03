import { useEffect, useState } from "react";
import { moodListApi } from "@/service";
import dayjs from "dayjs";

function humanTime(v: string) {
  const currTime = dayjs().add(1, "second");
  const itemTime = dayjs(v);
  if (itemTime.isBetween(currTime.subtract(60, "second"), currTime)) {
    return currTime.diff(itemTime, "second") + "秒前";
  } else if (
    itemTime.isBetween(
      currTime.subtract(60, "minute"),
      currTime.subtract(1, "minute")
    )
  ) {
    return currTime.diff(itemTime, "minute") + "分钟前";
  } else if (
    itemTime.isBetween(currTime.startOf("day"), currTime.endOf("day"))
  ) {
    return "今天" + itemTime.format("HH:mm");
  } else if (
    itemTime.isBetween(
      currTime.subtract(1, "day").startOf("day"),
      currTime.subtract(1, "day").endOf("day")
    )
  ) {
    return "昨天" + itemTime.format("HH:mm");
  } else if (
    itemTime.isBetween(
      currTime.startOf("year"),
      currTime.subtract(1, "day").endOf("day")
    )
  ) {
    return itemTime.format("MM月DD日 HH:mm");
  } else {
    return itemTime.format("YYYY-MM-DD HH:mm");
  }
}

export function Mood() {
  const [moods, setMoods] = useState<any[]>([]);
  const [index, setIndex] = useState(0);
  useEffect(() => {
    (async () => {
      const ret = await moodListApi({ page: 1 });
      setMoods(ret.list || []);
    })();
  }, []);
  const prev = () => setIndex((i) => (i - 1 >= 0 ? i - 1 : i));
  const next = () => setIndex((i) => (i + 1 < moods.length ? i + 1 : i));
  const m = moods[index];
  return (
    <div id="info" className="flex items-start">
      <div id="avatar">
        <img title="莫一哲" alt="莫一哲" src="/assets/images/faceicon.jpg" />
      </div>
      <div id="latest" className="flex-1 min-w-0">
        {m && (
          <p className="current active">
            {m.content}
            <span className="stamp">
              <span className="method">
                {humanTime(m.created_at)} by {m.user.nick_name}
              </span>
            </span>
          </p>
        )}
      </div>
      <div className="handle">
        <i className="iconfont icon-left" title="上一条" onClick={prev} />
        <i className="iconfont icon-right" title="下一条" onClick={next} />
      </div>
    </div>
  );
}

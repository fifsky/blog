import { useEffect, useState } from "react";
import { moodListApi } from "@/service";
import dayjs from "dayjs";
import { MoodItem } from "@/types/openapi";

function humanTime(v: string) {
  const currTime = dayjs().add(1, "second");
  const itemTime = dayjs(v);
  if (itemTime.isBetween(currTime.subtract(60, "second"), currTime)) {
    return currTime.diff(itemTime, "second") + "秒前";
  } else if (itemTime.isBetween(currTime.subtract(60, "minute"), currTime.subtract(1, "minute"))) {
    return currTime.diff(itemTime, "minute") + "分钟前";
  } else if (itemTime.isBetween(currTime.startOf("day"), currTime.endOf("day"))) {
    return "今天" + itemTime.format("HH:mm");
  } else if (
    itemTime.isBetween(
      currTime.subtract(1, "day").startOf("day"),
      currTime.subtract(1, "day").endOf("day"),
    )
  ) {
    return "昨天" + itemTime.format("HH:mm");
  } else if (
    itemTime.isBetween(currTime.startOf("year"), currTime.subtract(1, "day").endOf("day"))
  ) {
    return itemTime.format("MM月DD日 HH:mm");
  } else {
    return itemTime.format("YYYY-MM-DD HH:mm");
  }
}

export function Mood() {
  const [moods, setMoods] = useState<MoodItem[]>([]);
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
    <div className="relative mb-[10px] flex items-start group">
      <div className="p-px border border-[#89d5ef] bg-white overflow-hidden">
        <img
          title="莫一哲"
          alt="莫一哲"
          src="/assets/images/faceicon.jpg"
          className="block w-[96px] h-[96px]"
        />
      </div>
      <div className="flex-1 min-w-0 min-h-[98px] ml-[20px] border border-[#89d5ef] bg-gradient-to-b from-white to-[#eeffde]">
        {/* 左侧箭头 */}
        <div className="absolute top-[0.9rem] left-[110px] w-0 h-0 border-t-[0.6rem] border-t-transparent border-b-[0.6rem] border-b-transparent border-r-[0.7rem] border-r-white"></div>
        {m && (
          <p className="p-[10px] line-[120%] break-all overflow-hidden text-ellipsis text-[#555]">
            {m.content}
            <span className="absolute right-[10px] bottom-[5px] line-[120%] text-[#8c8c8c] text-xs">
              {humanTime(m.created_at)} by {m.user.nick_name}
            </span>
          </p>
        )}
      </div>
      <div className="absolute bottom-0 left-[125px] cursor-pointer user-select-none opacity-0 group-hover:opacity-100 transition-opacity">
        <i
          className="iconfont icon-left text-[20px] text-[rgba(48,175,255,0.5)] hover:text-[rgba(48,175,255,1)] mr-[10px]"
          title="上一条"
          onClick={prev}
        />
        <i
          className="iconfont icon-right text-[20px] text-[rgba(48,175,255,0.5)] hover:text-[rgba(48,175,255,1)]"
          title="下一条"
          onClick={next}
        />
      </div>
    </div>
  );
}

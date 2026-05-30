import { useEffect, useState, useRef } from "react";
import { moodListApi, moodRandomApi } from "@/service";
import dayjs from "dayjs";
import { MoodItem } from "@/types/openapi";
import { Typewriter } from "react-simple-typewriter";
import { ListRestart, ChevronLeft, ChevronRight } from "lucide-react";

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
  const [mood, setMood] = useState<MoodItem | null>(null);
  const [key, setKey] = useState(0);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const playerRef = useRef<any>(null);

  // 初始化时获取最新的一条心情
  useEffect(() => {
    (async () => {
      const ret = await moodListApi({ page: 1 });
      if (ret.list && ret.list.length > 0) {
        setMood(ret.list[0]);
      }
    })();
  }, []);

  // 点击按钮时随机获取一条
  const fetchRandomMood = async () => {
    const ret = await moodRandomApi();
    setMood(ret);
    setKey((k) => k + 1);
  };

  const handlePrev = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (playerRef.current?.aplayer) {
      playerRef.current.aplayer.skipBack();
      playerRef.current.aplayer.play();
    }
  };

  const handleNext = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (playerRef.current?.aplayer) {
      playerRef.current.aplayer.skipForward();
      playerRef.current.aplayer.play();
    }
  };

  return (
    <div className="relative mb-[10px] flex items-start group">
      <div className="relative p-px border border-[#89d5ef] bg-white overflow-hidden w-[98px] h-[98px] group/player [&_meting-js]:block [&_meting-js]:w-full [&_meting-js]:h-full [&_.aplayer]:!w-full [&_.aplayer]:!h-full [&_.aplayer]:!m-0 [&_.aplayer]:!shadow-none [&_.aplayer]:!rounded-none [&_.aplayer-body]:!w-full [&_.aplayer-body]:!h-full [&_.aplayer-pic]:!w-full [&_.aplayer-pic]:!h-full [&_.aplayer-pic]:!rounded-none">
        {/* @ts-expect-error Custom element meting-js */}
        <meting-js ref={playerRef} server="netease" type="playlist" id="18003576240" mini="true" />

        {/* 上一首按钮 */}
        <div
          onClick={handlePrev}
          className="absolute left-1 top-1/2 -translate-y-1/2 w-6 h-6 rounded-full flex items-center justify-center bg-black/30 opacity-0 group-hover/player:opacity-100 transition-opacity cursor-pointer z-10 text-white hover:bg-black/60"
          title="上一首"
        >
          <ChevronLeft size={16} />
        </div>

        {/* 下一首按钮 */}
        <div
          onClick={handleNext}
          className="absolute right-1 top-1/2 -translate-y-1/2 w-6 h-6 rounded-full flex items-center justify-center bg-black/30 opacity-0 group-hover/player:opacity-100 transition-opacity cursor-pointer z-10 text-white hover:bg-black/60"
          title="下一首"
        >
          <ChevronRight size={16} />
        </div>
      </div>
      <div className="flex-1 min-w-0 min-h-[98px] ml-[20px] border border-[#89d5ef] bg-gradient-to-b from-white to-[#eeffde]">
        {/* 左侧箭头 */}
        <div className="absolute top-[0.9rem] left-[110px] w-0 h-0 border-t-[0.6rem] border-t-transparent border-b-[0.6rem] border-b-transparent border-r-[0.7rem] border-r-white"></div>
        {mood && (
          <p className="p-[10px] line-[120%] break-all overflow-hidden text-ellipsis text-[#555]">
            <Typewriter key={key} words={[mood.content]} typeSpeed={50} />
            <span className="absolute right-[10px] bottom-[5px] line-[120%] text-[#8c8c8c] text-xs">
              {humanTime(mood.created_at)} by {mood.user.nick_name}
            </span>
          </p>
        )}
      </div>
      <div
        className="absolute bottom-1 left-[125px] cursor-pointer user-select-none opacity-0 group-hover:opacity-100 transition-opacity"
        title="随机一条"
        onClick={fetchRandomMood}
      >
        <ListRestart className="w-[20px] h-[20px] text-[rgba(48,175,255,0.5)] hover:text-[rgba(48,175,255,1)] transition-colors" />
      </div>
    </div>
  );
}

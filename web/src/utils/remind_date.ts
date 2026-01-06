import dayjs from "dayjs";

export const remindType: Record<number, string> = {
  0: "固定",
  1: "每分钟",
  2: "每小时",
  3: "每天",
  4: "每周",
  5: "每月",
  6: "每年",
};

export const monthFormat: Record<number, string> = {
  1: "01",
  2: "02",
  3: "03",
  4: "04",
  5: "05",
  6: "06",
  7: "07",
  8: "08",
  9: "09",
  10: "10",
  11: "11",
  12: "12",
};

export const weekFormat: Record<number, string> = {
  1: "一",
  2: "二",
  3: "三",
  4: "四",
  5: "五",
  6: "六",
  7: "日",
};

export const numFormat = (n: number) => (n < 10 ? "0" + n : String(n));

export const remindTimeFormat = (v: any) => {
  let str = "";
  switch (v.type) {
    case 0:
      str =
        dayjs(v.created_at).year() +
        "年" +
        monthFormat[v.month] +
        "月" +
        numFormat(v.day) +
        "日 " +
        numFormat(v.hour) +
        "时" +
        numFormat(v.minute) +
        "分";
      break;
    case 3:
      str = numFormat(v.hour) + "时" + numFormat(v.minute) + "分";
      break;
    case 4:
      str =
        "周" +
        weekFormat[v.week] +
        " " +
        numFormat(v.hour) +
        "时" +
        numFormat(v.minute) +
        "分";
      break;
    case 5:
      str =
        numFormat(v.day) +
        "日 " +
        numFormat(v.hour) +
        "时" +
        numFormat(v.minute) +
        "分";
      break;
    case 6:
      str =
        monthFormat[v.month] +
        "月" +
        numFormat(v.day) +
        "日 " +
        numFormat(v.hour) +
        "时" +
        numFormat(v.minute) +
        "分";
      break;
  }
  return str;
};

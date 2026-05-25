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

export const formToCron = (values: any): string => {
  const { type, month = 1, week = 1, day = 1, hour = 0, minute = 0 } = values;
  switch (Number(type)) {
    case 0:
      const currentYear = dayjs().year();
      return `${currentYear}-${numFormat(month)}-${numFormat(day)} ${numFormat(hour)}:${numFormat(minute)}:00`;
    case 1:
      return `* * * * *`;
    case 2:
      return `${minute} * * * *`;
    case 3:
      return `${minute} ${hour} * * *`;
    case 4:
      const cronWeek = week === 7 ? 0 : week;
      return `${minute} ${hour} * * ${cronWeek}`;
    case 5:
      return `${minute} ${hour} ${day} * *`;
    case 6:
      return `${minute} ${hour} ${day} ${month} *`;
    default:
      return "";
  }
};

export const cronToForm = (cron: string): any => {
  const defaultValues = { type: 0, month: 1, week: 1, day: 1, hour: 0, minute: 0 };
  if (!cron) return defaultValues;

  if (cron.includes("-") && cron.includes(":")) {
    const d = dayjs(cron);
    if (d.isValid()) {
      return {
        type: 0,
        month: d.month() + 1,
        day: d.date(),
        hour: d.hour(),
        minute: d.minute(),
        week: 1,
      };
    }
  }

  const parts = cron.split(" ");
  if (parts.length === 5) {
    const [minute, hour, day, month, week] = parts;
    if (cron === "* * * * *") {
      return { ...defaultValues, type: 1 };
    }
    if (hour === "*" && day === "*" && month === "*" && week === "*") {
      return { ...defaultValues, type: 2, minute: Number(minute) };
    }
    if (day === "*" && month === "*" && week === "*") {
      return { ...defaultValues, type: 3, minute: Number(minute), hour: Number(hour) };
    }
    if (day === "*" && month === "*" && week !== "*") {
      const formWeek = Number(week) === 0 ? 7 : Number(week);
      return {
        ...defaultValues,
        type: 4,
        minute: Number(minute),
        hour: Number(hour),
        week: formWeek,
      };
    }
    if (month === "*" && week === "*") {
      return {
        ...defaultValues,
        type: 5,
        minute: Number(minute),
        hour: Number(hour),
        day: Number(day),
      };
    }
    if (week === "*") {
      return {
        ...defaultValues,
        type: 6,
        minute: Number(minute),
        hour: Number(hour),
        day: Number(day),
        month: Number(month),
      };
    }
  }

  return defaultValues;
};

export const remindTimeFormat = (record: any) => {
  if (!record.cron) return "";
  const v = cronToForm(record.cron);
  let str = "";
  switch (v.type) {
    case 0:
      str =
        dayjs(record.cron).year() +
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
    case 1:
      str = "每分钟";
      break;
    case 2:
      str = numFormat(v.minute) + "分";
      break;
    case 3:
      str = numFormat(v.hour) + "时" + numFormat(v.minute) + "分";
      break;
    case 4:
      str = "周" + weekFormat[v.week] + " " + numFormat(v.hour) + "时" + numFormat(v.minute) + "分";
      break;
    case 5:
      str = numFormat(v.day) + "日 " + numFormat(v.hour) + "时" + numFormat(v.minute) + "分";
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

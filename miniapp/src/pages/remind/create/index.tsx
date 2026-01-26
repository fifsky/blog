import { useState } from "react";
import { Picker, Text, View } from "@tarojs/components";
import Taro from "@tarojs/taro";
import { AtButton, AtCard, AtTextarea } from "taro-ui";
import { remindCreateApi } from "../../../service";

const TYPES = [
  { value: 0, label: "固定时间" },
  { value: 3, label: "每周" },
  { value: 4, label: "每天" },
  { value: 5, label: "每月" },
  { value: 6, label: "每年" },
];

export default function RemindCreatePage() {
  const [typeIndex, setTypeIndex] = useState(0);
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const backToList = () => {
    const pages = Taro.getCurrentPages();
    if (pages.length > 1) {
      Taro.navigateBack();
      return;
    }
    Taro.switchTab({ url: "/pages/remind/index" });
  };

  const submit = async () => {
    if (!content.trim()) {
      Taro.showToast({ title: "请输入提醒内容", icon: "none" });
      return;
    }
    if (submitting) return;

    setSubmitting(true);
    try {
      const type = TYPES[typeIndex].value;
      await remindCreateApi({ type, content });
      Taro.showToast({ title: "创建成功", icon: "success" });
      backToList();
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "创建失败", icon: "none" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", padding: "24rpx" }}>
      <AtCard title="新增提醒" isFull>
        <Picker
          mode="selector"
          range={TYPES}
          rangeKey="label"
          onChange={(e: any) => setTypeIndex(Number(e.detail.value))}
        >
          <View
            style={{
              padding: "20rpx",
              borderRadius: "16rpx",
              backgroundColor: "#fff",
              border: "1rpx solid #eee",
            }}
          >
            <Text style={{ fontSize: "26rpx", color: "#333" }}>{TYPES[typeIndex].label}</Text>
          </View>
        </Picker>

        <View style={{ marginTop: "24rpx" }}>
          <AtTextarea
            value={content}
            onChange={(v) => setContent(String(v))}
            maxLength={200}
            placeholder="提醒内容，比如：每周一早上 9 点喝水"
            height={220}
          />
        </View>

        <View style={{ marginTop: "24rpx" }}>
          <AtButton type="primary" loading={submitting} disabled={submitting} onClick={submit}>
            创建提醒
          </AtButton>
        </View>
      </AtCard>
    </View>
  );
}


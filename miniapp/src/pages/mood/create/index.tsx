import { useState } from "react";
import { View } from "@tarojs/components";
import Taro from "@tarojs/taro";
import { AtButton, AtCard, AtTextarea } from "taro-ui";
import { moodCreateApi } from "../../../service";

export default function MoodCreatePage() {
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const backToList = () => {
    const pages = Taro.getCurrentPages();
    if (pages.length > 1) {
      Taro.navigateBack();
      return;
    }
    Taro.switchTab({ url: "/pages/mood/index" });
  };

  const submit = async () => {
    if (!content.trim()) {
      Taro.showToast({ title: "请输入心情内容", icon: "none" });
      return;
    }
    if (submitting) return;

    setSubmitting(true);
    try {
      await moodCreateApi({ content });
      Taro.showToast({ title: "发布成功", icon: "success" });
      backToList();
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "发布失败", icon: "none" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", padding: "24rpx" }}>
      <AtCard title="写下此刻的心情">
        <AtTextarea
          value={content}
          onChange={(v) => setContent(String(v))}
          maxLength={500}
          placeholder="写点什么…"
          height={220}
        />
        <View style={{ marginTop: "24rpx" }}>
          <AtButton type="primary" loading={submitting} disabled={submitting} onClick={submit}>
            发布
          </AtButton>
        </View>
      </AtCard>
    </View>
  );
}

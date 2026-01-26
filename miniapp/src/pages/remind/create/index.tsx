import { useState } from "react";
import { View } from "@tarojs/components";
import Taro from "@tarojs/taro";
import { AtButton, AtCard, AtTextarea } from "taro-ui";
import { remindSmartCreateApi } from "../../../service";

export default function RemindCreatePage() {
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
      await remindSmartCreateApi(content);
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
      <AtCard title="新增提醒">
        <View style={{ marginTop: "12rpx" }}>
          <AtTextarea
            value={content}
            onChange={(v) => setContent(String(v))}
            maxLength={200}
            placeholder="例如：每周一早上 9 点喝水；明天下午 3 点开会"
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

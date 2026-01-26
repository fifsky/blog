import { useEffect, useState } from "react";
import { ScrollView, Text, View } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import { AtFab, AtIcon } from "taro-ui";
import type { MoodItem } from "../../types/openapi";
import { moodListApi } from "../../service";

export default function MoodPage() {
  const [list, setList] = useState<MoodItem[]>([]);

  const load = async () => {
    try {
      const resp = await moodListApi({ page: 1 });
      setList(resp.list || []);
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "加载失败", icon: "none" });
    }
  };

  useDidShow(() => {
    void load();
  });

  useEffect(() => {
    void load();
  }, []);

  const goCreate = () => {
    Taro.navigateTo({ url: "/pages/mood/create/index" });
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", paddingBottom: "120rpx" }}>
      <ScrollView scrollY style={{ height: "100vh" }}>
        {list.map((m) => (
          <View
            key={m.id}
            style={{
              padding: "24rpx",
              margin: "24rpx",
              backgroundColor: "#fff",
              borderRadius: "12rpx",
            }}
          >
            <View>
              <Text style={{ fontSize: "24rpx", color: "#999" }}>{m.created_at}</Text>
            </View>
            <View>
              <Text style={{ fontSize: "28rpx", color: "#333" }}>{m.content}</Text>
            </View>
          </View>
        ))}
        {list.length === 0 ? (
          <View style={{ padding: "80rpx 24rpx", textAlign: "center" }}>
            <Text style={{ fontSize: "28rpx", color: "#999" }}>暂无心情</Text>
          </View>
        ) : null}
      </ScrollView>

      <View
        style={{
          position: "fixed",
          right: "32rpx",
          bottom: "140rpx",
          zIndex: 1000,
        }}
      >
        <AtFab onClick={goCreate}>
          <AtIcon value="add" size={24} />
        </AtFab>
      </View>
    </View>
  );
}

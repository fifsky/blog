import { useEffect, useState } from "react";
import { ScrollView, Text, View } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import { AtButton, AtCard, AtFab, AtIcon } from "taro-ui";
import type { RemindItem } from "../../types/openapi";
import { remindDeleteApi, remindListApi } from "../../service";

export default function RemindPage() {
  const [list, setList] = useState<RemindItem[]>([]);

  const load = async () => {
    try {
      const resp = await remindListApi({ page: 1 });
      setList(resp.list || []);
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "加载失败", icon: "none" });
    }
  };

  useEffect(() => {
    void load();
  }, []);

  useDidShow(() => {
    void load();
  });

  const goCreate = () => {
    Taro.navigateTo({ url: "/pages/remind/create/index" });
  };

  const remove = async (id: number) => {
    try {
      await remindDeleteApi({ id });
      Taro.showToast({ title: "删除成功", icon: "success" });
      await load();
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "删除失败", icon: "none" });
    }
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", paddingBottom: "120rpx" }}>
      <ScrollView scrollY style={{ height: "100vh" }}>
        {list.map((r) => (
          <View key={r.id} style={{ padding: "0 24rpx 24rpx" }}>
            <AtCard title={r.content} note={`下次：${r.next_time} · 类型：${r.type}`} isFull>
              <View style={{ display: "flex", justifyContent: "flex-end" }}>
                <AtButton
                  size="small"
                  onClick={() => {
                    void remove(r.id);
                  }}
                >
                  删除
                </AtButton>
              </View>
            </AtCard>
          </View>
        ))}
        {list.length === 0 ? (
          <View style={{ padding: "80rpx 24rpx", textAlign: "center" }}>
            <Text style={{ fontSize: "28rpx", color: "#999" }}>暂无提醒</Text>
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

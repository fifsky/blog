import { useEffect, useState } from "react";
import { ScrollView, Text, View } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import { AtActionSheet, AtActionSheetItem, AtFab, AtIcon } from "taro-ui";
import type { RemindItem } from "../../types/openapi";
import { remindDeleteApi, remindListApi } from "../../service";

export default function RemindPage() {
  const [list, setList] = useState<RemindItem[]>([]);
  const [actionOpen, setActionOpen] = useState(false);
  const [activeId, setActiveId] = useState<number | null>(null);

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

  const openActions = (id: number) => {
    setActiveId(id);
    setActionOpen(true);
  };

  const closeActions = () => {
    setActionOpen(false);
  };

  const handleEdit = () => {
    if (!activeId) return;
    setActionOpen(false);
    Taro.navigateTo({ url: `/pages/remind/create/index?id=${activeId}` });
  };

  const handleDelete = async () => {
    if (!activeId) return;
    const res = await Taro.showModal({
      title: "确认删除",
      content: "删除后将无法恢复，是否确认删除该提醒？",
      confirmText: "删除",
      cancelText: "取消",
    });
    if (!res.confirm) {
      setActionOpen(false);
      return;
    }
    setActionOpen(false);
    await remove(activeId);
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", paddingBottom: "120rpx" }}>
      <ScrollView scrollY style={{ height: "100vh" }}>
        {list.map((r) => (
          <View
            key={r.id}
            style={{
              padding: "24rpx",
              margin: "24rpx",
              backgroundColor: "#fff",
              borderRadius: "12rpx",
            }}
          >
            <View
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
              }}
            >
              <View style={{ display: "flex", alignItems: "center" }}>
                <AtIcon value="clock" size={12} color="#999" />
                <Text style={{ fontSize: "24rpx", color: "#999", paddingLeft: "8rpx" }}>
                  {r.next_time}
                </Text>
              </View>
              <View
                style={{
                  padding: "8rpx 0",
                  display: "flex",
                  flexDirection: "row",
                  alignItems: "center",
                }}
                onClick={() => {
                  openActions(r.id);
                }}
              >
                <View
                  style={{
                    width: "6rpx",
                    height: "6rpx",
                    borderRadius: "9999rpx",
                    backgroundColor: "#999",
                    marginLeft: "4rpx",
                  }}
                />
                <View
                  style={{
                    width: "6rpx",
                    height: "6rpx",
                    borderRadius: "9999rpx",
                    backgroundColor: "#999",
                    marginLeft: "4rpx",
                  }}
                />
                <View
                  style={{
                    width: "6rpx",
                    height: "6rpx",
                    borderRadius: "9999rpx",
                    backgroundColor: "#999",
                    marginLeft: "4rpx",
                  }}
                />
              </View>
            </View>
            <View style={{ paddingTop: "12rpx" }}>
              <Text style={{ fontSize: "28rpx", color: "#333" }}>{r.content}</Text>
            </View>
            <View style={{ width: "100%", display: "flex", justifyContent: "flex-end" }}>
              <Text style={{ fontSize: "24rpx", color: "#999" }}>{r.created_at}</Text>
            </View>
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

      <AtActionSheet isOpened={actionOpen} cancelText="取消" onClose={closeActions}>
        <AtActionSheetItem
          onClick={() => {
            handleEdit();
          }}
        >
          编辑
        </AtActionSheetItem>
        <AtActionSheetItem
          onClick={() => {
            void handleDelete();
          }}
        >
          删除
        </AtActionSheetItem>
      </AtActionSheet>
    </View>
  );
}

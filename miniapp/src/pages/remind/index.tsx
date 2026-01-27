import { useEffect, useState } from "react";
import { ScrollView, Text, View } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import { ActionSheet, Button, Cell, Field, Popup, Textarea } from "@taroify/core";
import { Ellipsis } from "@taroify/icons";
import type { RemindItem } from "../../types/openapi";
import { remindDeleteApi, remindListApi } from "../../service";

export default function RemindPage() {
  const [list, setList] = useState<RemindItem[]>([]);
  const [actionOpen, setActionOpen] = useState(false);
  const [activeId, setActiveId] = useState<number | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);

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

  const openCreate = () => {
    setCreateOpen(true);
  };

  const closeCreate = () => {
    if (submitting) return;
    setCreateOpen(false);
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
    // 这里暂不支持编辑已有提醒，仅预留逻辑
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

  const submitCreate = async () => {
    if (!content.trim()) {
      Taro.showToast({ title: "请输入提醒内容", icon: "none" });
      return;
    }
    if (submitting) return;

    setSubmitting(true);
    try {
      await (await import("../../service")).remindSmartCreateApi(content);
      Taro.showToast({ title: "创建成功", icon: "success" });
      setContent("");
      setCreateOpen(false);
      await load();
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "创建失败", icon: "none" });
    } finally {
      setSubmitting(false);
    }
  };

  const handleActionSelect = (e: any) => {
    const value = e?.detail?.value ?? e?.value ?? e?.action?.value;
    if (value === "edit") {
      handleEdit();
      return;
    }
    if (value === "delete") {
      void handleDelete();
      return;
    }
    setActionOpen(false);
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f2f4f6", paddingBottom: "120rpx" }}>
      <ScrollView scrollY style={{ height: "100vh" }}>
        <View style={{ padding: "24rpx" }}>
          <Button
            color="primary"
            shape="round"
            style={{ width: "100%", height: "88rpx", fontSize: "28rpx" } as any}
            onClick={openCreate}
          >
            + 新增提醒
          </Button>
        </View>
        {list.map((r) => (
          <View
            key={r.id}
            style={{
              padding: "24rpx",
              margin: "24rpx",
              backgroundColor: "#fff",
              borderRadius: "16rpx",
              display: "flex",
              alignItems: "center",
            }}
          >
            {(() => {
              const raw = String(r.next_time || "");
              const parts = raw.split(" ");
              const date = parts[0] || "";
              const timePart = (parts[1] || "").slice(0, 5);
              const timeText = timePart || "--:--";
              return (
                <>
                  <View
                    style={{
                      width: "96rpx",
                      height: "96rpx",
                      borderRadius: "9999rpx",
                      backgroundColor: "#ff8a65",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                      flexShrink: 0,
                    }}
                  >
                    <Text style={{ fontSize: "26rpx", color: "#fff", fontWeight: 600 }}>
                      {timeText}
                    </Text>
                  </View>

                  <View style={{ flex: 1, paddingLeft: "24rpx" }}>
                    <View
                      style={{
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                      }}
                    >
                      <Text
                        style={{
                          fontSize: "28rpx",
                          color: "#333",
                          wordBreak: "break-all",
                          overflowWrap: "break-word",
                          paddingRight: "16rpx",
                          flex: 1,
                        }}
                      >
                        {r.content}
                      </Text>

                      <View
                        style={{
                          padding: "8rpx 0",
                          display: "flex",
                          alignItems: "center",
                          flexShrink: 0,
                        }}
                        onClick={() => {
                          openActions(r.id);
                        }}
                      >
                        <Ellipsis style={{ color: "#999", fontSize: "36rpx" } as any} />
                      </View>
                    </View>

                    <View style={{ marginTop: "10rpx" }}>
                      <Text style={{ fontSize: "24rpx", color: "#999" }}>{date}</Text>
                    </View>
                  </View>
                </>
              );
            })()}
          </View>
        ))}
        {list.length === 0 ? (
          <View style={{ padding: "80rpx 24rpx", textAlign: "center" }}>
            <Text style={{ fontSize: "28rpx", color: "#999" }}>暂无提醒</Text>
          </View>
        ) : null}
      </ScrollView>
      <ActionSheet
        open={actionOpen}
        onClose={setActionOpen}
        onCancel={closeActions}
        cancelText="取消"
        actions={[
          { name: "编辑", value: "edit" },
          { name: "删除", value: "delete", style: { color: "#ee0a24" } },
        ]}
        onSelect={handleActionSelect}
      />

      <Popup open={createOpen} rounded placement="bottom" onClose={setCreateOpen}>
        <View style={{ padding: "32rpx 24rpx 40rpx" }}>
          <Text style={{ fontSize: "32rpx", fontWeight: 600, color: "#333" }}>新增提醒</Text>
          <View style={{ marginTop: "24rpx" }}>
            <Cell.Group inset>
              <Field align="start">
                <Textarea
                  value={content}
                  onChange={(v: any) => setContent(String(v?.detail?.value ?? v))}
                  limit={200}
                  autoHeight
                  autoFocus
                  placeholder="输入提醒的内容，例如：每周一早上 9 点喝水"
                />
              </Field>
            </Cell.Group>
          </View>
          <View
            style={{
              marginTop: "32rpx",
              display: "flex",
              justifyContent: "space-between",
              gap: "24rpx",
            }}
          >
            <Button
              size="small"
              shape="round"
              onClick={closeCreate}
              style={{ flex: 1, backgroundColor: "#f5f5f5", border: "none" } as any}
            >
              取消
            </Button>
            <Button
              size="small"
              shape="round"
              color="primary"
              loading={submitting}
              disabled={submitting}
              onClick={submitCreate}
              style={{ flex: 1 } as any}
            >
              保存
            </Button>
          </View>
        </View>
      </Popup>
    </View>
  );
}

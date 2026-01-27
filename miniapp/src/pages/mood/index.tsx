import { useEffect, useState } from "react";
import { ScrollView, Text, View } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import { Button, Cell, Field, Popup, Textarea } from "@taroify/core";
import type { MoodItem } from "../../types/openapi";
import { moodCreateApi, moodListApi } from "../../service";

export default function MoodPage() {
  const [list, setList] = useState<MoodItem[]>([]);
  const [createOpen, setCreateOpen] = useState(false);
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);

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

  const openCreate = () => {
    setCreateOpen(true);
  };

  const closeCreate = () => {
    if (submitting) return;
    setCreateOpen(false);
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
      setContent("");
      setCreateOpen(false);
      await load();
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "发布失败", icon: "none" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f2f4f6", paddingBottom: "120rpx" }}>
      <ScrollView scrollY style={{ height: "100vh" }}>
        <View
          style={{
            padding: "24rpx",
          }}
        >
          <Button
            color="primary"
            shape="round"
            style={{ width: "100%", height: "88rpx", fontSize: "28rpx" } as any}
            onClick={openCreate}
          >
            + 写心情
          </Button>
        </View>
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
      <Popup open={createOpen} rounded placement="bottom" onClose={setCreateOpen}>
        <View style={{ padding: "32rpx 24rpx 40rpx" }}>
          <Text style={{ fontSize: "32rpx", fontWeight: 600, color: "#333" }}>写下此刻的心情</Text>
          <View style={{ marginTop: "24rpx" }}>
            <Cell.Group inset>
              <Field align="start">
                <Textarea
                  value={content}
                  onChange={(v: any) => setContent(String(v?.detail?.value ?? v))}
                  limit={500}
                  autoHeight
                  placeholder="写点什么…"
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
              color="primary"
              size="small"
              shape="round"
              loading={submitting}
              disabled={submitting}
              onClick={submit}
              style={{ flex: 1 } as any}
            >
              发布
            </Button>
          </View>
        </View>
      </Popup>
    </View>
  );
}

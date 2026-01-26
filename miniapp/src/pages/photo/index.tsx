import { useEffect, useState } from "react";
import { Image, ScrollView, Text, View } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import { AtFab, AtIcon } from "taro-ui";
import type { PhotoItem } from "../../types/openapi";
import { photoListApi } from "../../service";

export default function PhotoPage() {
  const [photos, setPhotos] = useState<PhotoItem[]>([]);

  const loadPhotos = async () => {
    try {
      const resp = await photoListApi({ page: 1 });
      setPhotos(resp.list || []);
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "加载相册失败", icon: "none" });
    }
  };

  useEffect(() => {
    void loadPhotos();
  }, []);

  useDidShow(() => {
    void loadPhotos();
  });

  const goCreate = () => {
    Taro.navigateTo({ url: "/pages/photo/create/index" });
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", paddingBottom: "120rpx" }}>
      <ScrollView scrollY style={{ height: "100vh" }}>
        {photos.map((p) => (
          <View
            key={p.id}
            style={{
              padding: "24rpx",
              margin: "24rpx",
              backgroundColor: "#fff",
              borderRadius: "12rpx",
            }}
          >
            <View>
              <Image
                mode="aspectFill"
                style={{ width: "100%", height: "360rpx" }}
                src={p.thumbnail || p.src}
              />
            </View>
            <View>
              <Text
                style={{ fontSize: "24rpx", color: "#999" }}
              >{`${p.province_name} · ${p.city_name} · ${p.created_at}`}</Text>
            </View>
          </View>
        ))}
        {photos.length === 0 ? (
          <View style={{ padding: "80rpx 24rpx", textAlign: "center" }}>
            <Text style={{ fontSize: "28rpx", color: "#999" }}>暂无照片，先上传一张吧</Text>
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

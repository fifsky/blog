import { useEffect, useState } from "react";
import { Image, Picker, ScrollView, Text, View } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import { Button, Popup } from "@taroify/core";
import type { PhotoItem } from "../../types/openapi";
import type { RegionItem } from "../../types/openapi";
import {
  nearestRegionApi,
  ossPresignApi,
  photoCreateApi,
  photoListApi,
  regionListApi,
} from "../../service";
import { putFileToPresignUrl } from "../../utils/upload";

export default function PhotoPage() {
  const [photos, setPhotos] = useState<PhotoItem[]>([]);
  const [provinceOptions, setProvinceOptions] = useState<RegionItem[]>([]);
  const [cityOptions, setCityOptions] = useState<RegionItem[]>([]);
  const [provinceIndex, setProvinceIndex] = useState(0);
  const [cityIndex, setCityIndex] = useState(0);
  const [uploading, setUploading] = useState(false);
  const [autoLocated, setAutoLocated] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);

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
    void loadRegions();
  }, []);

  useDidShow(() => {
    void loadPhotos();
  });

  const loadRegions = async () => {
    try {
      const provinces = await regionListApi({ parent_id: 0 });
      setProvinceOptions(provinces.list || []);
      if (provinces.list.length > 0) {
        const first = provinces.list[0];
        const cities = await regionListApi({ parent_id: first.region_id });
        setCityOptions(cities.list || []);
      }
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "加载地区失败", icon: "none" });
    }
  };

  useEffect(() => {
    if (!autoLocated && provinceOptions.length > 0) {
      setAutoLocated(true);
      void autoDetectRegion();
    }
  }, [autoLocated, provinceOptions]);

  const onRegionChange = (e: any) => {
    const value: number[] = e.detail.value || [];
    const pIndex = value[0] ?? 0;
    const cIndex = value[1] ?? 0;
    setProvinceIndex(pIndex);
    setCityIndex(cIndex);
  };

  const onRegionColumnChange = async (e: any) => {
    const column: number = e.detail.column;
    const value: number = e.detail.value;

    if (column === 0) {
      const pIndex = value;
      setProvinceIndex(pIndex);
      const province = provinceOptions[pIndex];
      if (!province) return;
      const cities = await regionListApi({ parent_id: province.region_id });
      setCityOptions(cities.list || []);
      setCityIndex(0);
    }

    if (column === 1) {
      setCityIndex(value);
    }
  };

  const autoDetectRegion = async () => {
    try {
      const location = await Taro.getLocation({ type: "gcj02" });
      const nearest = await nearestRegionApi(location.latitude, location.longitude);
      const pIndex = provinceOptions.findIndex((p) => p.region_id === nearest.province_id);
      if (pIndex >= 0) {
        setProvinceIndex(pIndex);
        const cities = await regionListApi({ parent_id: nearest.province_id });
        setCityOptions(cities.list || []);
        const cIndex = cities.list.findIndex((c) => c.region_id === nearest.city_id);
        if (cIndex >= 0) {
          setCityIndex(cIndex);
        }
      }
      Taro.showToast({ title: "已自动定位城市", icon: "success" });
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "自动定位失败，请手动选择", icon: "none" });
    }
  };

  const openCreate = () => {
    setCreateOpen(true);
  };

  const closeCreate = () => {
    if (uploading) return;
    setCreateOpen(false);
  };

  const doUpload = async (sourceType: ("camera" | "album")[]) => {
    if (uploading) return;

    const province = provinceOptions[provinceIndex];
    const city = cityOptions[cityIndex];
    if (!province || !city) {
      Taro.showToast({ title: "请选择省市", icon: "none" });
      return;
    }

    setUploading(true);
    try {
      const choose = await Taro.chooseImage({ count: 1, sourceType });
      if (choose.tempFilePaths.length === 0) {
        return;
      }
      const filePath = choose.tempFilePaths[0];
      const fileName = filePath.split("/").pop() || "photo.jpg";

      const presign = await ossPresignApi({ filename: fileName });
      await putFileToPresignUrl(presign.url, filePath);

      await photoCreateApi({
        title: city.region_name || "打卡照片",
        description: `${province.region_name} ${city.region_name}`,
        srcs: [presign.cdn_url],
        province: province.region_id,
        city: city.region_id,
      });

      Taro.showToast({ title: "上传成功", icon: "success" });
      setCreateOpen(false);
      await loadPhotos();
    } catch (e: any) {
      const msg = String(e?.errMsg || e?.message || "");
      if (msg.includes("cancel")) {
        return;
      }
      Taro.showToast({ title: e?.message || "上传失败", icon: "none" });
    } finally {
      setUploading(false);
    }
  };

  const provinceName = provinceOptions[provinceIndex]?.region_name || "选择省份";
  const cityName = cityOptions[cityIndex]?.region_name || "选择城市";

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
            + 上传相册
          </Button>
        </View>
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

      <Popup open={createOpen} rounded placement="bottom" onClose={setCreateOpen}>
        <View style={{ padding: "32rpx 24rpx 40rpx" }}>
          <Text style={{ fontSize: "32rpx", fontWeight: 600, color: "#333" }}>上传相册</Text>
          <View style={{ marginTop: "24rpx" }}>
            <View
              style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}
            >
              <Picker
                mode="multiSelector"
                range={[
                  provinceOptions.map((p) => p.region_name),
                  cityOptions.map((c) => c.region_name),
                ]}
                value={[provinceIndex, cityIndex]}
                onChange={onRegionChange}
                onColumnChange={onRegionColumnChange}
              >
                <View
                  style={{
                    flex: 1,
                    padding: "20rpx 0",
                  }}
                >
                  <Text style={{ fontSize: "28rpx", color: "#333" }}>
                    <Text style={{ fontWeight: "bold" }}>打卡城市</Text>{" "}
                    {provinceName &&
                    cityName &&
                    provinceName !== "选择省份" &&
                    cityName !== "选择城市"
                      ? `${provinceName} ${cityName}`
                      : "请选择"}
                  </Text>
                </View>
              </Picker>

              <Button variant="text" onClick={autoDetectRegion}>
                📍 定位
              </Button>
            </View>
          </View>

          <View style={{ marginTop: "24rpx" }}>
            <View
              style={{
                display: "flex",
                justifyContent: "space-between",
                gap: "24rpx",
              }}
            >
              <Button
                color="primary"
                size="small"
                shape="round"
                loading={uploading}
                disabled={uploading}
                onClick={() => void doUpload(["camera"])}
                style={{ flex: 1 } as any}
              >
                拍照
              </Button>
              <Button
                color="primary"
                size="small"
                shape="round"
                loading={uploading}
                disabled={uploading}
                onClick={() => void doUpload(["album"])}
                style={{ flex: 1 } as any}
              >
                上传
              </Button>
            </View>
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
          </View>
        </View>
      </Popup>
    </View>
  );
}

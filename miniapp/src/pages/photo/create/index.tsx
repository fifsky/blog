import { useEffect, useState } from "react";
import { Picker, Text, View } from "@tarojs/components";
import Taro from "@tarojs/taro";
import { AtButton, AtCard } from "taro-ui";
import type { RegionItem } from "../../../types/openapi";
import { nearestRegionApi, ossPresignApi, photoCreateApi, regionListApi } from "../../../service";
import { putFileToPresignUrl } from "../../../utils/upload";

export default function PhotoCreatePage() {
  const [provinceOptions, setProvinceOptions] = useState<RegionItem[]>([]);
  const [cityOptions, setCityOptions] = useState<RegionItem[]>([]);
  const [provinceIndex, setProvinceIndex] = useState(0);
  const [cityIndex, setCityIndex] = useState(0);
  const [uploading, setUploading] = useState(false);

  const backToList = () => {
    const pages = Taro.getCurrentPages();
    if (pages.length > 1) {
      Taro.navigateBack();
      return;
    }
    Taro.switchTab({ url: "/pages/photo/index" });
  };

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
    void loadRegions();
  }, []);

  const onProvinceChange = async (e: any) => {
    const index = Number(e.detail.value);
    setProvinceIndex(index);
    const province = provinceOptions[index];
    if (!province) return;
    const cities = await regionListApi({ parent_id: province.region_id });
    setCityOptions(cities.list || []);
    setCityIndex(0);
  };

  const onCityChange = (e: any) => {
    setCityIndex(Number(e.detail.value));
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

  const upload = async () => {
    if (uploading) return;

    const province = provinceOptions[provinceIndex];
    const city = cityOptions[cityIndex];
    if (!province || !city) {
      Taro.showToast({ title: "请选择省市", icon: "none" });
      return;
    }

    setUploading(true);
    try {
      const choose = await Taro.chooseImage({ count: 1 });
      const filePath = choose.tempFilePaths[0];
      const fileName = filePath.split("/").pop() || "photo.jpg";

      const presign = await ossPresignApi({ filename: fileName });
      await putFileToPresignUrl(presign.url, filePath);

      await photoCreateApi({
        title: "小程序上传",
        description: "",
        srcs: [presign.cdn_url],
        province: province.region_id,
        city: city.region_id,
      });

      Taro.showToast({ title: "上传成功", icon: "success" });
      backToList();
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "上传失败", icon: "none" });
    } finally {
      setUploading(false);
    }
  };

  const provinceName = provinceOptions[provinceIndex]?.region_name || "选择省份";
  const cityName = cityOptions[cityIndex]?.region_name || "选择城市";

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", padding: "24rpx" }}>
      <AtCard title="上传相册">
        <View style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
          <Text style={{ fontSize: "28rpx", color: "#333" }}>打卡城市</Text>
          <AtButton size="small" onClick={autoDetectRegion}>
            自动识别
          </AtButton>
        </View>

        <View style={{ marginTop: "24rpx", display: "flex", gap: "24rpx" }}>
          <Picker
            mode="selector"
            range={provinceOptions}
            rangeKey="region_name"
            onChange={onProvinceChange}
          >
            <View
              style={{
                flex: 1,
                padding: "20rpx",
                borderRadius: "16rpx",
                backgroundColor: "#fff",
                border: "1rpx solid #eee",
              }}
            >
              <Text style={{ fontSize: "26rpx", color: "#333" }}>{provinceName}</Text>
            </View>
          </Picker>

          <Picker
            mode="selector"
            range={cityOptions}
            rangeKey="region_name"
            onChange={onCityChange}
          >
            <View
              style={{
                flex: 1,
                padding: "20rpx",
                borderRadius: "16rpx",
                backgroundColor: "#fff",
                border: "1rpx solid #eee",
              }}
            >
              <Text style={{ fontSize: "26rpx", color: "#333" }}>{cityName}</Text>
            </View>
          </Picker>
        </View>

        <View style={{ marginTop: "24rpx" }}>
          <AtButton type="primary" loading={uploading} disabled={uploading} onClick={upload}>
            选择图片并上传
          </AtButton>
        </View>
      </AtCard>
    </View>
  );
}

import Taro, { useDidShow } from "@tarojs/taro";
import { View, Text } from "@tarojs/components";
import { AtButton, AtCard } from "taro-ui";
import { miniappLoginApi } from "../../service";

const ACCESS_TOKEN_STORAGE_KEY = "access_token";
const MOOD_TAB_URL = "/pages/mood/index";

export default function LoginPage() {
  const redirectToMoodIfLoggedIn = (): boolean => {
    const token = Taro.getStorageSync(ACCESS_TOKEN_STORAGE_KEY);
    if (!token) {
      return false;
    }

    // 本地已存在 access token，说明已登录过；直接进入“心情”页，跳过登录页
    Taro.switchTab({ url: MOOD_TAB_URL });
    return true;
  };

  useDidShow(() => {
    redirectToMoodIfLoggedIn();
  });

  const handleLogin = async () => {
    // 兜底：如果用户从其他入口进入登录页，但本地已有 token，直接跳转即可
    if (redirectToMoodIfLoggedIn()) {
      return;
    }

    // 通过微信登录接口拿到 code，再由后端换取 access token
    const wxLogin = await Taro.login();
    if (!wxLogin.code) {
      Taro.showToast({ title: "获取登录 code 失败", icon: "none" });
      return;
    }

    try {
      const resp = await miniappLoginApi({ code: wxLogin.code });
      Taro.setStorageSync(ACCESS_TOKEN_STORAGE_KEY, resp.access_token);
      Taro.switchTab({ url: MOOD_TAB_URL });
    } catch (e: any) {
      Taro.showToast({ title: e?.message || "登录失败", icon: "none" });
    }
  };

  return (
    <View style={{ minHeight: "100vh", backgroundColor: "#f5f5f5", padding: "48rpx 32rpx" }}>
      <AtCard title="fifsky" note="使用微信一键登录，登录后可发表心情、上传相册、管理提醒">
        <AtButton type="primary" onClick={handleLogin}>
          微信一键登录
        </AtButton>
      </AtCard>
      <View style={{ marginTop: "32rpx" }}>
        <Text style={{ fontSize: "24rpx", color: "#999" }}>
          本地：127.0.0.1:8080 / 线上：api.fifsky.com
        </Text>
      </View>
    </View>
  );
}

export default defineAppConfig({
  pages: ["pages/mood/index", "pages/photo/index", "pages/remind/index", "pages/login/index"],
  tabBar: {
    color: "#64748b",
    selectedColor: "#4f46e5",
    backgroundColor: "#ffffff",
    borderStyle: "black",
    list: [
      {
        pagePath: "pages/mood/index",
        text: "心情",
        iconPath: "images/mood.png",
        selectedIconPath: "images/mood-selected.png",
      },
      {
        pagePath: "pages/photo/index",
        text: "相册",
        iconPath: "images/photo.png",
        selectedIconPath: "images/photo-selected.png",
      },
      {
        pagePath: "pages/remind/index",
        text: "提醒",
        iconPath: "images/remind.png",
        selectedIconPath: "images/remind-selected.png",
      },
    ],
  },
  usingComponents: {},
  window: {
    backgroundTextStyle: "light",
    navigationBarBackgroundColor: "#89d5ef",
    navigationBarTitleText: "無處告別",
    navigationBarTextStyle: "black",
    backgroundColor: "#89d5ef",
    enablePullDownRefresh: true,
  },
  requiredPrivateInfos: ["getLocation"],
  permission: {
    "scope.userLocation": {
      desc: "你的位置信息将用于打卡坐标",
    },
  },
});

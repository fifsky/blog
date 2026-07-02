import SwiftUI

/// 主内容视图
/// 包含四个标签页：博文、心情、提醒、足迹
struct ContentView: View {

    @State private var selectedTab = 0

    var body: some View {
        TabView(selection: $selectedTab) {
            // 博文
            NavigationStack {
                ArticleListView()
            }
            .tabItem {
                Label("博文", systemImage: "doc.text")
            }
            .tag(0)

            // 心情
            NavigationStack {
                MoodListView()
            }
            .tabItem {
                Label("心情", systemImage: "face.smiling")
            }
            .tag(1)

            // 提醒
            NavigationStack {
                RemindListView()
            }
            .tabItem {
                Label("提醒", systemImage: "bell")
            }
            .tag(2)

            // 足迹
            NavigationStack {
                FootprintMapView(
                    onShowListView: {
                        // TODO: 后续接入足迹列表页切换
                    },
                    onAddFootprint: {
                        // 由 FootprintMapView 内部 sheet 处理
                    }
                )
            }
            .tabItem {
                Label("足迹", systemImage: "map")
            }
            .tag(3)
        }
    }
}

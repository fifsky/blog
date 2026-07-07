import SwiftUI
import UIKit

/// 主界面的 UIKit TabBar 外壳
final class AppTabBarController: UITabBarController {

    override func viewDidLoad() {
        super.viewDidLoad()
        // UIKit 外壳直接绑定选中色，避免依赖 SwiftUI AccentColor 的隐式传递。
        tabBar.tintColor = UIColor(named: "AccentColor") ?? UIColor(
            red: 0x00 / 255.0,
            green: 0x7D / 255.0,
            blue: 0xFE / 255.0,
            alpha: 1.0
        )

        viewControllers = [
            AppNavigationController(
                rootView: ArticleListView(),
                title: "博文",
                systemImage: "doc.text.fill"
            ),
            AppNavigationController(
                rootView: MoodListView(),
                title: "心情",
                systemImage: "face.smiling"
            ),
            AppNavigationController(
                rootView: RemindListView(),
                title: "提醒",
                systemImage: "bell.fill"
            ),
            AppNavigationController(
                rootView: FootprintMapView(),
                title: "足迹",
                systemImage: "map.fill",
                rootTitle: "足迹地图"
            )
        ]

        selectedIndex = 0
    }
}

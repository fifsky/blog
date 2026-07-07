import SwiftUI
import UIKit

/// 单个 Tab 内部的 UIKit 导航栈
final class AppNavigationController: UINavigationController {

    let appNavigator: AppNavigator

    /// 初始化一个承载 SwiftUI 根页面的导航栈
    /// - Parameters:
    ///   - rootView: Tab 根页面
    ///   - title: TabBar 标题
    ///   - systemImage: TabBar 图标
    ///   - rootTitle: 根页面导航栏标题
    init<Root: View>(
        rootView: Root,
        title: String,
        systemImage: String,
        rootTitle: String? = nil
    ) {
        let navigator = AppNavigator()
        self.appNavigator = navigator

        let controller = NavigationHostingController(
            rootView: rootView
                .environment(\.appNavigator, navigator)
                .tint(Color.themePrimary)
        )
        controller.title = rootTitle ?? title
        controller.navigationItem.largeTitleDisplayMode = .always
        controller.prefersBottomBarHidden = false
        super.init(rootViewController: controller)

        navigator.navigationController = self
        delegate = self
        navigationBar.prefersLargeTitles = true
        tabBarItem = UITabBarItem(
            title: title,
            image: UIImage(systemName: systemImage),
            selectedImage: UIImage(systemName: systemImage)
        )
    }

    @available(*, unavailable)
    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
}

extension AppNavigationController: UINavigationControllerDelegate {

    func navigationController(
        _ navigationController: UINavigationController,
        willShow viewController: UIViewController,
        animated: Bool
    ) {
        let shouldHideTabBar = shouldHideBottomBar(for: viewController)
        setTabBarHidden(shouldHideTabBar)

        transitionCoordinator?.animate(alongsideTransition: nil) { [weak self] context in
            guard context.isCancelled,
                  let self,
                  let topViewController = self.topViewController else {
                return
            }
            self.setTabBarHidden(self.shouldHideBottomBar(for: topViewController))
        }
    }

    /// 判断指定页面是否需要隐藏底部 TabBar
    private func shouldHideBottomBar(for viewController: UIViewController) -> Bool {
        guard viewController !== viewControllers.first else { return false }
        return (viewController as? BottomBarVisibilityProviding)?.prefersBottomBarHidden ?? true
    }

    /// 直接切换真实 TabBar，避开 `hidesBottomBarWhenPushed` 在 iOS 26 浮动 TabBar 上的空白过渡
    private func setTabBarHidden(_ hidden: Bool) {
        guard let tabBar = tabBarController?.tabBar,
              tabBar.isHidden != hidden else {
            return
        }
        UIView.performWithoutAnimation {
            tabBar.isHidden = hidden
            tabBar.superview?.layoutIfNeeded()
        }
    }
}

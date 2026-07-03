import SwiftUI
import UIKit

extension View {

    /// 点击视图空白区域收起键盘
    ///
    /// 使用 simultaneousGesture 同时识别手势，不影响子视图（Button、Form、列表项等）的原有交互。
    /// 采用 minimumDistance 为 0 的 DragGesture，可同时响应轻点与轻划，对 ScrollView/Form 更稳健。
    func hideKeyboardOnTap() -> some View {
        simultaneousGesture(
            DragGesture(minimumDistance: 0)
                .onEnded { _ in
                    UIApplication.shared.sendAction(
                        #selector(UIResponder.resignFirstResponder),
                        to: nil,
                        from: nil,
                        for: nil
                    )
                }
        )
    }
}

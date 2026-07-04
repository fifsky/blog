import SwiftUI
import UIKit

extension View {

    /// 点击视图空白区域收起键盘
    ///
    /// 使用 simultaneousGesture 同时识别手势，不影响子视图（Button、Form、列表项等）的原有交互。
    /// 拖拽收起键盘由调用方的 scrollDismissesKeyboard 处理，避免 DragGesture 抢占表单控件点击。
    func hideKeyboardOnTap() -> some View {
        simultaneousGesture(
            TapGesture()
                .onEnded {
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

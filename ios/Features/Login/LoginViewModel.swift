import Foundation

/// 登录步骤
enum LoginStep {
    /// 第一步：输入用户名密码
    case credentials
    /// 第二步：输入 TOTP 验证码
    case totp
}

/// 登录视图模型
@Observable
class LoginViewModel {

    // MARK: - 状态

    /// 当前登录步骤
    var step: LoginStep = .credentials

    /// 用户名
    var userName = ""

    /// 密码
    var password = ""

    /// TOTP 验证码（第二步输入）
    var totpCode = ""

    /// 是否正在加载中
    var isLoading = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    /// 登录是否成功
    var isLoginSuccess = false

    // MARK: - 私有属性

    private let authManager = AuthManager.shared

    // MARK: - 计算属性

    /// 第一步登录按钮是否可用
    var isCredentialsButtonDisabled: Bool {
        userName.isEmpty || password.isEmpty
    }

    /// 第二步验证按钮是否可用（6 位验证码）
    var isTotpButtonDisabled: Bool {
        totpCode.count != 6
    }

    // MARK: - 登录方法

    /// 第一步：提交用户名密码
    /// 若账号开启了 TOTP，后端返回 require_totp=true，进入第二步
    func submitCredentials() async {
        isLoading = true
        errorMessage = nil
        showError = false

        do {
            let response = try await authManager.login(
                userName: userName,
                password: password,
                totpCode: nil
            )

            // 需要二次验证，进入第二步
            if response.require_totp == true {
                step = .totp
                totpCode = ""
                isLoading = false
                return
            }

            // 未开启 TOTP，直接登录成功
            isLoginSuccess = true
        } catch {
            if !error.isCancellation {
                errorMessage = error.localizedDescription
                showError = true
            }
        }

        isLoading = false
    }

    /// 第二步：提交 TOTP 验证码完成登录
    func submitTotp() async {
        guard totpCode.count == 6 else { return }

        isLoading = true
        errorMessage = nil
        showError = false

        do {
            let response = try await authManager.login(
                userName: userName,
                password: password,
                totpCode: totpCode
            )

            // 验证码错误时后端返回错误（走 catch），到这里即为登录成功
            if response.require_totp != true {
                isLoginSuccess = true
            }
        } catch {
            if !error.isCancellation {
                errorMessage = error.localizedDescription
                showError = true
            }
        }

        isLoading = false
    }

    /// 返回第一步（修改用户名密码）
    func backToCredentials() {
        step = .credentials
        totpCode = ""
        errorMessage = nil
        showError = false
    }
}

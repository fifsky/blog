import SwiftUI

/// 登录视图
/// 分两步登录：第一步输入用户名密码，第二步（账号开启 TOTP 时）输入 6 位验证码
struct LoginView: View {

    @State private var viewModel = LoginViewModel()

    /// 登录成功回调
    var onLoginSuccess: (LoginViewModel) -> Void

    var body: some View {
        NavigationStack {
            ZStack {
                // 用 switch + transition 实现两步页面切换
                switch viewModel.step {
                case .credentials:
                    credentialsStep
                        .transition(.asymmetric(
                            insertion: .move(edge: .trailing).combined(with: .opacity),
                            removal: .move(edge: .leading).combined(with: .opacity)
                        ))
                case .totp:
                    totpStep
                        .transition(.asymmetric(
                            insertion: .move(edge: .trailing).combined(with: .opacity),
                            removal: .move(edge: .leading).combined(with: .opacity)
                        ))
                }
            }
            .animation(.easeInOut(duration: 0.25), value: viewModel.step)
            .alert("登录失败", isPresented: $viewModel.showError) {
                Button("确定", role: .cancel) {}
            } message: {
                Text(viewModel.errorMessage ?? "未知错误")
            }
            .onChange(of: viewModel.isLoginSuccess) {
                if viewModel.isLoginSuccess {
                    onLoginSuccess(viewModel)
                }
            }
        }
    }

    // MARK: - 第一步：用户名密码

    /// 第一步登录页面
    private var credentialsStep: some View {
        VStack(spacing: 24) {
            Spacer()

            // 标题
            VStack(spacing: 8) {
                Text("博客")
                    .font(.largeTitle)
                    .bold()
                Text("登录以管理你的内容")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }

            // 登录表单
            VStack(spacing: 16) {
                // 用户名
                TextField("用户名", text: $viewModel.userName)
                    .textFieldStyle(.roundedBorder)
                    .textContentType(.username)
                    .textInputAutocapitalization(.never)
                    .autocorrectionDisabled()
                    .submitLabel(.next)
                    .onSubmit { submitCredentials() }

                // 密码
                SecureField("密码", text: $viewModel.password)
                    .textFieldStyle(.roundedBorder)
                    .textContentType(.password)
                    .submitLabel(.go)
                    .onSubmit { submitCredentials() }

                // 登录按钮
                Button {
                    submitCredentials()
                } label: {
                    HStack {
                        if viewModel.isLoading {
                            ProgressView()
                                .tint(.white)
                        }
                        Text("登录")
                    }
                    .frame(maxWidth: .infinity)
                }
                .buttonStyle(.borderedProminent)
                .disabled(viewModel.isCredentialsButtonDisabled || viewModel.isLoading)
            }
            .padding(.horizontal, 24)

            Spacer()
            Spacer()
        }
        .navigationTitle("登录")
    }

    // MARK: - 第二步：TOTP 验证码

    /// 第二步验证码页面
    private var totpStep: some View {
        VStack(spacing: 24) {
            Spacer()

            // 图标 + 标题
            VStack(spacing: 12) {
                Image(systemName: "lock.shield")
                    .font(.system(size: 56))
                    .foregroundStyle(.blue)

                VStack(spacing: 8) {
                    Text("两步验证")
                        .font(.title2)
                        .bold()
                    Text("请输入身份验证器中的 6 位验证码")
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                        .multilineTextAlignment(.center)
                }
            }

            // 6 位验证码输入
            VStack(spacing: 20) {
                // 大号居中的验证码输入框，限制 6 位数字
                TextField("", text: $viewModel.totpCode)
                    .keyboardType(.numberPad)
                    .textContentType(.oneTimeCode)
                    .multilineTextAlignment(.center)
                    .font(.system(size: 34, weight: .semibold, design: .monospaced))
                    .frame(maxWidth: 220)
                    .padding(.vertical, 14)
                    .background(
                        RoundedRectangle(cornerRadius: 12)
                            .strokeBorder(totpFieldColor, lineWidth: 1.5)
                    )
                    .onChange(of: viewModel.totpCode) { _, newValue in
                        // 仅保留数字，最多 6 位
                        let filtered = newValue.filter { $0.isNumber }
                        viewModel.totpCode = String(filtered.prefix(6))
                        // 输入满 6 位自动提交
                        if viewModel.totpCode.count == 6 {
                            submitTotp()
                        }
                    }

                // 提示输入位数
                Text("\(viewModel.totpCode.count) / 6")
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }

            // 验证按钮
            Button {
                submitTotp()
            } label: {
                HStack {
                    if viewModel.isLoading {
                        ProgressView()
                            .tint(.white)
                    }
                    Text("验证")
                }
                .frame(maxWidth: .infinity)
            }
            .buttonStyle(.borderedProminent)
            .disabled(viewModel.isTotpButtonDisabled || viewModel.isLoading)
            .padding(.horizontal, 24)

            // 返回第一步
            Button {
                viewModel.backToCredentials()
            } label: {
                Text("返回重新输入账号")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }

            Spacer()
            Spacer()
        }
        .navigationTitle("两步验证")
    }

    // MARK: - 计算属性

    /// 验证码输入框边框颜色：未满 6 位为次要色，满 6 位为蓝色
    private var totpFieldColor: Color {
        viewModel.totpCode.count == 6 ? .blue : .secondary.opacity(0.4)
    }

    // MARK: - 辅助方法

    /// 提交第一步
    private func submitCredentials() {
        guard !viewModel.isCredentialsButtonDisabled else { return }
        Task {
            await viewModel.submitCredentials()
        }
    }

    /// 提交第二步验证码
    private func submitTotp() {
        guard !viewModel.isTotpButtonDisabled else { return }
        Task {
            await viewModel.submitTotp()
        }
    }
}

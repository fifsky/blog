import SwiftUI

/// 表情选择器视图
/// 以底部弹出的方式展示，按分类展示常用表情
struct EmojiPicker: View {

    /// 选中表情的回调
    var onSelect: (String) -> Void

    /// 搜索关键词
    @State private var searchText = ""

    /// 当前选中的分类
    @State private var selectedCategory: EmojiCategory = .frequent

    var body: some View {
        NavigationStack {
            VStack(spacing: 0) {
                // 分类选择器
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 12) {
                        ForEach(EmojiCategory.allCases) { category in
                            Button {
                                selectedCategory = category
                            } label: {
                                Text(category.label)
                                    .font(.subheadline)
                                    .foregroundStyle(
                                        selectedCategory == category ? Color.accentColor : .secondary
                                    )
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 6)
                                    .background(
                                        selectedCategory == category
                                            ? Color.accentColor.opacity(0.12)
                                            : Color.clear
                                    )
                                    .clipShape(Capsule())
                            }
                        }
                    }
                    .padding(.horizontal, 16)
                    .padding(.vertical, 8)
                }

                Divider()

                // 表情网格
                ScrollView {
                    let emojis = selectedCategory.emojis
                    LazyVGrid(columns: Array(repeating: GridItem(.flexible(), spacing: 8), count: 8), spacing: 8) {
                        ForEach(emojis, id: \.self) { emoji in
                            Button {
                                onSelect(emoji)
                            } label: {
                                Text(emoji)
                                    .font(.title2)
                                    .frame(width: 44, height: 44)
                                    .background(Color(.systemGray6))
                                    .clipShape(RoundedRectangle(cornerRadius: 8))
                            }
                        }
                    }
                    .padding(12)
                }
            }
            .navigationTitle("选择表情")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .topBarLeading) {
                    Button("完成") {
                        // 由父视图通过 isPresented 控制
                    }
                }
            }
        }
        .presentationDetents([.medium, .large])
        .presentationDragIndicator(.visible)
    }
}

// MARK: - 表情分类

/// 表情分类枚举
enum EmojiCategory: String, CaseIterable, Identifiable {
    case frequent
    case expressions
    case gestures
    case animals
    case food
    case nature
    case objects

    var id: String { rawValue }

    /// 分类显示名称
    var label: String {
        switch self {
        case .frequent: return "常用"
        case .expressions: return "表情"
        case .gestures: return "手势"
        case .animals: return "动物"
        case .food: return "食物"
        case .nature: return "自然"
        case .objects: return "物品"
        }
    }

    /// 分类下的表情列表
    var emojis: [String] {
        switch self {
        case .frequent:
            return ["😊", "😂", "🥰", "😍", "😢", "😭", "😤", "🥺",
                    "👍", "👎", "❤️", "🔥", "✨", "🎉", "💪", "🙏",
                    "😅", "😎", "🤔", "😴", "🤗", "😢", "😈", "🤡"]
        case .expressions:
            return ["😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂",
                    "🙂", "🙃", "😉", "😊", "😇", "🥰", "😍", "🤩",
                    "😘", "😗", "😚", "😙", "🥲", "😋", "😛", "😜",
                    "🤪", "😝", "🤑", "🤗", "🤭", "🤫", "🤔", "🫡",
                    "😐", "😑", "😶", "🫥", "😏", "😒", "🙄", "😬",
                    "🤥", "😌", "😔", "😪", "🤤", "😴", "😷", "🤒"]
        case .gestures:
            return ["👋", "🤚", "🖐️", "✋", "🖖", "👌", "🤌", "🤏",
                    "✌️", "🤞", "🤟", "🤘", "🤙", "👈", "👉", "👆",
                    "👇", "☝️", "👍", "👎", "✊", "👊", "🤛", "🤜",
                    "👏", "🙌", "🫶", "👐", "🤲", "🤝", "🙏", "💪",
                    "🦾", "🦿", "🦵", "🦶", "👂", "🦻", "👃", "🧠"]
        case .animals:
            return ["🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼",
                    "🐨", "🐯", "🦁", "🐮", "🐷", "🐸", "🐵", "🐔",
                    "🐧", "🐦", "🐤", "🦆", "🦅", "🦉", "🦇", "🐺",
                    "🐗", "🐴", "🦄", "🐝", "🐛", "🦋", "🐌", "🐞",
                    "🐜", "🦗", "🕷️", "🦂", "🐢", "🐍", "🦎", "🦖"]
        case .food:
            return ["🍎", "🍐", "🍊", "🍋", "🍌", "🍉", "🍇", "🍓",
                    "🍈", "🍒", "🍑", "🥭", "🍍", "🥥", "🥝", "🍅",
                    "🍔", "🍕", "🌭", "🥪", "🌮", "🌯", "🥙", "🧆",
                    "🥚", "🍳", "🥘", "🍲", "🥣", "🥗", "🍿", "🧈",
                    "🍱", "🍙", "🍚", "🍛", "🍜", "🍝", "🍠", "🍢",
                    "🍣", "🍤", "🍡", "🥟", "🥠", "🥡", "🍦", "🍾"]
        case .nature:
            return ["🌸", "💮", "🏵️", "🌹", "🥀", "🌺", "🌻", "🌼",
                    "🌷", "🌱", "🪴", "🌲", "🌳", "🌴", "🌵", "🌾",
                    "🌿", "☘️", "🍀", "🍁", "🍂", "🍃", "🍄", "🌰",
                    "🌍", "🌏", "🌐", "🌑", "🌒", "🌓", "🌔", "🌕",
                    "🌖", "🌗", "🌘", "🌙", "⭐", "🌟", "✨", "💫",
                    "🔥", "🌈", "☀️", "🌤️", "⛅", "🌥️", "☁️", "🌧️"]
        case .objects:
            return ["💡", "🔦", "📱", "💻", "⌨️", "🖥️", "🖨️", "🖱️",
                    "💾", "📷", "📹", "🎥", "📺", "📻", "📚", "📖",
                    "📝", "✏️", "🖋️", "🖊️", "🖌️", "💬", "💭", "🗑️",
                    "🔒", "🔑", "🔔", "🎵", "🎶", "🎸", "🎹", "🥁",
                    "🎺", "🎻", "🏆", "🥇", "🥈", "🥉", "⚽", "🏀",
                    "🏈", "⚾", "🎾", "🏐", "🎱", "🎲", "🎯", "🎮"]
        }
    }
}

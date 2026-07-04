import SwiftUI

/// 足迹列表视图
/// 以卡片形式展示所有足迹，支持下拉刷新和切换地图视图
struct FootprintListView: View {

    @State private var viewModel = FootprintListViewModel()

    /// 切换到地图视图的回调
    var onShowMapView: () -> Void

    /// 新建足迹的回调
    var onAddFootprint: () -> Void

    /// 点击足迹卡片的回调
    var onSelectFootprint: (Footprint) -> Void

    var body: some View {
        Group {
            if viewModel.isLoading && viewModel.footprints.isEmpty {
                // 首次加载中
                loadingView
            } else if viewModel.footprints.isEmpty {
                // 空状态
                emptyView
            } else {
                // 足迹列表
                footprintGrid
            }
        }
        .navigationTitle("足迹")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .topBarTrailing) {
                HStack(spacing: 16) {
                    // 切换到地图视图
                    Button {
                        onShowMapView()
                    } label: {
                        Image(systemName: "map")
                    }

                    // 新建足迹
                    Button {
                        onAddFootprint()
                    } label: {
                        Image(systemName: "plus")
                    }
                }
            }
        }
        .task {
            await viewModel.loadFootprints()
        }
        .alert("错误", isPresented: $viewModel.showError) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(viewModel.errorMessage ?? "未知错误")
        }
    }

    // MARK: - 加载视图

    /// 加载中占位视图
    private var loadingView: some View {
        VStack(spacing: 12) {
            ProgressView()
            Text("加载中...")
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }

    // MARK: - 空状态

    /// 空状态视图
    private var emptyView: some View {
        VStack(spacing: 16) {
            Image(systemName: "map")
                .font(.system(size: 48))
                .foregroundStyle(.secondary)

            Text("还没有足迹")
                .font(.title3)
                .foregroundStyle(.secondary)

            Text("点击右上角 + 开始记录你的足迹吧")
                .font(.subheadline)
                .foregroundStyle(.tertiary)

            Button {
                onAddFootprint()
            } label: {
                Text("添加足迹")
                    .fontWeight(.medium)
            }
            .buttonStyle(.borderedProminent)
            .padding(.top, 8)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }

    // MARK: - 足迹网格

    /// 足迹卡片网格
    private var footprintGrid: some View {
        ScrollView {
            LazyVGrid(
                columns: [
                    GridItem(.flexible(), spacing: 12),
                    GridItem(.flexible(), spacing: 12)
                ],
                spacing: 12
            ) {
                ForEach(viewModel.footprints) { footprint in
                    FootprintCardView(footprint: footprint)
                        .onTapGesture {
                            onSelectFootprint(footprint)
                        }
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 8)

            if viewModel.isLoadingMore {
                HStack {
                    Spacer()
                    ProgressView()
                        .padding()
                    Spacer()
                }
            } else if viewModel.hasMore {
                Color.clear
                    .frame(height: 1)
                    .onAppear {
                        Task { await viewModel.loadMore() }
                    }
            }
        }
        .refreshable {
            await viewModel.refresh()
        }
    }
}

// MARK: - 足迹卡片视图

/// 单个足迹的卡片视图
struct FootprintCardView: View {

    /// 足迹数据
    let footprint: Footprint

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // 缩略图
            thumbnailSection

            // 信息区域
            VStack(alignment: .leading, spacing: 6) {
                // 名称
                Text(footprint.name ?? "未命名")
                    .font(.headline)
                    .lineLimit(1)

                // 日期
                if let date = footprint.date, !date.isEmpty {
                    Text(date)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }

                // 描述
                if let desc = footprint.description, !desc.isEmpty {
                    Text(desc)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                        .lineLimit(2)
                }

                // 照片数量标记
                if let photos = footprint.photos, !photos.isEmpty {
                    HStack(spacing: 4) {
                        Image(systemName: "photo.on.rectangle.angled")
                            .font(.caption2)
                        Text("\(photos.count)张照片")
                            .font(.caption2)
                    }
                    .foregroundStyle(.tertiary)
                }
            }
            .padding(10)
        }
        .background(Color(.systemBackground))
        .clipShape(RoundedRectangle(cornerRadius: 12))
        .shadow(color: .black.opacity(0.08), radius: 4, x: 0, y: 2)
    }

    /// 缩略图区域
    private var thumbnailSection: some View {
        Group {
            if let photos = footprint.photos, let firstPhoto = photos.first {
                let thumb = firstPhoto.thumbnail ?? ""
                let src = firstPhoto.src ?? ""
                AsyncImage(url: URL(string: thumb.isEmpty ? src : thumb)) { phase in
                    switch phase {
                    case .success(let image):
                        image
                            .resizable()
                            .aspectRatio(contentMode: .fill)
                    case .failure:
                        placeholderImage
                    default:
                        ProgressView()
                            .frame(height: 140)
                    }
                }
                .frame(height: 140)
                .clipShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
            } else {
                placeholderImage
            }
        }
        .frame(height: 140)
    }

    /// 占位图片（无照片时显示）
    private var placeholderImage: some View {
        let markerColor = footprint.marker_color ?? "#FF3B30"
        return ZStack {
            Rectangle()
                .fill(Color(.secondarySystemBackground))

            VStack(spacing: 6) {
                Image(systemName: markerColor.isEmpty ? "mappin" : "mappin.and.square")
                    .font(.title2)
                    .foregroundStyle(Color(hex: markerColor) ?? .blue)

                Text(footprint.name ?? "未命名")
                    .font(.caption)
                    .foregroundStyle(.secondary)
                    .lineLimit(1)
                    .padding(.horizontal, 8)
            }
        }
        .clipShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
    }
}

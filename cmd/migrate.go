package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"app/config"

	"github.com/urfave/cli/v3"
)

func NewMigrate(db *sql.DB, conf *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "数据迁移工具",
		Commands: []*cli.Command{
			{
				Name:  "photo-to-footprint",
				Usage: "将 photos 表数据按城市分组迁移到 footprints 表",
				Action: func(ctx context.Context, cli *cli.Command) error {
					return migratePhotoToFootprint(ctx, db)
				},
			},
			{
				Name:  "drop-photos",
				Usage: "删除 photos 表（迁移完成后执行）",
				Action: func(ctx context.Context, cli *cli.Command) error {
					return dropPhotosTable(ctx, db)
				},
			},
		},
	}
}

type photoRow struct {
	id          int
	title       string
	description string
	src         string
	thumbnail   string
	province    int
	city        int
	createdAt   string
}

type regionRow struct {
	regionID   int
	regionName string
	longitude  string
	latitude   string
}

func migratePhotoToFootprint(ctx context.Context, db *sql.DB) error {
	// 检查 footprints 表是否存在
	var tableExists bool
	err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'footprints'",
	).Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("检查 footprints 表失败: %w", err)
	}
	if !tableExists {
		log.Println("footprints 表不存在，开始创建...")
		_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS footprints (
			id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
			name varchar(256) NOT NULL COMMENT '地点名称',
			description varchar(1024) NOT NULL DEFAULT '' COMMENT '描述',
			longitude varchar(20) NOT NULL COMMENT '经度',
			latitude varchar(20) NOT NULL COMMENT '纬度',
			date varchar(128) NOT NULL DEFAULT '' COMMENT '到访日期',
			marker_color varchar(32) NOT NULL DEFAULT '' COMMENT '标记颜色',
			categories json DEFAULT NULL COMMENT '分类标签数组',
			url varchar(1024) NOT NULL DEFAULT '' COMMENT '关联链接',
			url_label varchar(128) NOT NULL DEFAULT '' COMMENT '链接按钮文案',
			photos json DEFAULT NULL COMMENT '照片URL数组',
			created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='旅行足迹'`)
		if err != nil {
			return fmt.Errorf("创建 footprints 表失败: %w", err)
		}
		log.Println("footprints 表创建成功")
	}

	// 检查 photos 表是否存在
	var photosExists bool
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'photos'",
	).Scan(&photosExists)
	if err != nil {
		return fmt.Errorf("检查 photos 表失败: %w", err)
	}
	if !photosExists {
		log.Println("photos 表不存在，跳过迁移")
		return nil
	}

	// 查询所有照片
	rows, err := db.QueryContext(ctx,
		"SELECT id, title, description, src, thumbnail, province, city, DATE_FORMAT(created_at, '%Y-%m-%d') FROM photos ORDER BY city, id",
	)
	if err != nil {
		return fmt.Errorf("查询 photos 失败: %w", err)
	}
	defer rows.Close()

	var photos []photoRow
	for rows.Next() {
		var p photoRow
		if err := rows.Scan(&p.id, &p.title, &p.description, &p.src, &p.thumbnail, &p.province, &p.city, &p.createdAt); err != nil {
			return fmt.Errorf("扫描照片数据失败: %w", err)
		}
		photos = append(photos, p)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(photos) == 0 {
		log.Println("photos 表为空，跳过迁移")
		return nil
	}
	log.Printf("发现 %d 条照片记录\n", len(photos))

	// 查询城市区域信息
	regionMap := make(map[int]regionRow)
	regionRows, err := db.QueryContext(ctx, "SELECT region_id, region_name, longitude, latitude FROM regions WHERE level IN (1, 2)")
	if err != nil {
		return fmt.Errorf("查询 regions 失败: %w", err)
	}
	defer regionRows.Close()
	for regionRows.Next() {
		var r regionRow
		if err := regionRows.Scan(&r.regionID, &r.regionName, &r.longitude, &r.latitude); err != nil {
			return fmt.Errorf("扫描区域数据失败: %w", err)
		}
		regionMap[r.regionID] = r
	}
	if err := regionRows.Err(); err != nil {
		return err
	}

	// 按城市分组
	type cityGroup struct {
		cityID int
		photos []photoRow
	}
	groupMap := make(map[int]*cityGroup)
	var groupOrder []int

	for _, p := range photos {
		g, ok := groupMap[p.city]
		if !ok {
			g = &cityGroup{cityID: p.city}
			groupMap[p.city] = g
			groupOrder = append(groupOrder, p.city)
		}
		g.photos = append(g.photos, p)
	}

	// 逐组插入 footprints
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	insertStmt, err := tx.PrepareContext(ctx,
		"INSERT INTO footprints (name, description, longitude, latitude, date, marker_color, categories, photos) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return fmt.Errorf("准备插入语句失败: %w", err)
	}
	defer insertStmt.Close()

	for _, cityID := range groupOrder {
		g := groupMap[cityID]
		region := regionMap[cityID]

		name := region.regionName
		if name == "" {
			name = fmt.Sprintf("未知地点-%d", cityID)
		}

		lng := region.longitude
		lat := region.latitude
		if lng == "" || lat == "" {
			// 尝试用省份坐标
			if len(g.photos) > 0 {
				provRegion := regionMap[g.photos[0].province]
				if provRegion.longitude != "" {
					lng = provRegion.longitude
					lat = provRegion.latitude
				}
			}
		}
		if lng == "" || lat == "" {
			lng = "116.4074"
			lat = "39.9042"
		}

		desc := g.photos[0].description
		if desc == "" {
			desc = g.photos[0].title
		}

		date := g.photos[0].createdAt[:7] // YYYY-MM

		photoItems := make([]map[string]string, 0, len(g.photos))
		for _, p := range g.photos {
			photoItems = append(photoItems, map[string]string{
				"src":       p.src,
				"thumbnail": p.thumbnail,
			})
		}
		photosJSON, _ := json.Marshal(photoItems)

		categories, _ := json.Marshal([]string{"旅行"})

		result, err := insertStmt.ExecContext(ctx, name, desc, lng, lat, date, "", categories, photosJSON)
		if err != nil {
			return fmt.Errorf("插入足迹失败 (city=%d): %w", cityID, err)
		}
		id, _ := result.LastInsertId()
		log.Printf("  插入足迹: id=%d, name=%s, lng=%s, lat=%s, photos=%d", id, name, lng, lat, len(g.photos))
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	log.Printf("\n迁移完成! 共插入 %d 条足迹记录\n", len(groupOrder))
	log.Println("确认无误后可执行: ./app migrate drop-photos")
	return nil
}

func dropPhotosTable(ctx context.Context, db *sql.DB) error {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM footprints").Scan(&count)
	if err != nil {
		return fmt.Errorf("请先完成迁移: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("footprints 表为空，请先执行 photo-to-footprint 迁移")
	}

	log.Printf("footprints 表有 %d 条记录，确认删除 photos 表...\n", count)
	_, err = db.ExecContext(ctx, "DROP TABLE IF EXISTS photos")
	if err != nil {
		return fmt.Errorf("删除 photos 表失败: %w", err)
	}
	log.Println("photos 表已删除")
	return nil
}

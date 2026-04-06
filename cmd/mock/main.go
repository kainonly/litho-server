package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"server/bootstrap"
	"server/common"
	"server/model"

	"github.com/kainonly/go/help"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

const (
	productCount = 200
	orderCount   = 500000
	batchSize    = 1000
)

// ---- 虚构活动名素材（可组合出数千种名称）----

var eventPrefixes = []string{
	"星辰", "幻境", "霓虹", "苍穹", "晨曦", "极光", "银河", "彩云",
	"碧海", "云端", "山野", "浮光", "暗涌", "晚风", "流萤", "烟火",
	"深蓝", "白夜", "赤道", "雨林", "冰川", "沙漠", "熔岩", "峡谷",
	"暮色", "黎明", "正午", "子夜", "黄昏", "朝霞", "紫陌", "翠微",
}

var eventMids = []string{
	"音乐节", "演出季", "狂欢夜", "盛典", "艺术展", "嘉年华",
	"巡演", "演唱会", "沉浸展", "文化节", "创意市集", "电子音乐节",
	"戏剧节", "舞蹈节", "影像节", "光影展", "装置艺术展", "潮流展",
	"脱口秀节", "喜剧节", "爵士节", "古典音乐会", "民谣节", "说唱节",
	"街舞大赛", "竞技赛", "电竞赛", "马拉松", "嘉年华巡游", "烟火晚会",
}

var eventSuffixes = []string{
	"第%d届", "Vol.%d", "%d周年纪念场", "限定场", "特别场",
	"首演场", "收官场", "夏季场", "冬季场", "跨年场",
	"城市限定", "户外版", "线上直播场", "%d城联动", "回归场",
}

var ticketTypes = []struct {
	name  string
	ratio float64
}{
	{"普通票", 1.0},
	{"看台票", 1.2},
	{"内场票", 1.8},
	{"VIP票", 3.0},
	{"SVIP票", 5.0},
	{"早鸟票", 0.8},
	{"学生票", 0.6},
	{"家庭套票", 2.2},
	{"双人票", 1.9},
	{"亲子票", 1.5},
	{"尊享票", 4.0},
	{"联票（两场）", 1.7},
	{"全程通票", 2.5},
	{"单日票", 1.0},
	{"三日通票", 2.8},
}

var baseDescriptions = []string{
	"虚构活动，限量发售，售完即止。",
	"模拟数据，含电子凭证，入场核验。",
	"演示专用，支持退改，详见须知。",
	"示例票券，实名绑定，一人一票。",
	"测试数据，凭证二维码当日有效。",
	"虚构场次，含导览手册，现场领取。",
	"模拟票务，禁止转让，遗失不补。",
	"演示票券，含餐饮折扣券一张。",
}

var remarks = []string{
	"", "", "", "", "", "", "", // 多数无备注
	"需要无障碍通道", "团体购票", "生日场次", "含周边礼包", "指定座区",
	"儿童需持证", "企业团建", "学生证核验", "早入场通道", "残障人士陪护",
	"摄影证申请", "媒体证核验", "嘉宾邀请码", "会员专属场", "城市卡优惠",
}

func main() {
	configPath := flag.String("config", "config/values.yml", "配置文件路径")
	flag.Parse()

	values, err := loadValues(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	db, err := bootstrap.UseGorm(values)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接数据库失败: %v\n", err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var org model.Org
	if err := db.First(&org).Error; err != nil {
		fmt.Fprintf(os.Stderr, "获取组织失败: %v\n", err)
		os.Exit(1)
	}

	var userIDs []string
	if err := db.Model(&model.User{}).Pluck("id", &userIDs).Error; err != nil {
		fmt.Fprintf(os.Stderr, "获取用户失败: %v\n", err)
		os.Exit(1)
	}
	if len(userIDs) == 0 {
		fmt.Fprintf(os.Stderr, "没有用户数据，请先运行 seeders\n")
		os.Exit(1)
	}

	fmt.Printf("使用组织: %s (%s)\n", org.Name, org.ID)
	fmt.Printf("可用用户: %d 个\n", len(userIDs))

	fmt.Println("生成票券产品...")
	products, err := generateProducts(db, org.ID, rng)
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成票券产品失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("生成票券产品: %d 个\n", len(products))

	fmt.Printf("生成订单数据（共 %d 条）...\n", orderCount)
	if err := generateOrders(db, org.ID, userIDs, products, rng); err != nil {
		fmt.Fprintf(os.Stderr, "生成订单失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Mock 数据生成完成！")
}

func randomEventName(rng *rand.Rand) string {
	prefix := eventPrefixes[rng.Intn(len(eventPrefixes))]
	mid := eventMids[rng.Intn(len(eventMids))]
	suffix := fmt.Sprintf(eventSuffixes[rng.Intn(len(eventSuffixes))], rng.Intn(10)+1)
	return prefix + mid + "·" + suffix
}

func generateProducts(db *gorm.DB, orgID string, rng *rand.Rand) ([]model.Product, error) {
	if err := db.Where("org_id = ?", orgID).Delete(&model.Product{}).Error; err != nil {
		return nil, err
	}

	active := true
	inactive := false
	products := make([]model.Product, 0, productCount)
	usedNames := make(map[string]bool)

	for len(products) < productCount {
		eventName := randomEventName(rng)
		tt := ticketTypes[rng.Intn(len(ticketTypes))]
		fullName := eventName + "·" + tt.name
		if usedNames[fullName] {
			continue
		}
		usedNames[fullName] = true

		basePrice := float64(rng.Intn(300) + 100) // 100~400 基础价
		price := math.Round(basePrice*tt.ratio*100) / 100
		stock := rng.Intn(2000) + 100

		// 约 15% 产品下架，模拟可随时调控
		isActive := &active
		if rng.Intn(100) < 15 {
			isActive = &inactive
		}

		products = append(products, model.Product{
			ID:          help.SID(),
			OrgID:       orgID,
			Name:        fullName,
			Description: baseDescriptions[rng.Intn(len(baseDescriptions))],
			Price:       price,
			Stock:       int32(stock),
			Active:      isActive,
		})
	}

	if err := db.CreateInBatches(products, 100).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func generateOrders(db *gorm.DB, orgID string, userIDs []string, products []model.Product, rng *rand.Rand) error {
	if err := db.Exec(`TRUNCATE TABLE order_item`).Error; err != nil {
		return err
	}
	if err := db.Exec(`TRUNCATE TABLE "order"`).Error; err != nil {
		return err
	}

	now := time.Now()
	statuses := []int16{0, 1, 2, 3}
	// 待付款5%，已付款10%，已完成75%，已取消10%
	statusWeights := []int{5, 10, 75, 10}

	orders := make([]model.Order, 0, batchSize)
	items := make([]model.OrderItem, 0, batchSize*2)
	total := 0

	for i := 0; i < orderCount; i++ {
		orderID := help.SID()
		userID := userIDs[rng.Intn(len(userIDs))]

		// 时间分布：近30天占40%（实时演示用），其余分布在过去两年
		var createdAt time.Time
		if rng.Intn(10) < 4 {
			createdAt = now.Add(-time.Duration(rng.Intn(30*24)) * time.Hour)
		} else {
			daysAgo := rng.Intn(700) + 30
			createdAt = now.AddDate(0, 0, -daysAgo)
		}
		createdAt = createdAt.Add(time.Duration(rng.Intn(86400)) * time.Second)

		// 活动场次：下单后 1~90 天，固定在 14/19/20/21 点整
		scheduledAt := createdAt.AddDate(0, 0, rng.Intn(90)+1)
		hour := []int{14, 19, 20, 21}[rng.Intn(4)]
		scheduledAt = scheduledAt.Truncate(24*time.Hour).Add(time.Duration(hour) * time.Hour)

		status := weightedRandom(rng, statuses, statusWeights)

		p := products[rng.Intn(len(products))]
		qty := int32(rng.Intn(4) + 1)
		subtotal := p.Price * float64(qty)
		amount := subtotal

		orderItems := []model.OrderItem{
			{
				ID:          help.SID(),
				OrderID:     parseID(orderID),
				ProductID:   parseID(p.ID),
				ProductName: p.Name,
				Price:       p.Price,
				Quantity:    &qty,
				Subtotal:    subtotal,
			},
		}

		var paidAt, closedAt *time.Time
		if status >= 1 {
			t := createdAt.Add(time.Duration(rng.Intn(1800)+30) * time.Second)
			paidAt = &t
		}
		if status == 2 {
			t := scheduledAt.Add(time.Duration(rng.Intn(3600)) * time.Second)
			closedAt = &t
		}
		if status == 3 {
			t := createdAt.Add(time.Duration(rng.Intn(86400)) * time.Second)
			closedAt = &t
		}

		no := fmt.Sprintf("TK%s%08d", createdAt.Format("20060102"), i%100000000)

		order := model.Order{
			ID:          orderID,
			CreatedAt:   &createdAt,
			UpdatedAt:   &createdAt,
			OrgID:       orgID,
			UserID:      userID,
			No:          no,
			Amount:      amount,
			Status:      status,
			ScheduledAt: scheduledAt,
			Remark:      remarks[rng.Intn(len(remarks))],
			PaidAt:      paidAt,
			ClosedAt:    closedAt,
		}
		orders = append(orders, order)
		items = append(items, orderItems...)

		if len(orders) >= batchSize {
			if err := db.CreateInBatches(&orders, batchSize).Error; err != nil {
				return fmt.Errorf("插入 order 失败: %w", err)
			}
			if err := db.CreateInBatches(&items, batchSize).Error; err != nil {
				return fmt.Errorf("插入 order_item 失败: %w", err)
			}
			total += len(orders)
			fmt.Printf("\r已生成: %d / %d", total, orderCount)
			orders = orders[:0]
			items = items[:0]
		}
	}

	if len(orders) > 0 {
		if err := db.CreateInBatches(&orders, batchSize).Error; err != nil {
			return fmt.Errorf("插入 order 失败: %w", err)
		}
		if err := db.CreateInBatches(&items, batchSize).Error; err != nil {
			return fmt.Errorf("插入 order_item 失败: %w", err)
		}
		total += len(orders)
	}

	fmt.Printf("\r已生成: %d / %d\n", total, orderCount)
	return nil
}

func parseID(id string) int64 {
	var v int64
	fmt.Sscanf(id, "%d", &v)
	return v
}

func weightedRandom(rng *rand.Rand, statuses []int16, weights []int) int16 {
	total := 0
	for _, w := range weights {
		total += w
	}
	r := rng.Intn(total)
	for i, w := range weights {
		r -= w
		if r < 0 {
			return statuses[i]
		}
	}
	return statuses[len(statuses)-1]
}

func loadValues(path string) (*common.Values, error) {
	absPath, err := resolvePath(path)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	v := new(common.Values)
	if err := yaml.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}

func resolvePath(path string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		execPath, _ = os.Getwd()
	}
	candidates := []string{
		filepath.Join(filepath.Dir(execPath), "..", "..", path),
		filepath.Join(filepath.Dir(execPath), path),
		path,
	}
	for _, c := range candidates {
		abs, _ := filepath.Abs(c)
		if _, err := os.Stat(abs); err == nil {
			return abs, nil
		}
	}
	return "", fmt.Errorf("路径不存在: %s", path)
}

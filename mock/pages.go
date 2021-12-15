package main

import (
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productsList = [][]interface{}{
	{"name", model.
		NewField("商品名称", "string").
		SetRequired()},
	{"description", model.
		NewField("商品描述", "text")},
	{"pictures", model.
		NewField("商品图", "picture").
		SetDescription("建议尺寸：800*800像素")},
	{"videos", model.
		NewField("主图视频", "video").
		SetDescription("添加主图视频可提升商品的成交转化，有利于获取更多流量，建议时长 9-30 秒、视频宽高和商品图一致。")},
	{"group", model.
		NewField("商品分组", "select").
		SetSpec(&model.FieldSpec{Reference: "", Target: "", Multiple: model.Bool(true)})},
	{"spec", model.
		NewField("商品规格", "json")},
	{"price", model.
		NewField("价格", "number").
		SetRequired().
		SetSpec(&model.FieldSpec{Min: 0, Decimal: 2})},
	{"underlined_price", model.
		NewField("划线价", "number").
		SetSpec(&model.FieldSpec{Min: 0, Decimal: 2})},
	{"stock_mode", model.
		NewField("库存扣减方式", "radio").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "拍下减库存", Value: "place_an_order"},
				{Label: "付款减库存", Value: "payed"},
			},
		})},
	{"stock", model.
		NewField("库存", "number").
		SetRequired().
		SetSpec(&model.FieldSpec{Min: 0})},
	{"weight", model.
		NewField("重量", "number").
		SetSpec(&model.FieldSpec{Min: 0})},
	{"member_discount", model.
		NewField("会员折扣", "bool").
		SetDescription("是否勾选不影响自定义会员价生效。")},
	{"code", model.
		NewField("商品编码", "string")},
	{"barcode", model.
		NewField("商品条码", "string")},
	{"cost_price", model.
		NewField("成本价", "number").
		SetSpec(&model.FieldSpec{Min: 0, Decimal: 2})},
	{"delivery", model.
		NewField("配送方式", "checkbox").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "快速发货", Value: "send"},
				{Label: "到店自提", Value: "in_store"},
			},
		})},
	{"express_shipping", model.
		NewField("快递运费", "radio").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "统一邮费", Value: "unified"},
				{Label: "运费模板", Value: "template"},
			},
		})},
	{"express_fee", model.
		NewField("统一邮费", "number").
		SetSpec(&model.FieldSpec{Min: 0, Decimal: 2})},
	{"express_ref", model.
		NewField("运费模板", "select").
		SetSpec(&model.FieldSpec{Reference: "", Target: ""})},
	{"sale_mode", model.
		NewField("开售方式", "radio").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "立即开售", Value: "now"},
				{Label: "定时开售", Value: "timing"},
				{Label: "放入仓库", Value: "warehouse"},
			},
		})},
	{"sale_time", model.
		NewField("开售时间", "datetime")},
	{"timed_off", model.
		NewField("定时下架", "bool")},
	{"timed_off_time", model.
		NewField("定时下架时间", "datetime")},
	{"after_sale", model.
		NewField("售后服务", "checkbox").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "支持买家申请换货", Value: "support_exchange"},
				{Label: "7天无理由退货", Value: "7_days_no_reason"},
				{Label: "自动退款", Value: "automatic_refund"},
			},
		})},
	{"pre_sale", model.
		NewField("预售方式", "radio").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "全款预售", Value: "support_exchange"},
				{Label: "定金预售", Value: "deposit"},
				{Label: "关闭", Value: "off"},
			},
		})},
	{"delivery_time_mode", model.
		NewField("发货时间", "radio").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "指定时间", Value: "specify"},
				{Label: "预估时间", Value: "estimate"},
			},
		})},
	{"delivery_time", model.
		NewField("开始发货", "datetime")},
	{"delivery_day", model.
		NewField("天", "number").
		SetSpec(&model.FieldSpec{Min: 0})},
	{"purchase_limit", model.
		NewField("限购", "checkbox").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "限制每人可购买数量", Value: "limit_the_number"},
				{Label: "只允许特定用户购买", Value: "limit_the_group"},
			},
		})},
	{"purchase_limit_strategy", model.
		NewField("限购条件", "radio").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "终身限购", Value: "lifelong"},
				{Label: "按周期限购", Value: "period"},
			},
		})},
	{"purchase_limit_period", model.
		NewField("限购周期", "select").
		SetSpec(&model.FieldSpec{
			Values: []model.Enum{
				{Label: "每天", Value: "day"},
				{Label: "每周", Value: "week"},
				{Label: "每月", Value: "month"},
			},
		})},
	{"purchase_limit_card", model.
		NewField("权益卡用户限定", "select").
		SetSpec(&model.FieldSpec{Reference: "", Target: "", Multiple: model.Bool(true)})},
	{"purchase_limit_member", model.
		NewField("会员用户限定", "select").
		SetSpec(&model.FieldSpec{Reference: "", Target: "", Multiple: model.Bool(true)})},
	{"purchase_limit_label", model.
		NewField("用户标签限定", "select").
		SetSpec(&model.FieldSpec{Reference: "", Target: "", Multiple: model.Bool(true)})},
	{"min_sale_quantity", model.
		NewField("起售数量", "number").
		SetSpec(&model.FieldSpec{Min: 0})},
}

func MockPages(db *mongo.Database) (result *mongo.InsertManyResult, err error) {
	ctx := context.Background()
	if err = db.Collection("pages").Drop(ctx); err != nil {
		return
	}
	productId := primitive.NewObjectID()
	orderId := primitive.NewObjectID()
	productsFields := make(model.SchemaFields, len(productsList))
	for i, v := range productsList {
		productsFields[v[0].(string)] = v[1].(*model.Field).SetSort(int64(i))
	}
	data := []interface{}{
		model.NewPage("商品管理", "group").
			SetID(productId).
			SetIcon("shopping"),
		model.NewPage("商品清单", "default").
			SetParent(&productId).
			SetSchema(model.NewSchema("products", productsFields)),
		model.NewPage("商品分组", "default").
			SetParent(&productId).
			SetSchema(model.NewSchema("product_group", model.SchemaFields{})),
		model.NewPage("商品设置", "form").
			SetParent(&productId).
			SetSchema(model.NewSchema("product_values", model.SchemaFields{})),
		model.NewPage("订单管理", "group").
			SetID(orderId).
			SetIcon("profile"),
		model.NewPage("订单列表", "default").
			SetParent(&orderId).
			SetSchema(model.NewSchema("orders", model.SchemaFields{})),
		model.NewPage("售后维权", "default").
			SetParent(&orderId).
			SetSchema(model.NewSchema("after_sale_orders", model.SchemaFields{})),
	}
	if result, err = db.
		Collection("pages").
		InsertMany(ctx, data); err != nil {
		return
	}
	return
}

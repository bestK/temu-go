package entity

import "gopkg.in/guregu/null.v4"

// ParentOrder PO 单
type ParentOrder struct {
	ParentOrderMap ParentOrderMap `json:"parentOrderMap"` // 父单信息
	OrderList      []Order        `json:"orderList"`      // 商品信息（子单列表）
}

type ParentOrderMap struct {
	ParentOrderSn string `json:"parentOrderSn"` // 父订单号
	// name: 具体枚举如下
	// soon_to_be_overdue-即将逾期
	// past_due-已逾期
	// pending_buyer_cancellation-买家取消待确认订单
	// pending_buyer_address_change-买家改地址待确认订单
	// value:
	// 是否有标签：0=无标签，1=有标签
	ParentOrderLabel             []Label   `json:"parentOrderLabel"`             // 标签名称
	ParentOrderStatus            int       `json:"parentOrderStatus"`            // 订单状态
	ParentOrderTime              int64     `json:"parentOrderTime"`              // 订单创建时间
	ParentOrderPendingFinishTime int64     `json:"parentOrderPendingFinishTime"` // 订单结束pending转为自发货时间
	ExpectShipLatestTime         int64     `json:"expectShipLatestTime"`         // 要求最晚发货时间
	ParentShippingTime           int64     `json:"parentShippingTime"`           // 父单发货时间
	UpdateTime                   int64     `json:"updateTime"`                   // 订单更新时间（秒级时间戳）
	LatestDeliveryTime           int64     `json:"latestDeliveryTime"`           // 最后发货时间
	FulfillmentWarning           []string  `json:"fulfillmentWarning"`           // 履约相关提醒: SUGGEST_SIGNATURE_ON_DELIVERY-建议发货时购买签名服务
	RegionId                     int       `json:"regionId"`                     // 区域 ID
	SiteId                       int       `json:"siteId"`                       // 站点 ID
	HasShippingFee               null.Bool `json:"hasShippingFee"`               // 有运费？
}

type Label struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Order struct {
	OrderSn string `json:"orderSn"` // 子订单号
	SkuId   int64  `json:"skuId"`   // SKU ID
	GoodsId int64  `json:"goodsId"` // 商品 ID
	// 备注：代表商家实际需要发货件数，在O单部分取消时：
	// 应履约件数=下单件数-发货前售后件数
	Quantity                        int    `json:"quantity"`                        // O 单应履约件数
	OriginalOrderQuantity           int    `json:"originalOrderQuantity"`           // 用户初始下单时的 O 单件数
	CanceledQuantityBeforeShipment  int    `json:"canceledQuantityBeforeShipment"`  // O 单发货前，用户发起部分取消的件数（用户申请且退款已受理）
	InventoryDeductionWarehouseId   string `json:"inventoryDeductionWarehouseId"`   // 库存扣减仓库id
	InventoryDeductionWarehouseName string `json:"inventoryDeductionWarehouseName"` // 库存扣减仓库名称
	// name: 标签名称
	// {customized_products：定制品, US_to_CA：美发加BC，is_US_to_CA_BBC：美发加BBC}
	// value
	// 是否有标签：0=无标签，1=有标签
	// BBC 订单需要结合is_US_to_CA_BBC判断
	OrderLabel  []Label `json:"orderLabel"`  // 子订单 O 单标签，内部请求异常返回为空，返回为空时请重试
	Spec        string  `json:"spec"`        // 商品信息描述
	ThumbUrl    string  `json:"thumbUrl"`    // 商品缩略图图片
	OrderStatus int     `json:"orderStatus"` // 订单状态（3 是已取消）
	// 卖家履约订单值返回：fulfillBySeller
	// 合作仓履约订单返回：fulfillByCooperativeWarehouse
	FulfillmentType    string    `json:"fulfillmentType"`    // 子订单履约类型
	GoodsName          string    `json:"goodsName"`          // 商品名称
	ProductList        []Product `json:"productList"`        // 货品信息
	FulfillmentWarning []string  `json:"fulfillmentWarning"` // 履约相关提醒: SUGGEST_SIGNATURE_ON_DELIVERY-建议发货时购买签名服务
}

type Product struct {
	ProductSkuId int64  `json:"productSkuId"` // 货品 SKU ID
	SoldFactor   int    `json:"soldFactor"`   // 商品和货品数量转换系数，商品数量(quantity)乘以转换系数，代表货品数量
	ProductId    int64  `json:"productId"`    // 货品 ID
	ExtCode      string `json:"extCode"`      // 货品编码
}

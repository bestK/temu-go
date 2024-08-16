package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/normal"
)

type shipOrderPackingService service

type ShipOrderPackingSendRequest struct {
	normal.Parameter
	DeliverMethod       int64    `json:"deliverMethod"`       // 发货方式
	DeliveryAddressId   int64    `json:"deliveryAddressId"`   // 发货地址 ID
	DeliveryOrderSnList []string `json:"deliveryOrderSnList"` // 发货单号
	SelfDeliveryInfo    struct {
		DriverUid             int    `json:"driverUid"`             // 司机uid
		DriverName            string `json:"driverName"`            // 司机姓名
		PlateNumber           string `json:"plateNumber"`           // 车牌号
		DeliveryContactNumber string `json:"deliveryContactNumber"` // 电话号码
		DeliveryContactAreaNo string `json:"deliveryContactAreaNo"` // 电话区号
		ExpressPackageNum     int    `json:"expressPackageNum"`     // 发货总箱数
	} `json:"selfDeliveryInfo"` // 自送信息
	ThirdPartyDeliveryInfo struct {
		ExpressCompanyId          int64   `json:"expressCompanyId"`          // 快递公司Id
		TmsChannelId              int64   `json:"tmsChannelId"`              // TMS快递产品类型ID
		ExpressCompanyName        string  `json:"expressCompanyName"`        // 快递公司名称
		StandbyExpress            bool    `json:"standbyExpress"`            // 是否是备用快递公司
		ExpressDeliverySn         string  `json:"expressDeliverySn"`         // 快递单号
		PredictTotalPackageWeight int64   `json:"predictTotalPackageWeight"` // 预估总包裹重量不能为空,单位克.总量必须大于等于1千克且为整千克值
		ExpectPickUpGoodsTime     int64   `json:"expectPickUpGoodsTime"`     // 预约取货时间
		ExpressPackageNum         int     `json:"expressPackageNum"`         // 交接给快递公司的包裹数
		MinChargeAmount           float64 `json:"minChargeAmount"`           // 最小预估运费（单位元）
		MaxChargeAmount           float64 `json:"maxChargeAmount"`           // 最大预估运费（单位元）
		PredictId                 int64   `json:"predictId"`                 // 预测ID
	} `json:"thirdPartyDeliveryInfo"` // 公司指定物流
	ThirdPartyExpressDeliveryInfoVO struct {
		ExpressCompanyId   int64  `json:"expressCompanyId"`   // 快递公司Id
		ExpressCompanyName string `json:"expressCompanyName"` // 快递公司名称
		ExpressDeliverySn  string `json:"expressDeliverySn"`  // 快递单号
		ExpressPackageNum  int    `json:"expressPackageNum"`  // 发货总箱数
	} `json:"thirdPartyExpressDeliveryInfoVO"`
}

func (m ShipOrderPackingSendRequest) Validate() error {
	return nil
	// return validation.ValidateStruct(&m,
	// 	validation.Field(&m.Request, validation.When(m.Request != nil, validation.By(func(value interface{}) error {
	//
	// 		return nil
	// 	}))),
	// )
}

// Send 装箱发货接口
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#ezXrHy
func (s shipOrderPackingService) Send(request ShipOrderPackingSendRequest) (number string, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			ExpressBatchSn string `json:"expressBatchSn"` // 创建生成的发货批次号
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(request).SetResult(&result).Post("bg.shiporder.packing.send")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	number = result.Result.ExpressBatchSn
	return
}

// 装箱发货校验（bg.shiporder.packing.match）

type ShipOrderPackingMatchRequest struct {
	normal.Parameter
	DeliveryOrderSnList []string `json:"deliveryOrderSnList"` // 发货单号
}

func (m ShipOrderPackingMatchRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSnList, validation.Required.Error("发货单号列表不能为空。")),
	)
}

func (s shipOrderPackingService) Match(request ShipOrderPackingMatchRequest) (items []any, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result []any `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(request).SetResult(&result).Post("bg.shiporder.packing.match")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	items = result.Result
	return
}

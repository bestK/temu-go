package temu

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
	"time"
)

// 装箱发货
type shipOrderPackingService service

// ShipOrderPackingSendSelfDeliveryInformation 自行配送信息
type ShipOrderPackingSendSelfDeliveryInformation struct {
	DriverUid             int    `json:"driverUid,omitempty"`             // 司机 uid
	DriverName            string `json:"driverName,omitempty"`            // 司机姓名
	PlateNumber           string `json:"plateNumber,omitempty"`           // 车牌号
	DeliveryContactNumber string `json:"deliveryContactNumber,omitempty"` // 电话号码
	DeliveryContactAreaNo string `json:"deliveryContactAreaNo,omitempty"` // 电话区号
	ExpressPackageNum     int    `json:"expressPackageNum,omitempty"`     // 发货总箱数
}

func (m ShipOrderPackingSendSelfDeliveryInformation) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryContactNumber,
			validation.When(len(m.DeliveryContactNumber) != 0, validation.By(is.MobilePhoneOrTelNumber())),
		),
		validation.Field(&m.DeliveryContactAreaNo,
			validation.When(len(m.DeliveryContactAreaNo) != 0, validation.By(is.TelNumberAreaCode())),
		),
		validation.Field(&m.ExpressPackageNum, validation.Min(1).Error("发货总箱数不能小于 {.min}")),
	)
}

// ShipOrderPackingSendPlatformRecommendationDeliveryInformation 平台推荐服务商配送信息
type ShipOrderPackingSendPlatformRecommendationDeliveryInformation struct {
	ExpressCompanyId          int64   `json:"expressCompanyId,omitempty"`          // 快递公司 Id
	TmsChannelId              int     `json:"tmsChannelId,omitempty"`              // TMS 快递产品类型 Id
	ExpressCompanyName        string  `json:"expressCompanyName,omitempty"`        // 快递公司名称
	StandbyExpress            bool    `json:"standbyExpress"`                      // 是否是备用快递公司
	ExpressDeliverySn         string  `json:"expressDeliverySn,omitempty"`         // 快递单号
	PredictTotalPackageWeight int64   `json:"predictTotalPackageWeight,omitempty"` // 预估总包裹重量不能为空,单位克.总量必须大于等于1千克且为整千克值
	ExpectPickUpGoodsTime     int64   `json:"expectPickUpGoodsTime,omitempty"`     // 预约取货时间
	ExpressPackageNum         int     `json:"expressPackageNum,omitempty"`         // 交接给快递公司的包裹数
	MinChargeAmount           float64 `json:"minChargeAmount,omitempty"`           // 最小预估运费（单位元）
	MaxChargeAmount           float64 `json:"maxChargeAmount,omitempty"`           // 最大预估运费（单位元）
	PredictId                 int64   `json:"predictId,omitempty"`                 // 预测 ID
}

func (m ShipOrderPackingSendPlatformRecommendationDeliveryInformation) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ExpressCompanyId, validation.Required.Error("快递公司 Id 不能为空")),
		validation.Field(&m.ExpressCompanyName, validation.Required.Error("快递公司名称不能为空")),
		validation.Field(&m.ExpressDeliverySn, validation.Required.Error("快递单号不能为空")),
		validation.Field(&m.PredictTotalPackageWeight, validation.Min(1).Error("预估总包裹重量不能小于 {.min}")),
		validation.Field(&m.ExpectPickUpGoodsTime,
			validation.Required.Error("预约取货时间不能为空"),
			validation.By(func(value interface{}) error {
				v, ok := value.(int64)
				if !ok {
					return fmt.Errorf("无效的预约取货时间：%v", value)
				}

				t := time.UnixMilli(v)
				if t.IsZero() || t.Before(time.Now()) {
					return fmt.Errorf("无效的预约取货时间：%v", value)
				}

				return nil
			}),
		),
		validation.Field(&m.ExpressPackageNum, validation.Min(1).Error("发货总箱数不能小于 {.min}")),
		validation.Field(&m.PredictId, validation.Required.Error("预测 Id 不能为空")),
	)
}

// ShipOrderPackingSendThirdPartyDeliveryInformation 自行委托第三方物流配送信息
type ShipOrderPackingSendThirdPartyDeliveryInformation struct {
	ExpressCompanyId   int64  `json:"expressCompanyId"`            // 快递公司 Id
	ExpressCompanyName string `json:"expressCompanyName"`          // 快递公司名称
	ExpressDeliverySn  string `json:"expressDeliverySn"`           // 快递单号
	ExpressPackageNum  int    `json:"expressPackageNum,omitempty"` // 发货总箱数
}

func (m ShipOrderPackingSendThirdPartyDeliveryInformation) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ExpressCompanyId, validation.Required.Error("快递公司 Id 不能为空")),
		validation.Field(&m.ExpressCompanyName, validation.Required.Error("快递公司名称不能为空")),
		validation.Field(&m.ExpressDeliverySn, validation.Required.Error("快递单号不能为空")),
		validation.Field(&m.ExpressPackageNum, validation.Min(1).Error("发货总箱数不能小于 {.min}")),
	)
}

type ShipOrderPackingSendRequest struct {
	normal.Parameter
	DeliverMethod                   null.Int                                                       `json:"deliverMethod"`                             // 发货方式
	DeliveryAddressId               int64                                                          `json:"deliveryAddressId"`                         // 发货地址 ID
	DeliveryOrderSnList             []string                                                       `json:"deliveryOrderSnList"`                       // 发货单号
	SelfDeliveryInfo                *ShipOrderPackingSendSelfDeliveryInformation                   `json:"selfDeliveryInfo,omitempty"`                // 自送信息
	ThirdPartyDeliveryInfo          *ShipOrderPackingSendPlatformRecommendationDeliveryInformation `json:"thirdPartyDeliveryInfo,omitempty"`          // 公司指定物流
	ThirdPartyExpressDeliveryInfoVO *ShipOrderPackingSendThirdPartyDeliveryInformation             `json:"thirdPartyExpressDeliveryInfoVO,omitempty"` // 第三方配送
}

func (m ShipOrderPackingSendRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliverMethod,
			validation.By(func(value interface{}) error {
				v, ok := value.(null.Int)
				if !ok || !v.Valid {
					return errors.New("无效的发货方式")
				}

				err := validation.Validate(int(v.Int64), validation.In(
					entity.DeliveryMethodSelf,
					entity.DeliveryMethodPlatformRecommendation,
					entity.DeliveryMethodThirdParty,
				).Error("无效的发货方式"))
				if err != nil {
					return err
				}

				switch v.Int64 {
				case entity.DeliveryMethodSelf:
					if m.SelfDeliveryInfo == nil {
						return errors.New("装箱发货自送信息不能为空")
					} else {
						return m.SelfDeliveryInfo.validate()
					}
				case entity.DeliveryMethodPlatformRecommendation:
					if m.ThirdPartyDeliveryInfo == nil {
						return errors.New("装箱发货物流信息不能为空")
					} else {
						return m.ThirdPartyDeliveryInfo.validate()
					}
				case entity.DeliveryMethodThirdParty:
					if m.ThirdPartyExpressDeliveryInfoVO == nil {
						return errors.New("装箱发货第三方配送信息不能为空")
					} else {
						return m.ThirdPartyExpressDeliveryInfoVO.validate()
					}
				}
				return nil
			}),
		),
		validation.Field(&m.DeliveryAddressId, validation.Required.Error("发货地址不能为空")),
		validation.Field(&m.DeliveryOrderSnList,
			validation.Required.Error("发货单号不能为空"),
			validation.Each(validation.By(is.ShipOrderNumber())),
		),
	)
}

type ShipOrderPackingSendCreateResult struct {
	DeliveryOrderCreateInfos []struct {
		SubPurchaseOrderSn string `json:"subPurchaseOrderSn"` // 采购子单号
		PackageInfos       []struct {
			PackageDetailSaveInfos []struct {
				ProductSkuId int64 `json:"productSkuId"` // skuId
				SkuNum       int   `json:"skuNum"`       // 发货 sku 数目
			} `json:"packageDetailSaveInfos"` // 包裹明细
		} `json:"packageInfos"` // 包裹信息列表
	} `json:"deliveryOrderCreateInfos"` // 采购单创建信息列表
	DeliverOrderDetailInfos []struct {
		ProductSkuId  int64 `json:"productSkuId"`  // skuId
		DeliverSkuNum int   `json:"deliverSkuNum"` // 发货 sku 数目
	} `json:"deliverOrderDetailInfos"` // 发货单详情列表
	// DeliveryAddressId  int64                  `json:"deliveryAddressId"`  // 发货地址 ID
	ReceiveAddressInfo              *entity.ReceiveAddress                                         `json:"receiveAddressInfo"`              // 收货地址
	SubWarehouseId                  int64                                                          `json:"subWarehouseId"`                  // 子仓 ID
	DeliverMethod                   int                                                            `json:"deliverMethod"`                   // 发货方式
	DeliveryAddressId               int64                                                          `json:"deliveryAddressId"`               // 发货地址 ID
	SelfDeliveryInfo                *ShipOrderPackingSendSelfDeliveryInformation                   `json:"selfDeliveryInfo"`                // 自送信息
	ThirdPartyDeliveryInfo          *ShipOrderPackingSendPlatformRecommendationDeliveryInformation `json:"thirdPartyDeliveryInfo"`          // 公司指定物流
	ThirdPartyExpressDeliveryInfoVO *ShipOrderPackingSendThirdPartyDeliveryInformation             `json:"thirdPartyExpressDeliveryInfoVO"` // 第三方配送
}

type ShipOrderPackingSendResult struct {
	CreateExpressErrorRequestList []ShipOrderPackingSendCreateResult `json:"createExpressErrorRequestList"` // 创建快递运单失败的请求列表
	ExpressBatchSn                string                             `json:"expressBatchSn"`                // 创建生成的发货批次号
}

// Send 装箱发货接口
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#ezXrHy
func (s shipOrderPackingService) Send(ctx context.Context, request ShipOrderPackingSendRequest) (number string, err error) {
	if err = request.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result ShipOrderPackingSendResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.shiporder.packing.send")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	number = result.Result.ExpressBatchSn
	return
}

// 装箱发货校验

type ShipOrderPackingMatchRequest struct {
	normal.Parameter
	DeliveryOrderSnList []string `json:"deliveryOrderSnList"` // 发货单号
}

func (m ShipOrderPackingMatchRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSnList,
			validation.Required.Error("发货单号列表不能为空"),
			validation.Each(validation.By(is.ShipOrderNumber())),
		),
	)
}

// Match 装箱发货校验
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#TDP3qU
func (s shipOrderPackingService) Match(ctx context.Context, request ShipOrderPackingMatchRequest) (item entity.ShipOrderPackingMatchResult, err error) {
	if err = request.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result entity.ShipOrderPackingMatchResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.shiporder.packing.match")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	item = result.Result
	return
}

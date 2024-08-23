package temu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var deliveryOrderSnList = []string{"FH2408151907154", "FH2408161777912"}

func Test_barcodeService_BoxMarkPrintUrl(t *testing.T) {
	url, err := temuClient.Services.Barcode.BoxMarkPrintUrl(deliveryOrderSnList...)
	assert.Equal(t, nil, err, "Services.Barcode.BoxMarkPrintUrl(%#v)", deliveryOrderSnList)
	assert.Contains(t, url, "https://openapi.kuajingmaihuo.com/tool/print?dataKey=")
}

func Test_barcodeService_BoxMark(t *testing.T) {
	_, err := temuClient.Services.Barcode.BoxMark(deliveryOrderSnList...)
	assert.Equal(t, nil, err, "Services.Barcode.BoxMark(%#v)", deliveryOrderSnList)
}

func Test_barcodeService_NormalGoods(t *testing.T) {
	params := NormalGoodsBarcodeQueryParams{
		ProductSkcIdList: []int{8972250969},
	}
	_, err := temuClient.Services.Barcode.NormalGoods(params)
	assert.Equalf(t, nil, err, "Services.Barcode.NormalGoods(%#v)", params)
}

func Test_barcodeService_CustomGoods(t *testing.T) {
	params := CustomGoodsBarcodeQueryParams{
		PersonalProductSkuIdList: []int{60294097402138},
	}
	items, err := temuClient.Services.Barcode.CustomGoods(params)
	assert.Equalf(t, nil, err, "Services.Barcode.CustomGoods(%#v)", params)
	assert.Equal(t, true, len(items) != 0, "Services.Barcode.CustomGoods(%#v)", params)
}

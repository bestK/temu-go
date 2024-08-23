package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsSalesService_BoxMark(t *testing.T) {
	params := GoodsSalesQueryParams{}
	items, err := temuClient.Services.GoodsSales.All(params)
	assert.Equal(t, nil, err, "Services.GoodsSales.All")

	if len(items) != 0 {
		item := items[0]
		var sales entity.GoodsSales
		// 根据商品 SKC ID 查询
		sales, err = temuClient.Services.GoodsSales.One(item.ProductSkcID)
		assert.Equalf(t, nil, err, "Services.PurchaseOrder.One(%d)", item.ProductSkcID)
		assert.Equalf(t, item, sales, "Services.PurchaseOrder.One(%d)", item.ProductSkcID)
	}
}

package temu

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var temuClient *Client
var ctx context.Context

func TestMain(m *testing.M) {
	b, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var cfg config.Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	temuClient = NewClient(cfg)
	ctx = context.Background()
	m.Run()
}

func TestClient_SetRegionId(t *testing.T) {
	tests := []struct {
		name           string
		regionId       int
		expectRegionId int
	}{
		{"t1", 1, 0},
		{"t2", entity.AmericanRegionId, entity.AmericanRegionId},
	}
	for _, tt := range tests {
		temuClient.SetRegionId(tt.regionId)
		assert.Equalf(t, tt.expectRegionId, temuClient.RegionId, tt.name)
	}
}

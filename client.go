package temu

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/hiscaler/gox/stringx"
	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/normal"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	Version   = "0.0.1"
	UserAgent = "Temu API Client-Golang/" + Version
)

const (
	BadRequestError           = 400     // 错误的请求
	UnauthorizedError         = 401     // 身份验证或权限错误
	NotFoundError             = 404     // 访问资源不存在
	InternalServerError       = 500     // 服务器内部错误
	MethodNotImplementedError = 501     // 方法未实现
	SystemExceptionError      = 200000  // 系统异常
	InvalidSignError          = 7000015 // 签名无效
	NoAppKeyError             = 7000002 // 未设置 App Key
	NoAccessTokenError        = 7000003 // 未设置 Access Token
	InvalidAccessTokenError   = 7000018 // 无效的 Access Token
	AccessTokenKeyUnmatched   = 7000006 // Access Token 和 Key 不匹配
)

var ErrNotFound = errors.New("数据不存在")
var ErrInvalidSign = errors.New("无效的签名")
var ErrInvalidParameters = errors.New("无效的参数")

type service struct {
	debug      bool          // Is debug mode
	logger     *slog.Logger  // Log
	httpClient *resty.Client // HTTP client
}

type services struct {
	ShipOrder               shipOrderService
	ShipOrderStaging        shipOrderStagingService
	ShipOrderPacking        shipOrderPackingService
	ShipOrderPackage        shipOrderPackageService
	PurchaseOrder           purchaseOrderService
	Logistics               logisticsService
	GoodsSizeChart          goodsSizeChartService
	GoodsSizeChartClass     goodsSizeChartClassService
	GoodsSizeChartSetting   goodsSizeChartSettingService
	MallAddress             mallAddressService
	ShipOrderReceiveAddress shipOrderReceiveAddressService
	Goods                   goodsService
	Mall                    mallService
	JitVirtualInventory     jitVirtualInventoryService
	JitMode                 jitModeService
	JitPresaleRule          jitPresaleRuleService
}

type Client struct {
	Debug        bool           // Is debug mode
	Logger       *slog.Logger   // Log
	Services     services       // API services
	TimeLocation *time.Location // 时区
}

// generate sign string
func generateSign(values map[string]any, appSecret string) map[string]any {
	delete(values, "sign")
	values["timestamp"] = time.Now().Unix()
	size := len(values)
	keys := make([]string, size)
	for k := range values {
		size--
		keys[size] = k
	}
	sort.Strings(keys)
	sb := strings.Builder{}
	sb.WriteString(appSecret)
	for _, key := range keys {
		value := stringx.String(values[key])
		if value == "" {
			delete(values, key)
			continue
		}
		sb.WriteString(key)
		sb.WriteString(value)
	}
	sb.WriteString(appSecret)
	values["sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
	return values
}

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	time.Local = loc
}

func NewClient(config config.Config) *Client {
	var l *slog.Logger
	if config.Logger != nil {
		l = config.Logger
	} else {
		if config.Debug {
			l = slog.New(slog.NewTextHandler(os.Stdout, nil))
		} else {
			l = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		}
	}
	client := &Client{
		Debug:        config.Debug,
		Logger:       l,
		TimeLocation: loc,
	}
	httpClient := resty.New().
		SetDebug(config.Debug).
		EnableTrace().
		SetBaseURL("https://openapi.kuajingmaihuo.com/openapi/router").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   UserAgent,
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(config.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: config.Timeout * time.Second,
			}).DialContext,
		}).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			values := make(map[string]any)
			if request.Body != nil {
				b, e := json.Marshal(request.Body)
				if e != nil {
					return e
				}

				e = json.Unmarshal(b, &values)
				if e != nil {
					return e
				}
			}
			values["app_key"] = config.AppKey
			values["app_secret"] = config.AppSecret
			values["access_token"] = config.AccessToken
			values["data_type"] = "JSON"
			values["version"] = "V1"
			values["timestamp"] = time.Now().Unix()
			values["type"] = request.URL
			request.URL = ""
			request.SetBody(generateSign(values, config.AppSecret))
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
			statusCode := response.StatusCode()
			if statusCode != http.StatusOK && statusCode != http.StatusCreated {
				params := response.Request.Body
				endpoint := ""
				if v, ok := params.(map[string]any); ok {
					if vv, ok := v["type"]; ok {
						endpoint = vv.(string)
					}
				}
				l.Error(fmt.Sprintf(`%s %s
   ENDPOINT: %s
 PARAMETERS: %s
STATUS CODE: %d
       BODY: %s
`,
					response.Request.Method,
					response.Request.URL,
					endpoint,
					params,
					statusCode,
					response.String(),
				))
			}
			return nil
		}).
		SetRetryCount(3).
		SetRetryWaitTime(time.Duration(500) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(1) * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if response == nil {
				return false
			}

			retry := response.StatusCode() == http.StatusTooManyRequests
			if !retry {
				r := struct {
					Success   bool   `json:"success"`
					ErrorCode int    `json:"errorCode"`
					ErrorMsg  string `json:"errorMsg"`
				}{}
				retry = json.Unmarshal(response.Body(), &r) == nil &&
					r.Success == false &&
					r.ErrorCode == 4000000 &&
					strings.EqualFold(r.ErrorMsg, "SYSTEM_EXCEPTION")
			}
			if retry {
				body := response.Request.Body
				endpoint := ""
				if body != nil {
					var values map[string]any
					var b []byte
					var e error
					if b, e = json.Marshal(body); e == nil {
						if e = json.Unmarshal(b, &values); e == nil {
							if v, ok := values["type"]; ok {
								endpoint = v.(string)
							}
							response.Request.SetBody(generateSign(values, config.AppSecret))
						}
					}
					retry = e == nil
				}
				if retry {
					messages := make([]string, 0)
					messages = append(messages, "URL: "+response.Request.URL)
					if endpoint != "" {
						messages = append(messages, "Type: "+endpoint)
					}
					l.Info("Retry " + strings.Join(messages, " "))
				}
			}
			return retry
		})
	if config.Debug {
		httpClient.SetBaseURL("https://openapi.kuajingmaihuo.com/openapi/router")
	}
	httpClient.JSONMarshal = json.Marshal
	httpClient.JSONUnmarshal = json.Unmarshal
	xService := service{
		debug:      config.Debug,
		logger:     l,
		httpClient: httpClient,
	}
	client.Services = services{
		ShipOrder:        (shipOrderService)(xService),
		ShipOrderStaging: (shipOrderStagingService)(xService),
		ShipOrderPacking: (shipOrderPackingService)(xService),
		ShipOrderPackage: (shipOrderPackageService)(xService),
		PurchaseOrder:    (purchaseOrderService)(xService),
		Goods: goodsService{
			service:       xService,
			Category:      (goodsCategoryService)(xService),
			Brand:         (goodsBrandService)(xService),
			LifeCycle:     (goodsLifeCycleService)(xService),
			Sales:         (goodsSalesService)(xService),
			TopSelling:    (goodsTopSellingService)(xService),
			Certification: (goodsCertificationService)(xService),
			Warehouse:     (goodsWarehouseService)(xService),
			Barcode:       (goodsBarcodeService)(xService),
		},
		Logistics:               (logisticsService)(xService),
		GoodsSizeChart:          (goodsSizeChartService)(xService),
		GoodsSizeChartClass:     (goodsSizeChartClassService)(xService),
		GoodsSizeChartSetting:   (goodsSizeChartSettingService)(xService),
		MallAddress:             (mallAddressService)(xService),
		ShipOrderReceiveAddress: (shipOrderReceiveAddressService)(xService),
		Mall:                    (mallService)(xService),
		JitVirtualInventory:     (jitVirtualInventoryService)(xService),
		JitMode:                 (jitModeService)(xService),
		JitPresaleRule:          (jitPresaleRuleService)(xService),
	}

	return client
}

func parseResponseTotal(currentPage, pageSize, total int) (n, totalPages int, isLastPage bool) {
	if currentPage == 0 {
		currentPage = 1
	}

	totalPages = (total / pageSize) + 1
	return total, totalPages, currentPage >= totalPages
}

func invalidInput(e error) error {
	var errs validation.Errors
	if !errors.As(e, &errs) {
		return e
	}

	size := len(errs)
	if size == 0 {
		return errors.New("未知错误")
	}

	fields := make([]string, size)
	messages := make([]string, size)
	for field := range errs {
		size--
		fields[size] = field
	}
	sort.Strings(fields)
	for i, field := range fields {
		messages[i] = errs[field].Error()
	}
	return errors.New(strings.Join(messages, ". "))
}

func recheckError(resp *resty.Response, result normal.Response, e error) (err error) {
	if e != nil {
		return e
	}

	if resp.IsError() {
		errorMessage := strings.TrimSpace(result.ErrorMessage)
		if errorMessage == "" {
			return errorWrap(resp.StatusCode(), resp.Error().(string))
		}
		return errors.New(errorMessage)
	}

	if !result.Success {
		return errorWrap(result.ErrorCode, result.ErrorMessage)
	}
	return nil
}

// errorWrap wrap an error, if status code is 200, return nil, otherwise return an error
func errorWrap(code int, message string) error {
	if code == http.StatusOK {
		return nil
	}

	if code == NotFoundError {
		return ErrNotFound
	}

	message = strings.TrimSpace(message)
	switch code {
	case BadRequestError:
		message = "请求错误"
	case UnauthorizedError:
		message = "认证失败，请确认您是否有相应的权限"
	case InternalServerError:
		message = "服务器内部错误"
	case MethodNotImplementedError:
		message = "请求方法未实现"
	case SystemExceptionError, 4000000:
		message = "Temu 平台异常"
	case InvalidSignError:
		return ErrInvalidSign
	case NoAppKeyError:
		message = "未设置 App Key"
	case NoAccessTokenError:
		message = "未设置 Access Token"
	case InvalidAccessTokenError, 7000020:
		message = "无效的 Access Token"
	case AccessTokenKeyUnmatched:
		message = "Access Token 和 Key 不匹配"
	case 7000016:
		return errors.New("无效的请求地址")
	case 2000000, 2000090, 3000000:
		return errors.New(message)
	default:
		message = fmt.Sprintf("%d: %s", code, message)
	}

	return errors.New(message)
}

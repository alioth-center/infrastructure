package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alioth-center/infrastructure/cli"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/trace"
	vs "github.com/alioth-center/infrastructure/utils/values"
	"github.com/atotto/clipboard"
	"github.com/olekukonko/tablewriter"
)

func initCommands() {
	cli.RegisterHandler("GetGeoIP", GetGeoIP)
	cli.RegisterHandler("FormatJSON", JsonFormat)
	cli.RegisterHandler("UnixTime", UnixTime)
}

func toString(value any) string {
	return fmt.Sprintf("%v", value)
}

type values struct {
	key   string
	value string
}

func printResult(values ...values) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoFormatHeaders(true)
	table.SetAutoMergeCells(true)
	table.SetHeader([]string{"Propriety", "Value"})
	for _, value := range values {
		table.Append([]string{value.key, value.value})
	}
	table.Render()
}

func GetGeoIP(input *cli.Input) {
	ip := input.Params.GetString("ip")
	if ip == "" {
		// no ip args, get localhost ip
		resp, err := http.NewSimpleClient().ExecuteRequest(http.NewRequestBuilder().
			WithContext(trace.NewContext()).
			WithMethod(http.GET).
			WithPath("https://api.ip.sb/ip").
			WithHeader(http.HeaderUserAgent, "Mozilla/5.0"),
		)
		if err != nil {
			printResult(values{"error", err.Error()})
			return
		}

		ip = strings.TrimSpace(string(resp.RawBody()))
	}
	// get ip info
	resp, err := http.NewSimpleClient().ExecuteRequest(http.NewRequestBuilder().
		WithContext(trace.NewContext()).
		WithMethod(http.GET).
		WithPath("https://api.ip.sb/geoip/"+ip).
		WithHeader(http.HeaderUserAgent, "Mozilla/5.0"),
	)
	if err != nil {
		printResult(values{"error", err.Error()})
		return
	}

	result := map[string]any{}
	if bindErr := resp.BindJson(&result); bindErr != nil {
		printResult(values{"error", bindErr.Error()})
		return
	}

	resultValues := []values{
		{
			key:   "ip",
			value: toString(result["ip"]),
		},
		{
			key:   "longitude",
			value: toString(result["longitude"]),
		},
		{
			key:   "latitude",
			value: toString(result["latitude"]),
		},
		{
			key:   "continent",
			value: toString(result["continent_code"]),
		},
		{
			key:   "country",
			value: fmt.Sprintf("[%v] %v", result["country_code"], result["country"]),
		},
		{
			key:   "city",
			value: fmt.Sprintf("[%v] %v", result["timezone"], result["city"]),
		},
		{
			key:   "isp",
			value: fmt.Sprintf("[%v] %v", result["asn"], result["isp"]),
		},
		{
			key:   "organization",
			value: toString(result["organization"]),
		},
		{
			key:   "asn_organization",
			value: toString(result["asn_organization"]),
		},
	}

	printResult(resultValues...)
}

// unescapeJSONString 递归解码被多次转义的 JSON 字符串
func unescapeJSONString(input string) (string, error) {
	decoded := input
	var err error

	// 尝试递归解码直到不能再解码为止
	for {
		var temp string
		err = json.Unmarshal([]byte(decoded), &temp)
		if err != nil {
			break
		}
		decoded = temp
	}

	return decoded, nil
}

// parseNestedJSON 递归解析嵌套的 JSON 字符串
func parseNestedJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		decoded, err := unescapeJSONString(v)
		if err == nil {
			var nestedData interface{}
			if json.Unmarshal([]byte(decoded), &nestedData) == nil {
				return parseNestedJSON(nestedData)
			}
		}
		return v
	case map[string]interface{}:
		for key, value := range v {
			v[key] = parseNestedJSON(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = parseNestedJSON(value)
		}
		return v
	default:
		return v
	}
}

// prettyPrintJSON 格式化打印 JSON
func prettyPrintJSON(input string) error {
	var jsonData interface{}

	// 先尝试将解码后的字符串解析为 JSON
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		return fmt.Errorf("cannot parse JSON: %v", err)
	}

	// 递归解析嵌套的 JSON 字符串
	jsonData = parseNestedJSON(jsonData)

	// 以缩进格式输出 JSON
	prettyJSON, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return fmt.Errorf("cannot format JSON: %v", err)
	}

	fmt.Println(string(prettyJSON))
	return nil
}

func JsonFormat(_ *cli.Input) {
	// 输入的字符串
	raw, readErr := clipboard.ReadAll()
	if readErr != nil {
		fmt.Println("Error reading clipboard:", readErr)
		return
	}

	// 第零步：输入可能不规整，处理
	invalid := false
	if strings.HasPrefix(raw, `{\"`) {
		invalid = true
		raw = `"` + raw
	}
	if invalid && !strings.HasSuffix(raw, `}"`) {
		raw = raw + `"`
	}

	// 第一步：解码被多次转义的 JSON 字符串
	decodedInput, err := unescapeJSONString(raw)
	if err != nil {
		log.Fatalf("无法解码输入: %v", err)
	}

	// 第二步：格式化并打印 JSON
	err = prettyPrintJSON(decodedInput)
	if err != nil {
		log.Fatalf("无法打印格式化 JSON: %v", err)
	}
}

func UnixTime(input *cli.Input) {
	ts := vs.StringToInt(input.Params.GetString("timestamp"), int64(0))
	if ts == 0 {
		printResult(values{"error", "invalid timestamp"})
		return
	}

	// check millisecond or second
	if ts < 1000000000000 {
		ts *= 1000
	}

	timestamp := time.UnixMilli(ts)
	printResult(values{"timestamp", fmt.Sprintf("%v", timestamp.Format("2006-01-02 15:04:05.000-07"))})
}

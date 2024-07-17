package domain

import "time"

type (
	AlertTime struct {
		time.Time
	}

	AlertInfo struct {
		Action      string `json:"action" mapstructure:"action"`
		CaIDA       string `json:"ca_id_a" mapstructure:"ca_id_a"`
		CaIDU       string `json:"ca_id_u" mapstructure:"ca_id_u"`
		Category    string `json:"category" mapstructure:"category"`
		CmGgA       int    `json:"cm_cg_a" mapstructure:"cm_cg_a"`
		CmCgU       int    `json:"cm_cg_u" mapstructure:"cm_cg_u"`
		CmID        string `json:"cm_id" mapstructure:"cm_id"`
		CmIDA       string `json:"cm_id_a" mapstructure:"cm_id_a"`
		CmIDU       string `json:"cm_id_u" mapstructure:"cm_id_u"`
		CmIDU2      string `json:"cm_id_u2" mapstructure:"cm_id_u2"`
		CmLvA       int    `json:"cm_lv_a" mapstructure:"cm_lv_a"`
		CmLvU       int    `json:"cm_lv_u" mapstructure:"cm_lv_u"`
		GID         int    `json:"gid" mapstructure:"gid"`
		QrID        string `json:"qr_id" mapstructure:"qr_id"`
		Rev         int    `json:"rev" mapstructure:"rev"`
		Severity    int    `json:"severity" mapstructure:"severity"`
		Signature   string `json:"signature" mapstructure:"signature"`
		SignatureID int    `json:"signature_id" mapstructure:"signature_id"`
		TxID        int    `json:"tx_id" mapstructure:"tx_id"`
		WpID        string `json:"wp_id" mapstructure:"wp_id"`
	}

	Location struct {
		City       string `json:"city" mapstructure:"city"`
		CityEn     string `json:"city_en" mapstructure:"city_en"`
		Country    string `json:"country" mapstructure:"country"`
		CountryEn  string `json:"country_en" mapstructure:"country_en"`
		Latitude   int    `json:"latitude" mapstructure:"latitude"`
		Longitude  int    `json:"longitude" mapstructure:"longitude"`
		Province   string `json:"province" mapstructure:"province"`
		ProvinceEn string `json:"province_en" mapstructure:"province_en"`
	}

	HTTPInfo struct {
		AnalyseIP              string `json:"analyse_ip" mapstructure:"analyse_ip"`
		AppName                string `json:"appname" mapstructure:"appname"`
		DestIPURLPathMd5       string `json:"dest_ip_url_path_md5" mapstructure:"dest_ip_url_path_md5"`
		DomaiName              string `json:"domainame" mapstructure:"domainame"`
		DType                  int    `json:"dtype" mapstructure:"dtype"`
		ExpandField            string `json:"expand_field" mapstructure:"expand_field"`
		HostName               string `json:"hostname" mapstructure:"hostname"`
		HTTPContentType        string `json:"http_content_type" mapstructure:"http_content_type"`
		HTTPMethod             string `json:"http_method" mapstructure:"http_method"`
		HTTPRequestHeaders     string `json:"http_request_headers" mapstructure:"http_request_headers"`
		HTTPRequestHeadersLen  int    `json:"http_request_headers_len" mapstructure:"http_request_headers_len"`
		HTTPResponseHeaders    string `json:"http_response_headers" mapstructure:"http_response_headers"`
		HTTPResponseHeadersLen int    `json:"http_response_headers_len" mapstructure:"http_response_headers_len"`
		HTTPUserAgent          string `json:"http_user_agent" mapstructure:"http_user_agent"`
		Length                 int    `json:"length" mapstructure:"length"`
		OpenService            string `json:"open_service" mapstructure:"open_service"`
		Protocol               string `json:"protocol" mapstructure:"protocol"`
		Redirect               string `json:"redirect" mapstructure:"redirect"`
		Status                 int    `json:"status" mapstructure:"status"`
		Type                   string `json:"type" mapstructure:"type"`
		URL                    string `json:"url" mapstructure:"url"`
		URLPath                string `json:"urlpath" mapstructure:"urlpath"`
		URLPathMd5             string `json:"urlpath_md5" mapstructure:"urlpath_md5"`
		UtID                   string `json:"utid" mapstructure:"utid"`
	}

	Monitor1 struct {
		Alert         AlertInfo `json:"alert" mapstructure:"alert"`
		AssetRole     int       `json:"asset_role" mapstructure:"asset_role"`
		ClusterID     string    `json:"clusterid" mapstructure:"clusterid"`
		Content       string    `json:"content" mapstructure:"content"`
		CreatedAt     string    `json:"created_at" mapstructure:"created_at"`
		DestAsset     int       `json:"dest_asset" mapstructure:"dest_asset"`
		DestIP        string    `json:"dest_ip" mapstructure:"dest_ip"`
		DestIPVal     int       `json:"dest_ip_val" mapstructure:"dest_ip_val"`
		DestPort      int       `json:"dest_port" mapstructure:"dest_port"`
		DipInfo       Location  `json:"dip_info" mapstructure:"dip_info"`
		Dispose       int       `json:"dispose" mapstructure:"dispose"`
		EventTime     AlertTime `json:"event_time" mapstructure:"event_time"`
		EventType     string    `json:"event_type" mapstructure:"event_type"`
		EventUnixtime int       `json:"event_unixtime" mapstructure:"event_unixtime"`
		FlowID        int       `json:"flow_id" mapstructure:"flow_id"`
		HTTP          HTTPInfo  `json:"http" mapstructure:"http"`
		ID            string    `json:"id" mapstructure:"id"`
		InIface       string    `json:"in_iface" mapstructure:"in_iface"`
		IPType        int       `json:"ip_type" mapstructure:"ip_type"`
		PktSrc        string    `json:"pkt_src" mapstructure:"pkt_src"`
		Proto         string    `json:"proto" mapstructure:"proto"`
		PV            float64   `json:"pv" mapstructure:"pv"`
		QrPayload     string    `json:"qr_payload" mapstructure:"qr_payload"`
		SIPInfo       Location  `json:"sip_info" mapstructure:"sip_info"`
		SrcCountry    string    `json:"src_country" mapstructure:"src_country"`
		SrcIP         string    `json:"src_ip" mapstructure:"src_ip"`
		SrcIPVal      int       `json:"src_ip_val" mapstructure:"src_ip_val"`
		SrcPort       int       `json:"src_port" mapstructure:"src_port"`
		Stream        int       `json:"stream" mapstructure:"stream"`
		TxID          int       `json:"tx_id" mapstructure:"tx_id"`
		UV            float64   `json:"uv" mapstructure:"uv"`
	}
)

const (
	// customTimeLayoutTZ = "2006-01-02T15:04:05.999999-0700"
	customTimeLayout = "2006-01-02T15:04:05.999999"
)

func (ct *AlertTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	t, err := time.Parse(`"`+customTimeLayout+`"`, s)
	if err != nil {
		return err
	}

	ct.Time = t
	return nil
}

func (ct *AlertTime) Decode() {

}

type (
	Monitor struct {
		ID            string   `json:"id"`
		EventTime     string   `json:"event_time"`
		FlowID        int64    `json:"flow_id"`
		EventUnixTime int      `json:"event_unixtime"`
		AlertTime     int      `json:"alert_time"`
		EventType     string   `json:"event_type"`
		Pid           int      `json:"pid"`
		Vlan          []int    `json:"vlan"`
		SrcIP         string   `json:"src_ip"`
		SrcPort       int      `json:"src_port"`
		DestIP        string   `json:"dest_ip"`
		DestPort      int      `json:"dest_port"`
		SrcMac        string   `json:"src_mac"`
		DestMac       string   `json:"dest_mac"`
		PktSrc        string   `json:"pkt_src"`
		Proto         string   `json:"proto"`
		AppProto      string   `json:"app_proto"`
		LogType       int      `json:"log_type"`
		AppName       []string `json:"app_name"`
		PcapPath      string   `json:"pcap_path"`
		SequenceNo    int      `json:"sequence_no"`
		Clusterid     string   `json:"clusterid"`
		TxID          int      `json:"tx_id"`
		Alert         Alert    `json:"alert"`
		HTTP          HTTP     `json:"http"`
		AssetRole     int      `json:"asset_role"`
		FlowEstab     int      `json:"flow_estab"`
		PayloadHex    string   `json:"payload_hex"`
		Stream        int      `json:"stream"`
		Tag           int      `json:"tag"`
	}
	Alert struct {
		Action      string `json:"action"`
		GID         int    `json:"gid"`
		QrID        int    `json:"qr_id"`
		SignatureID int    `json:"signature_id"`
		Rev         int    `json:"rev"`
		Signature   string `json:"signature"`
		Category    string `json:"category"`
		Severity    int    `json:"severity"`
		TxID        int    `json:"tx_id"`
	}
	HTTP struct {
		Hostname               string `json:"hostname"`
		URL                    string `json:"url"`
		URLPath                string `json:"urlpath"`
		HTTPUserAgent          string `json:"http_user_agent"`
		HTTPAcceptLanguage     string `json:"http_accept_language"`
		HTTPContentType        string `json:"http_content_type"`
		HTTPMethod             string `json:"http_method"`
		Protocol               string `json:"protocol"`
		Status                 int    `json:"status"`
		Length                 int    `json:"length"`
		HTTPRequestHeaders     string `json:"http_request_headers"`
		HTTPRequestHeadersLen  int    `json:"http_request_headers_len"`
		HTTPResponseCode       int    `json:"http_response_code"`
		HTTPResponseHeaders    string `json:"http_response_headers"`
		HTTPResponseHeadersLen int    `json:"http_response_headers_len"`
	}
)

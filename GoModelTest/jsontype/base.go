package models

const (
	CONFIG_PATH  = "config"
	BIN_PATH     = "bin"
	INITIAL_PATH = "reset"
	UPGRADE_PATH = "upgrade"
	COMMAND_PATH = "shell"

	GLOBAL_CURCFG_DOC = "/config/libraT.yaml"
	BUMBLE_CURCFG_DOC = "/config/bumble.yaml"
	ONVIF_CURCFG_DOC  = "/config/onvifserver.xml"

	GLOBAL_ORGCFG_DOC = "/reset/bumble/libraT.yaml"
	BUMBLE_ORGCFG_DOC = "/reset/bumble/bumble.yaml"

	RELOADSENSOR_TEMPLATE_DOC = "/bin/bumble/template/reset_sensor.template"
	RESETNETWORK_TEMPLATE_DOC = "/bin/bumble/template/set_static_ip.template"
	RELOADSENSOR_SHELL_DOC    = "/tmp/reset_sensor.sh"
	RESETNETWORK_SHELL_DOC    = "/tmp/set_static_ip.template"

	// config about machine,in libraT.yaml
	SENSOR_SN    = "/libraT/sensor_sn"
	SENSOR_UID   = "/libraT/sensor_uid"
	SENSOR_DESC  = "/libraT/sensor_desc"
	SERVICE_HOST = "/libraT/servicehost"

	MACHINE_DOC     = "/libraT/release"
	MACHINE_COMPANY = "/libraT/company"
	LIBRA_LEVEL     = "/libraT/libralevel"

	NETWORK_IFACE        = "/network/interface"
	SERVER_ADDRESS       = "/management/serverip"
	ENCODER_LIVE_ADDRESS = "/encoder/live_addr"
	TF_ENABLE            = "/mountpoint/tf"
	TF_DEVICE            = "/mountpoint/device"
	TF_MOUNTPOINT        = "/mountpoint/mount"
	NVR_CONFIG           = "/nvr"
	NVR_URL              = "/nvr/url"
	NVR_PASSAGEWAY       = "/nvr/passageway"

	RELEASE_VERSION     = "Release Version"
	MACHINE_MEMTYPE     = "Memory Type"
	MACHINE_DATE        = "Release Date"
	HARDWARE_VERSION    = "Hardware Version"
	MACHINE_MODEL       = "Model"
	MACHINE_MANUFACTURE = "Manufacture"

	/*-------------- for wlan card --------------*/
	WLAN_CARD_MODE        = "/wlan/mode"       // default to WLAN_CARD_AP_MODE
	WLAN_CARD_AP_SSID     = "/wlan/ap_ssid"    // default to opera_sn
	WLAN_CARD_AP_PASSWD   = "/wlan/ap_passwd"  // default to 12345678
	WLAN_CARD_AP_CHANNEL  = "/wlan/ap_channel" // default to 8
	WLAN_CARD_WIFI_SSID   = "/wlan/wifi_ssid"
	WLAN_CARD_WIFI_PASSWD = "/wlan/wifi_passwd"
	/*-------------- end wlan card --------------*/

	// config about bumble,in bumble.yaml
	BUMBLE_PORT     = "bumble/port"
	BUMBEL_INTERVAL = "bumble/interval"
	SUPERVISOR_CONN = "bumble/supervisord_rpc"
	REMOTE_ETCD_TLS = "bumble/remote_etcd_tls"
	DATE_MODE       = "date/datemode"
	NTP_ADDRESS     = "date/ntpaddr"
	//////////////////////////////////////////above not update///////////////////////////////////////

	USER_ID    = "/config/global/user_id"
	SAFENET_ID = "/config/global/safenet_id"
	TF_SIZE    = "config/global/tfsize"
	TF_AVIAL   = "config/global/tfavial"

	// factory default mode, 0 for factory, 1 for working
	FACTORY_DEFAULT = "/config/global/factory_default"

	// nsq address
	LOCAL_HOST = "/config/global/localhost"
	// LOCAL_ETCD_ADDRESS = "/config/global/local_etcd_addr"
	LOCAL_NSQ_ADDR_TCP = "/config/global/local_nsq_addr_tcp"
	// LOCAL_NSQ_ADDR_HTTP = "/config/global/local_nsq_addr_http"

	NSQ_ADDR_CMD_TCP    = "/config/nsqaddr/nsq_addr_cmd_tcp"
	NSQ_ADDR_CMD_HTTP   = "/config/nsqaddr/nsq_addr_cmd_http"
	NSQ_ADDR_RESP_TCP   = "/config/nsqaddr/nsq_addr_resp_tcp"
	NSQ_ADDR_REPORT_TCP = "/config/nsqaddr/nsq_addr_report_tcp"

	// topic
	CMD_TOPIC    = "/config/topic/cmd_topic"
	RESP_TOPIC   = "/config/topic/resp_topic"
	REPORT_TOPIC = "/config/topic/report_topic"

	REMOTE_ETCD_ADDRESS = "/config/global/remote_etcd_addr"

	// REMOTE_NSQ_ADDR_TCP  = "/config/global/remote_nsq_addr_tcp"
	// REMOTE_NSQ_ADDR_HTTP = "/config/global/remote_nsq_addr_http"

	/*-------------- for prime --------------*/
	PRIME_NSQ_ADDR_CMD_TCP    = "/config/nsqaddr/prime_nsq_addr_cmd"
	PRIME_NSQ_ADDR_CMD_HTTP   = "/config/nsqaddr/prime_nsq_addr_cmd_http"
	PRIME_NSQ_ADDR_RESP_TCP   = "/config/nsqaddr/prime_nsq_addr_resp"
	PRIME_NSQ_ADDR_REPORT_TCP = "/config/nsqaddr/prime_nsq_addr_report"

	PRIME_CMD_TOPIC          = "/config/topic/prime_cmd_topic"
	PRIME_RESP_TOPIC         = "/config/topic/prime_resp_topic"
	PRIME_RESP_DEFAULT_TOPIC = "/config/topic/prime_resp_default_topic"
	PRIME_REPORT_TOPIC       = "/config/topic/prime_report_topic"

	// PRIME_ETCD_ADDR = "/config/global/prime_etcd_addr"

	PRIME_DB_ADDR   = "/config/db/prime_db_addr"
	PRIME_DB_NAME   = "/config/db/prime_db"
	PRIME_DB_C      = "/config/db/prime_c"
	PRIME_DB_USER   = "/config/db/prime_user"
	PRIME_DB_PASSWD = "/config/db/prime_passwd"

	/*-------------- end prime --------------*/

	/*-------------- for libra --------------*/
	LIBRA_PRELOAD_PATH      = "/config/libra/data/preload_path"
	LIBRA_ENABLE_ABNORMAL   = "/config/libra/data/enable_abnormal_action_detection"
	LIBRA_ABNORMAL_THRES    = "/config/libra/data/abnormal_threshold"
	LIBRA_ENABLE_FALL       = "/config/libra/data/enable_fall_detection"
	LIBRA_ENABLE_INVALID    = "/config/libra/data/enable_invalid_operation_detection"
	LIBRA_INVALID_THRES     = "/config/libra/data/invalid_operation_time_threshold"
	LIBRA_ENABLE_LATCH      = "/config/libra/data/enable_latch"
	LIBRA_ENABLE_LENS       = "/config/libra/data/enable_lens_protection"
	LIBRA_CUTBOARD_DISTANCE = "/config/libra/data/cutboard_update_distance"
	LIBRA_MINIMUM_HEIGHT    = "/config/libra/data/detection_minimum_height"
	LIBRA_ENABLE_DARKNESS   = "/config/libra/data/enable_darkness_detection"

	/*-------------- end libra --------------*/

	/*-------------- for eventbrain --------------*/
	EVENTBRAIN_HOTSPOT_PREFIX = "/config/eventbrain/hotspots"
	EVENTBRAIN_VELOCITY       = "/config/eventbrain/FeatureSchemas/0/FeatureRules/0/UpperBound" // Velocity
	// "/config/eventbrain/FeatureSchemas/0/FeatureRules/1/UpperBound" //
	// "/config/eventbrain/FeatureSchemas/0/FeatureRules/2/UpperBound"
	EVENTBRAIN_DISTANCE = "/config/eventbrain/FeatureSchemas/0/FeatureRules/3/UpperBound" // Distance
	EVENTBRAIN_DWELLING = "/config/eventbrain/FeatureSchemas/0/FeatureRules/4/UpperBound" // DwellingTime
	// "/config/eventbrain/FeatureSchemas/0/FeatureRules/5/UpperBound"
	// "/config/eventbrain/FeatureSchemas/0/FeatureRules/6/UpperBound"
	EVENTBRAIN_POPULATION = "/config/eventbrain/FeatureSchemas/128/FeatureRules/0/UpperBound" // Population
	// "/config/eventbrain/FeatureSchemas/128/FeatureRules/1/UpperBound"
	EVENTBRAIN_APPROACH  = "/config/eventbrain/FeatureSchemas/128/FeatureRules/2/UpperBound" // Approach
	EVENTBRAIN_INTRUSION = "/config/eventbrain/FeatureSchemas/128/FeatureRules/3/UpperBound" // Intrusion

	//
	EVENTBRAIN_FEATURE_SCHEMA_PREFIX = "/config/eventbrain"
	EVENTBRAIN_VELOCITY_SUFFIX       = "FeatureSchemas/0/FeatureRules/0/UpperBound"
	EVENTBRAIN_DISTANCE_SUFFIX       = "FeatureSchemas/0/FeatureRules/3/UpperBound"
	EVENTBRAIN_DWELLING_SUFFIX       = "FeatureSchemas/0/FeatureRules/4/UpperBound"
	EVENTBRAIN_POPULATION_SUFFIX     = "FeatureSchemas/128/FeatureRules/0/UpperBound"
	EVENTBRAIN_APPROACH_SUFFIX       = "FeatureSchemas/128/FeatureRules/2/UpperBound"
	EVENTBRAIN_INTRUSION_SUFFIX      = "FeatureSchemas/128/FeatureRules/3/UpperBound"
	//
	EVENTBRAIN_ALERT_RULE_PREFIX   = "/config/eventbrain/alertrule"
	EVENTBRAIN_DWELLINGTIME_LEGACY = "dwellingtime_legacy"
	EVENTBRAIN_POPULATION_LEGACY   = "population_legacy"
	EVENTBRAIN_VELOCITY_LEGACY     = "velocity_legacy"
	EVENTBRAIN_APPROACHING_LEGACY  = "approaching_legacy"
	EVENTBRAIN_DISTANCE_LEGACY     = "distance_legacy"

	/*-------------- end eventbrain --------------*/

	/*-------------- for OPTIMUS --------------*/
	SENSOR_GROUP_PREFIX          = "/config/sensor"
	SENSOR_GROUP_TBD_PREFIX      = "/config/sensortbd"
	USER_GROUP_PREFIX            = "/config/user"
	OPTIMUS_SENSOR_2_SYNC_PREFIX = "/config/sensor2sync"
	OPTIMUS_SENSOR_GROUP         = "/config/group"
	OPTIMUS_STORAGES             = "config/optimus/storage"

	OPTIMUS_TUNER_CMD_SUFFIX  = "topic/prime_tuner_cmd_topic"
	OPTIMUX_TUNER_RESP_SUFFIX = "topic/prime_tuner_resp_topic"

	OPTIMUX_BUMBLE_CMD_SUFFIX    = "topic/prime_bumble_cmd_topic"
	OPTIMUX_BUMBLE_RESP_SUFFIX   = "topic/prime_bumble_resp_topic"
	OPTIMUX_BUMBLE_REPORT_SUFFIX = "topic/prime_bumble_report_topic"
	/*-------------- end OPTIMUS --------------*/

	/*-------------- for TUNER --------------*/
	NSQ_ADDR_TUNER_CMD_TCP  = "/config/nsqaddr/nsq_addr_tuner_cmd"
	NSQ_ADDR_TUNER_CMD_HTTP = "/config/nsqaddr/nsq_addr_tuner_cmd_http"
	NSQ_ADDR_TUNER_RESP_TCP = "/config/nsqaddr/nsq_addr_tuner_resp"

	TUNER_CMD_TOPIC  = "/config/topic/prime_tuner_cmd_topic"
	TUNER_RESP_TOPIC = "/config/topic/prime_tuner_resp_topic"

	PRIME_NSQ_ADDR_TUNER_CMD_TCP  = "/config/nsqaddr/prime_nsq_addr_tuner_cmd"
	PRIME_NSQ_ADDR_TUNER_CMD_HTTP = "/config/nsqaddr/prime_nsq_addr_tuner_cmd_http"
	PRIME_NSQ_ADDR_TUNER_RESP_TCP = "/config/nsqaddr/prime_nsq_addr_tuner_resp"
	/*-------------- end TUNER --------------*/

	/*-------------- for EVENTSERVER --------------*/
	PRIME_NSQ_ADDR_ES_EVENT_TCP = "config/nsqaddr/prime_event_topic_server"
	PRIME_NSQ_ADDR_ES_TSS_TCP   = "config/nsqaddr/prime_tss_topic_server"
	PRIME_NSQ_ADDR_ES_VIBO_TCP  = "config/nsqaddr/prime_vibo_topic_server"
	PRIME_NSQ_ADDR_ES_VIDEO_TCP = "config/nsqaddr/prime_video_topic_server"
	EVENTSERVER_PUB_2DB         = "config/eventserver/pub_db_url"
	EVENTSERVER_PUB_2VIBO       = "config/eventserver/pub_vibo_url"
	/*-------------- end EVENTSERVER --------------*/

	/*-------------- for IMAGE --------------*/
	PRIME_IMAGE_PREFIX         = "/config/image"
	SENSOR_DOWNLOAD_PACKAGE    = "/config/global/download"
	SENSOR_UPGRADE_PACKAGE     = "/config/global/upgrade"
	SENSOR_PACKAGE_CUR_VERSION = "/config/global/cur_version"
	SENSOR_PACKAGE_DWN_VERSION = "/config/global/dwn_version"
	SENSOR_PACKAGE_LST_VERSION = "/config/global/lst_version"
	/*-------------- end IMAGE --------------*/

	/*-------------- for TEMPLATE --------------*/
	PRIME_TEMPLATE_PREFIX = "config/template"
	PRIME_CONFIG_PREFIX   = "config"
)

const (
	WLAN_CARD_WIFI_MODE = "wifi"
	WLAN_CARD_AP_MODE   = "ap"
	WLAN_CARD_IDLE_MODE = "idle"
)

var (
	PROCVISOR_IGNORE = []string{
		"etcd",
		"nsqd",
		"muxerlab",
	}
)

const (
	PROCVISOR_ACTION_START   = "start"
	PROCVISOR_ACTION_STOP    = "stop"
	PROCVISOR_ACTION_RESTART = "restart"
	PROCVISOR_ACTION_LOGTAIL = "tail"

	PROCVISOR_ACTION_NEXTDESIGN = ""
)

const (
	CONTAINER_CMD_SYNC_TIME         = "synctime.sh"
	CONTAINER_CMD_SYNC_TIME_KEYWORD = "offset"

	CONTAINER_CMD_DATE = "date.sh"
)

const (
	SERVICE_ADDR_ANY      = ""
	SERVICE_ADDR_LOCAL    = "localhost"
	SERVICE_ADDR_EXTERNIP = "netaddr"

	SENSOR_DEFAULT_NETINTERFACE = "eth0"
	SENSOR_DEFAULT_DATEMODE     = 2
)

type RestFulResp struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

const (
	RESTFUL_RESPCODE_OK     = 0
	RESTFUL_RESPCODE_FAIL   = 1
	RESTFUL_RESPCODE_UNAUTH = 2
)

const (
	PRIME_SENSE_USB_ID = "1d27:0601"
)

// for onvif config
type OnvifConfig struct {
	server_ip   string           `xml:"server_ip"`
	server_port string           `xml:"server_port"`
	need_auth   string           `xml:"need_auth"`
	log_enable  string           `xml:"log_enable"`
	information OnvifInformation `xml:"information"`
	user        OnvifUser        `xml:"user"`
	profile     []OnvifProfile   `xml:"profile"`
	scope       []string         `xml:"scope"`
}

type OnvifInformation struct {
	Manufacturer    string `xml:"Manufacturer"`
	Model           string `xml:"Model"`
	FirmwareVersion string `xml:"FirmwareVersion"`
	SerialNumber    string `xml:"SerialNumber"`
	HardwareId      string `xml:"HardwareId"`
}

type OnvifUser struct {
	Username  string `xml:"username"`
	Password  string `xml:"password"`
	Userlevel string `xml:"userlevel"`
}

type OnvifProfile struct {
	VideoSource   OnvifVideo_source  `xml:"video_source"`
	Video_encoder OnvifVideo_encoder `xml:"video_encoder"`
	audio_source  string             `xml:"audio_source"`
	audio_encoder string             `xml:"audio_encoder"`
	stream_uri    string             `xml:"stream_uri"`
}

type OnvifVideo_source struct {
	Width  int `xml:"width"`
	height int `xml:"height"`
}

type OnvifVideo_encoder struct {
	width             int    `xml:"width"`
	height            int    `xml:"height"`
	quality           int    `xml:"quality"`
	session_timeout   int    `xml:"session_timeout"`
	framerate         int    `xml:"framerate"`
	encoding_interval int    `xml:"encoding_interval"`
	bitrate_limit     int    `xml:"bitrate_limit"`
	encoding          string `xml:"encoding"`
	h264              H264   `xml:"h264"`
}

type H264 struct {
	gov_length   int    `xml:"gov_length"`
	h264_profile string `xml:"h264_profile"`
}

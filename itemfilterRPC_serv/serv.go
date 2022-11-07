package itemfilterRPC_serv

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/featmate/itemfilterRPC/filterstore"
	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
	"github.com/featmate/itemfilterRPC/metadatastore"

	"github.com/Golang-Tools/idgener"
	log "github.com/Golang-Tools/loggerhelper/v2"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/xds"

	"net/http"

	_ "embed"

	"github.com/felixge/httpsnoop"
	"github.com/flowchartsman/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
)

// 这边需要指定文件位置
//
//go:embed service.swagger.json
var spec []byte

// Server grpc的服务器结构体
// 服务集成了如下特性:
// 设置收发最大消息长度
// 健康检测
// gzip做消息压缩
// 接口反射
// channelz支持
// TLS支持
// keep alive 支持
type Server struct {
	App_Name    string `json:"app_name,omitempty" jsonschema:"required,description=服务名,default=itemfilter"`
	App_Version string `json:"app_version,omitempty" jsonschema:"required,description=服务版本,default=0.0.1"`
	Address     string `json:"address,omitempty" jsonschema:"required,description=服务的主机和端口,default=0.0.0.0:5000"`
	Log_Level   string `json:"log_level,omitempty" jsonschema:"required,description=项目的log等级,enum=TRACE,enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR,default=DEBUG"`

	// 性能设置
	Max_Recv_Msg_Size                           int  `json:"max_recv_msg_size,omitempty" jsonschema:"description=允许接收的最大消息长度"`
	Max_Send_Msg_Size                           int  `json:"max_send_msg_size,omitempty" jsonschema:"description=允许发送的最大消息长度"`
	Initial_Window_Size                         int  `json:"initial_window_size,omitempty" jsonschema:"description=基于Stream的滑动窗口大小"`
	Initial_Conn_Window_Size                    int  `json:"initial_conn_window_size,omitempty" jsonschema:"description=基于Connection的滑动窗口大小"`
	Max_Concurrent_Streams                      int  `json:"max_concurrent_streams,omitempty" jsonschema:"description=一个连接中最大并发Stream数"`
	Max_Connection_Idle                         int  `json:"max_connection_idle,omitempty" jsonschema:"description=客户端连接的最大空闲时长"`
	Max_Connection_Age                          int  `json:"max_connection_age,omitempty" jsonschema:"description=如果连接存活超过n则发送goaway"`
	Max_Connection_Age_Grace                    int  `json:"max_connection_age_grace,omitempty" jsonschema:"description=强制关闭连接之前允许等待的rpc在n秒内完成"`
	Keepalive_Time                              int  `json:"keepalive_time,omitempty" jsonschema:"description=空闲连接每隔n秒ping一次客户端已确保连接存活"`
	Keepalive_Timeout                           int  `json:"keepalive_timeout,omitempty" jsonschema:"description=ping时长超过n则认为连接已死"`
	Keepalive_Enforcement_Min_Time              int  `json:"keepalive_enforement_min_time,omitempty" jsonschema:"description=如果客户端超过每n秒ping一次则终止连接"`
	Keepalive_Enforcement_Permit_Without_Stream bool `json:"keepalive_enforement_permit_without_stream,omitempty" jsonschema:"description=即使没有活动流也允许ping"`

	//TLS设置
	Server_Cert_Path string `json:"server_cert_path,omitempty" jsonschema:"description=使用TLS时服务端的证书位置"`
	Server_Key_Path  string `json:"server_key_path,omitempty" jsonschema:"description=使用TLS时服务端证书的私钥位置"`
	Ca_Cert_Path     string `json:"ca_cert_path,omitempty" jsonschema:"description=使用TLS时根证书位置"`
	Client_Crl_Path  string `json:"client_crl_path,omitempty" jsonschema:"description=客户端证书黑名单位置"`

	//使用XDS
	XDS                  bool `json:"xds,omitempty" jsonschema:"description=是否使用xDSAPIs"`
	XDS_CREDS            bool `json:"xds_creds,omitempty" jsonschema:"description=是否使用xDSAPIs来接收TLS设置"`
	XDS_Maintenance_Port int  `json:"xds_maintenance_port,omitempty" jsonschema:"description=maintenance服务的端口如果不设置则使用当前服务端口+1"`

	// 调试,目前admin接口不稳定,因此只使用channelz
	Use_Admin bool `json:"use_admin,omitempty" jsonschema:"description=是否使用grpc-admin方便调试"`

	UnaryInterceptors  []grpc.UnaryServerInterceptor  `json:"-" jsonschema:"nullable"`
	StreamInterceptors []grpc.StreamServerInterceptor `json:"-" jsonschema:"nullable"`

	//使用grpc-gateway
	Use_Gateway              bool   `json:"use_gateway,omitempty" jsonschema:"description=是否使用grpc-gateway提供restful接口"`
	Gateway_Use_TLS          bool   `json:"gateway_use_tls,omitempty" jsonschema:"description=grpc-gateway是否使用TLS以启动https服务"`
	Gateway_Client_Cert_Path string `json:"gateway_client_cert_path,omitempty" jsonschema:"description=指定gateway访问grpc的客户端证书当XDS_CREDS或设置了TLS设置时生效"`
	Gateway_Client_Key_Path  string `json:"gateway_client_key_path,omitempty" jsonschema:"description=指定gateway访问grpc的客户端证书对应的私钥位置"`

	//业务配置
	// FilterRedis配置
	FilterRedis_URL         string `json:"filterredis_url" jsonschema:"required,description=过滤器redis路径集群以逗号隔开,default=redis://localhost:6379"`
	FilterRedis_RouteMod    string `json:"filterredis_routemod" jsonschema:"description=redis集群的连接策略,enum=bylatency,enum=randomly,enum=none,default=none"`
	FilterRedis_ConnTimeout int    `json:"filterredis_conntimeout" jsonschema:"required,description=过滤器redis请求的最大超时时间单位ms,default=50"`

	// IDgenger
	IDGener string `json:"blacklistconfig_idgener,omitempty" jsonschema:"description=黑名单创建时的默认id生成器,enum=uuid4,enum=sonyflake,enum=snowflake,default=uuid4"`

	// 元数据存储
	MetaDataStore_Type        string `json:"metadatastore_type,omitempty" jsonschema:"description=信息元数据保存类型,enum=self,enum=redis,enum=etcd,default=self"`
	MetaDataStore_URL         string `json:"metadatastore_url,omitempty" jsonschema:"description=信息元数据保存服务位置"`
	MetaDataStore_RouteMod    string `json:"metadatastore_routemod,omitempty" jsonschema:"description=信息元数据保存在redis集群时的请求路由模式,enum=bylatency,enum=randomly,enum=none"`
	MetaDataStore_ConnTimeout int    `json:"metadatastore_conntimeout,omitempty" jsonschema:"description=请求信息元数据保存服务的最大超时单位ms,default=50"`

	// 黑名单
	BlacklistConfig_RangeInfoKey        string `json:"blacklistconfig_rangeinfokey,omitempty" jsonschema:"description=黑名单范围保存的key不包括命名空间,default=blacklists"`
	BlacklistConfig_InfoKeySubnamespace string `json:"blacklistconfig_infokeysubnamespace,omitempty" jsonschema:"description=黑名单描述信息保存的子命名空间不包括总命名空间,default=infoblacklist"`
	BlacklistConfig_KeySubnamespace     string `json:"blacklistconfig_keysubnamespace,omitempty" jsonschema:"description=黑名单保存的子命名空间不包括总命名空间,default=BID"`
	BlacklistConfig_DefaultTTLDays      int    `json:"blacklistconfig_defaultttldays,omitempty" jsonschema:"description=黑名单包括描述信息的默认保存时间"`

	// 范围
	RangeConfig_RangeInfoKey         string  `json:"rangeconfig_rangeinfokey,omitempty" jsonschema:"description=范围范围保存的key不包括命名空间,default=range"`
	RangeConfig_InfoKeySubnamespace  string  `json:"rangeconfig_infokeysubnamespace,omitempty" jsonschema:"description=范围描述信息保存的子命名空间不包括总命名空间,default=inforange"`
	RangeConfig_KeySubnamespace      string  `json:"rangeconfig_keysubnamespace,omitempty" jsonschema:"description=范围保存的子命名空间不包括总命名空间,default=RID"`
	RangeConfig_DefaultTTLDays       int     `json:"rangeconfig_defaultttldays,omitempty" jsonschema:"description=范围包括描述信息的默认保存时间"`
	RangeConfig_FilterType           string  `json:"rangeconfig_filtertype,omitempty" jsonschema:"description=范围使用的过滤器类型,enum=set,enum=bloom,enum=cuckoo,default=set"`
	RangeConfig_DefaultCapacity      int     `json:"rangeconfig_defaultcapacity,omitempty" jsonschema:"description=范围的默认预估容量,default=5000"`
	RangeConfig_DefaultErrorRate     float64 `json:"rangeconfig_defaulterrorrate,omitempty" jsonschema:"description=范围的默认碰撞率"`
	RangeConfig_DefaultExpansion     int     `json:"rangeconfig_defaultexpansion,omitempty" jsonschema:"description=范围的默认扩容倍数,default=1"`
	RangeConfig_DefaultNonScaling    bool    `json:"crangeconfig_defaultnonscaleing,omitempty" jsonschema:"description=范围默认是否会扩容过滤器"`
	RangeConfig_DefaultBucketSize    int     `json:"rangeconfig_defaultbucketsize,omitempty" jsonschema:"description=范围默认CUCKOO桶大小,default=1000"`
	RangeConfig_DefaultMaxIterations int     `json:"rangeconfig_defaultmaxiterations,omitempty" jsonschema:"description=范围默认交换项目的次数,default=4"`

	//namespace
	NamespaceConfig_EntitySourceTypeSubnamespace  string `json:"namespaceconfig_entitysourcetypeysubnamespace,omitempty" jsonschema:"description=实体类型的子命名空间名,default=EST"`
	NamespaceConfig_PickerSubnamespace            string `json:"namespaceconfig_pickersubnamespace,omitempty" jsonschema:"description=Picker的子命名空间名,default=PID"`
	NamespaceConfig_PickerIntervalSubnamespace    string `json:"namespaceconfig_pickerintervalsubnamespace,omitempty" jsonschema:"description=PickerInterval的子命名空间名,default=PI"`
	NamespaceConfig_PickerCounterSubnamespace     string `json:"namespaceconfig_pickercountersubnamespace,omitempty" jsonschema:"description=PickerCounter键的子命名空间,default=PC"`
	NamespaceConfig_PickerCounterDateSubnamespace string `json:"namespaceconfig_pickercounterdatesubnamespace,omitempty" jsonschema:"description=PickerCounter键的日期,default=PCD"`
	// 会话

	ContextConfig_KeySubnamespace      string  `json:"contextconfig_keysubnamespace,omitempty" jsonschema:"description=会话子命名空间不包括总命名空间,default=CID"`
	ContextConfig_FilterType           string  `json:"contextconfig_filtertype,omitempty" jsonschema:"description=过滤器类型,enum=set,enum=bloom,enum=cuckoo,default=set"`
	ContextConfig_DefaultCapacity      int     `json:"contextconfig_defaultcapacity,omitempty" jsonschema:"description=范围的预估容量,default=5000"`
	ContextConfig_DefaultErrorRate     float64 `json:"contextconfig_defaulterrorrate,omitempty" jsonschema:"description=碰撞率"`
	ContextConfig_DefaultExpansion     int     `json:"contextconfig_defaultexpansion,omitempty" jsonschema:"description=扩容倍数,default=1"`
	ContextConfig_DefaultNonScaling    bool    `json:"contextconfig_defaultnonscaleing,omitempty" jsonschema:"description=不会扩容过滤器"`
	ContextConfig_DefaultBucketSize    int     `json:"contextconfig_defaultbucketsize,omitempty" jsonschema:"description=CUCKOO桶大小,default=1000"`
	ContextConfig_DefaultMaxIterations int     `json:"contextconfig_defaultmaxiterations,omitempty" jsonschema:"description=交换项目的次数,default=4"`
	ContextConfig_DefaultTTLSeconds    int     `json:"contextconfig_defaultttlseconds,omitempty" jsonschema:"description=会话的的默认过期时间,default=900"`

	// Picker
	PickerConfig_TimeRange []string `json:"pickerconfig_timerange,omitempty" jsonschema:"description=维护的时间间隔,enum=day,enum=week,enum=month,enum=year,enum=infi,default=day,default=infi"`
	/// day
	PickerConfig_Day_FilterType           string  `json:"pickerconfig_day_filtertype,omitempty" jsonschema:"description=过滤器类型,enum=set,enum=bloom,enum=cuckoo,default=set"`
	PickerConfig_Day_DefaultCapacity      int     `json:"pickerconfig_day_defaultcapacity,omitempty" jsonschema:"description=范围的预估容量,default=5000"`
	PickerConfig_Day_DefaultErrorRate     float64 `json:"pickerconfig_day_defaulterrorrate,omitempty" jsonschema:"description=碰撞率"`
	PickerConfig_Day_DefaultExpansion     int     `json:"pickerconfig_day_defaultexpansion,omitempty" jsonschema:"description=扩容倍数,default=1"`
	PickerConfig_Day_DefaultNonScaling    bool    `json:"pickerconfig_day_defaultnonscaleing,omitempty" jsonschema:"description=不会扩容过滤器"`
	PickerConfig_Day_DefaultBucketSize    int     `json:"pickerconfig_day_defaultbucketsize,omitempty" jsonschema:"description=CUCKOO桶大小,default=1000"`
	PickerConfig_Day_DefaultMaxIterations int     `json:"pickerconfig_day_defaultmaxiterations,omitempty" jsonschema:"description=交换项目的次数,default=4"`

	/// week
	PickerConfig_Week_FilterType           string  `json:"pickerconfig_week_filtertype,omitempty" jsonschema:"description=过滤器类型,enum=set,enum=bloom,enum=cuckoo,default=set"`
	PickerConfig_Week_DefaultCapacity      int     `json:"pickerconfig_week_defaultcapacity,omitempty" jsonschema:"description=范围的预估容量,default=5000"`
	PickerConfig_Week_DefaultErrorRate     float64 `json:"pickerconfig_week_defaulterrorrate,omitempty" jsonschema:"description=碰撞率"`
	PickerConfig_Week_DefaultExpansion     int     `json:"pickerconfig_week_defaultexpansion,omitempty" jsonschema:"description=扩容倍数,default=1"`
	PickerConfig_Week_DefaultNonScaling    bool    `json:"pickerconfig_week_defaultnonscaleing,omitempty" jsonschema:"description=不会扩容过滤器"`
	PickerConfig_Week_DefaultBucketSize    int     `json:"pickerconfig_week_defaultbucketsize,omitempty" jsonschema:"description=CUCKOO桶大小,default=1000"`
	PickerConfig_Week_DefaultMaxIterations int     `json:"pickerconfig_week_defaultmaxiterations,omitempty" jsonschema:"description=交换项目的次数,default=4"`

	/// month
	PickerConfig_Month_FilterType           string  `json:"pickerconfig_month_filtertype,omitempty" jsonschema:"description=过滤器类型,enum=set,enum=bloom,enum=cuckoo,default=set"`
	PickerConfig_Month_DefaultCapacity      int     `json:"pickerconfig_month_defaultcapacity,omitempty" jsonschema:"description=范围的预估容量,default=5000"`
	PickerConfig_Month_DefaultErrorRate     float64 `json:"pickerconfig_month_defaulterrorrate,omitempty" jsonschema:"description=碰撞率"`
	PickerConfig_Month_DefaultExpansion     int     `json:"pickerconfig_month_defaultexpansion,omitempty" jsonschema:"description=扩容倍数,default=1"`
	PickerConfig_Month_DefaultNonScaling    bool    `json:"pickerconfig_month_defaultnonscaleing,omitempty" jsonschema:"description=不会扩容过滤器"`
	PickerConfig_Month_DefaultBucketSize    int     `json:"pickerconfig_month_defaultbucketsize,omitempty" jsonschema:"description=CUCKOO桶大小,default=1000"`
	PickerConfig_Month_DefaultMaxIterations int     `json:"pickerconfig_month_defaultmaxiterations,omitempty" jsonschema:"description=交换项目的次数,default=4"`

	/// year
	PickerConfig_Year_FilterType           string  `json:"pickerconfig_year_filtertype,omitempty" jsonschema:"description=过滤器类型,enum=set,enum=bloom,enum=cuckoo,default=set"`
	PickerConfig_Year_DefaultCapacity      int     `json:"pickerconfig_year_defaultcapacity,omitempty" jsonschema:"description=范围的预估容量,default=5000"`
	PickerConfig_Year_DefaultErrorRate     float64 `json:"pickerconfig_year_defaulterrorrate,omitempty" jsonschema:"description=碰撞率"`
	PickerConfig_Year_DefaultExpansion     int     `json:"pickerconfig_year_defaultexpansion,omitempty" jsonschema:"description=扩容倍数,default=1"`
	PickerConfig_Year_DefaultNonScaling    bool    `json:"pickerconfig_year_defaultnonscaleing,omitempty" jsonschema:"description=不会扩容过滤器"`
	PickerConfig_Year_DefaultBucketSize    int     `json:"pickerconfig_year_defaultbucketsize,omitempty" jsonschema:"description=CUCKOO桶大小,default=1000"`
	PickerConfig_Year_DefaultMaxIterations int     `json:"pickerconfig_year_defaultmaxiterations,omitempty" jsonschema:"description=交换项目的次数,default=4"`

	/// infi
	PickerConfig_Infi_FilterType           string  `json:"pickerconfig_infi_filtertype,omitempty" jsonschema:"description=过滤器类型,enum=set,enum=bloom,enum=cuckoo,default=set"`
	PickerConfig_Infi_DefaultCapacity      int     `json:"pickerconfig_infi_defaultcapacity,omitempty" jsonschema:"description=范围的预估容量,default=5000"`
	PickerConfig_Infi_DefaultErrorRate     float64 `json:"pickerconfig_infi_defaulterrorrate,omitempty" jsonschema:"description=碰撞率"`
	PickerConfig_Infi_DefaultExpansion     int     `json:"pickerconfig_infi_defaultexpansion,omitempty" jsonschema:"description=扩容倍数,default=1"`
	PickerConfig_Infi_DefaultNonScaling    bool    `json:"pickerconfig_infi_defaultnonscaleing,omitempty" jsonschema:"description=不会扩容过滤器"`
	PickerConfig_Infi_DefaultBucketSize    int     `json:"pickerconfig_infi_defaultbucketsize,omitempty" jsonschema:"description=CUCKOO桶大小,default=1000"`
	PickerConfig_Infi_DefaultMaxIterations int     `json:"pickerconfig_infi_defaultmaxiterations,omitempty" jsonschema:"description=交换项目的次数,default=4"`

	// PickerCounter功能设置
	PickerCounterOn                    bool `json:"pickercounter_on,omitempty" jsonschema:"description=是否启用PickerCounter功能"`
	PickerCounterConfig_TTLDays        int  `json:"pickercounterconfig_ttldays,omitempty" jsonschema:"description=PickerCounter保存的天数,default=7"`
	PickerCounterConfig_RangeCounterOn bool `json:"pickercounter_rangecounteron,omitempty" jsonschema:"description=是否启用范围PickerCounter"`
	PickerCounterConfig_RangeTTLDays   int  `json:"pickercounterconfig_rangettldays,omitempty" jsonschema:"description=范围统计PickerCounter保留的天数,default=7"`

	itemfilterRPC_pb.UnimplementedITEMFILTERRPCServer `json:"-"`
	opts                                              []grpc.ServerOption
	healthservice                                     *health.Server

	defaultContextFilterInfo *itemfilterRPC_pb.RedisFilterSetting
	defaultPickerFilterInfos map[string]*itemfilterRPC_pb.RedisFilterSetting
	metadatastore            metadatastore.MetadataStore
	filterstore              filterstore.FilterStore

	idGener idgener.IDGenInterface
}

// Main 服务的入口函数
func (s *Server) Main() {
	// 初始化log
	log.Set(log.WithLevel(s.Log_Level),
		log.AddExtField("app_name", s.App_Name),
		log.AddExtField("app_version", s.App_Version),
	)
	log.Info("grpc服务获得参数", log.Dict{"ServiceConfig": s})
	// 处理float64无法设置默认值的问题
	if s.ContextConfig_DefaultErrorRate <= 0 {
		s.ContextConfig_DefaultErrorRate = 0.001
	}
	if s.PickerConfig_Day_DefaultErrorRate <= 0 {
		s.PickerConfig_Day_DefaultErrorRate = 0.001
	}
	if s.PickerConfig_Week_DefaultErrorRate <= 0 {
		s.PickerConfig_Week_DefaultErrorRate = 0.001
	}
	if s.PickerConfig_Month_DefaultErrorRate <= 0 {
		s.PickerConfig_Month_DefaultErrorRate = 0.001
	}
	if s.PickerConfig_Year_DefaultErrorRate <= 0 {
		s.PickerConfig_Year_DefaultErrorRate = 0.001
	}
	if s.PickerConfig_Infi_DefaultErrorRate <= 0 {
		s.PickerConfig_Infi_DefaultErrorRate = 0.001
	}
	//默认blacklist过滤器
	defaultBlacklistFilterInfo := &itemfilterRPC_pb.SetFilterSetting{
		MaxTTLSeconds: int64(s.BlacklistConfig_DefaultTTLDays) * 24 * 3600,
	}
	// 默认range过滤器
	defaultRangeFilterInfo := &itemfilterRPC_pb.RedisFilterSetting{}
	switch s.RangeConfig_FilterType {
	case "cuckoo":
		{
			defaultRangeFilterInfo.Setting = &itemfilterRPC_pb.RedisFilterSetting_CuckoofilterSetting{
				CuckoofilterSetting: &itemfilterRPC_pb.CuckooFilterSetting{
					Capacity:      int64(s.RangeConfig_DefaultCapacity),
					Expansion:     int64(s.RangeConfig_DefaultExpansion),
					NonScaling:    s.RangeConfig_DefaultNonScaling,
					BucketSize:    int64(s.RangeConfig_DefaultBucketSize),
					MaxIterations: int64(s.RangeConfig_DefaultMaxIterations),
					MaxTTLSeconds: int64(s.RangeConfig_DefaultTTLDays) * 24 * 3600,
				},
			}
		}
	case "bloom":
		{
			defaultRangeFilterInfo.Setting = &itemfilterRPC_pb.RedisFilterSetting_BloomfilterSetting{
				BloomfilterSetting: &itemfilterRPC_pb.BloomFilterSetting{
					Capacity:      int64(s.RangeConfig_DefaultCapacity),
					ErrorRate:     s.RangeConfig_DefaultErrorRate,
					Expansion:     int64(s.RangeConfig_DefaultExpansion),
					NonScaling:    s.RangeConfig_DefaultNonScaling,
					MaxTTLSeconds: int64(s.RangeConfig_DefaultTTLDays) * 24 * 3600,
				},
			}
		}
	default:
		{
			defaultRangeFilterInfo.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
				SetfilterSetting: &itemfilterRPC_pb.SetFilterSetting{
					MaxTTLSeconds: int64(s.RangeConfig_DefaultTTLDays) * 24 * 3600,
				},
			}
		}
	}
	//初始化defaultContextFilterInfo
	s.defaultContextFilterInfo = &itemfilterRPC_pb.RedisFilterSetting{}
	switch s.ContextConfig_FilterType {
	case "cuckoo":
		{
			s.defaultContextFilterInfo.Setting = &itemfilterRPC_pb.RedisFilterSetting_CuckoofilterSetting{
				CuckoofilterSetting: &itemfilterRPC_pb.CuckooFilterSetting{
					Capacity:      int64(s.ContextConfig_DefaultCapacity),
					Expansion:     int64(s.ContextConfig_DefaultExpansion),
					NonScaling:    s.ContextConfig_DefaultNonScaling,
					BucketSize:    int64(s.ContextConfig_DefaultBucketSize),
					MaxIterations: int64(s.ContextConfig_DefaultMaxIterations),
					MaxTTLSeconds: int64(s.ContextConfig_DefaultTTLSeconds),
				},
			}
		}
	case "bloom":
		{
			s.defaultContextFilterInfo.Setting = &itemfilterRPC_pb.RedisFilterSetting_BloomfilterSetting{
				BloomfilterSetting: &itemfilterRPC_pb.BloomFilterSetting{
					Capacity:      int64(s.ContextConfig_DefaultCapacity),
					ErrorRate:     s.ContextConfig_DefaultErrorRate,
					Expansion:     int64(s.ContextConfig_DefaultExpansion),
					NonScaling:    s.ContextConfig_DefaultNonScaling,
					MaxTTLSeconds: int64(s.ContextConfig_DefaultTTLSeconds),
				},
			}
		}
	default:
		{
			s.defaultContextFilterInfo.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
				SetfilterSetting: &itemfilterRPC_pb.SetFilterSetting{
					MaxTTLSeconds: int64(s.ContextConfig_DefaultTTLSeconds),
				},
			}
		}
	}

	//初始化defaultPickerFilterInfos
	s.defaultPickerFilterInfos = map[string]*itemfilterRPC_pb.RedisFilterSetting{}
	for _, tr := range s.PickerConfig_TimeRange {
		switch tr {
		case "day":
			{
				filterdefaultsetting := itemfilterRPC_pb.RedisFilterSetting{}
				switch s.PickerConfig_Day_FilterType {
				case "cuckoo":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_CuckoofilterSetting{
							CuckoofilterSetting: &itemfilterRPC_pb.CuckooFilterSetting{
								Capacity:      int64(s.PickerConfig_Day_DefaultCapacity),
								Expansion:     int64(s.PickerConfig_Day_DefaultExpansion),
								NonScaling:    s.PickerConfig_Day_DefaultNonScaling,
								BucketSize:    int64(s.PickerConfig_Day_DefaultBucketSize),
								MaxIterations: int64(s.PickerConfig_Day_DefaultMaxIterations),
							},
						}
					}
				case "bloom":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_BloomfilterSetting{
							BloomfilterSetting: &itemfilterRPC_pb.BloomFilterSetting{
								Capacity:   int64(s.PickerConfig_Day_DefaultCapacity),
								ErrorRate:  s.PickerConfig_Day_DefaultErrorRate,
								Expansion:  int64(s.PickerConfig_Day_DefaultExpansion),
								NonScaling: s.PickerConfig_Day_DefaultNonScaling,
							},
						}
					}
				default:
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
							SetfilterSetting: &itemfilterRPC_pb.SetFilterSetting{},
						}
					}
				}
				s.defaultPickerFilterInfos["day"] = &filterdefaultsetting
			}
		case "week":
			{
				filterdefaultsetting := itemfilterRPC_pb.RedisFilterSetting{}
				switch s.PickerConfig_Week_FilterType {
				case "cuckoo":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_CuckoofilterSetting{
							CuckoofilterSetting: &itemfilterRPC_pb.CuckooFilterSetting{
								Capacity:      int64(s.PickerConfig_Week_DefaultCapacity),
								Expansion:     int64(s.PickerConfig_Week_DefaultExpansion),
								NonScaling:    s.PickerConfig_Week_DefaultNonScaling,
								BucketSize:    int64(s.PickerConfig_Week_DefaultBucketSize),
								MaxIterations: int64(s.PickerConfig_Week_DefaultMaxIterations),
							},
						}
					}
				case "bloom":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_BloomfilterSetting{
							BloomfilterSetting: &itemfilterRPC_pb.BloomFilterSetting{
								Capacity:   int64(s.PickerConfig_Week_DefaultCapacity),
								ErrorRate:  s.PickerConfig_Week_DefaultErrorRate,
								Expansion:  int64(s.PickerConfig_Week_DefaultExpansion),
								NonScaling: s.PickerConfig_Week_DefaultNonScaling,
							},
						}
					}
				default:
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
							SetfilterSetting: &itemfilterRPC_pb.SetFilterSetting{},
						}
					}
				}
				s.defaultPickerFilterInfos["week"] = &filterdefaultsetting
			}
		case "month":
			{
				filterdefaultsetting := itemfilterRPC_pb.RedisFilterSetting{}
				switch s.PickerConfig_Month_FilterType {
				case "cuckoo":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_CuckoofilterSetting{
							CuckoofilterSetting: &itemfilterRPC_pb.CuckooFilterSetting{
								Capacity:      int64(s.PickerConfig_Month_DefaultCapacity),
								Expansion:     int64(s.PickerConfig_Month_DefaultExpansion),
								NonScaling:    s.PickerConfig_Month_DefaultNonScaling,
								BucketSize:    int64(s.PickerConfig_Month_DefaultBucketSize),
								MaxIterations: int64(s.PickerConfig_Month_DefaultMaxIterations),
							},
						}
					}
				case "bloom":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_BloomfilterSetting{
							BloomfilterSetting: &itemfilterRPC_pb.BloomFilterSetting{
								Capacity:   int64(s.PickerConfig_Month_DefaultCapacity),
								ErrorRate:  s.PickerConfig_Month_DefaultErrorRate,
								Expansion:  int64(s.PickerConfig_Month_DefaultExpansion),
								NonScaling: s.PickerConfig_Month_DefaultNonScaling,
							},
						}
					}
				default:
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
							SetfilterSetting: &itemfilterRPC_pb.SetFilterSetting{},
						}
					}
				}
				s.defaultPickerFilterInfos["month"] = &filterdefaultsetting
			}
		case "year":
			{
				filterdefaultsetting := itemfilterRPC_pb.RedisFilterSetting{}
				switch s.PickerConfig_Year_FilterType {
				case "cuckoo":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_CuckoofilterSetting{
							CuckoofilterSetting: &itemfilterRPC_pb.CuckooFilterSetting{
								Capacity:      int64(s.PickerConfig_Year_DefaultCapacity),
								Expansion:     int64(s.PickerConfig_Year_DefaultExpansion),
								NonScaling:    s.PickerConfig_Year_DefaultNonScaling,
								BucketSize:    int64(s.PickerConfig_Year_DefaultBucketSize),
								MaxIterations: int64(s.PickerConfig_Year_DefaultMaxIterations),
							},
						}
					}
				case "bloom":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_BloomfilterSetting{
							BloomfilterSetting: &itemfilterRPC_pb.BloomFilterSetting{
								Capacity:   int64(s.PickerConfig_Year_DefaultCapacity),
								ErrorRate:  s.PickerConfig_Year_DefaultErrorRate,
								Expansion:  int64(s.PickerConfig_Year_DefaultExpansion),
								NonScaling: s.PickerConfig_Year_DefaultNonScaling,
							},
						}
					}
				default:
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
							SetfilterSetting: &itemfilterRPC_pb.SetFilterSetting{},
						}
					}
				}
				s.defaultPickerFilterInfos["year"] = &filterdefaultsetting
			}
		case "infi":
			{
				filterdefaultsetting := itemfilterRPC_pb.RedisFilterSetting{}
				switch s.PickerConfig_Infi_FilterType {
				case "cuckoo":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_CuckoofilterSetting{
							CuckoofilterSetting: &itemfilterRPC_pb.CuckooFilterSetting{
								Capacity:      int64(s.PickerConfig_Infi_DefaultCapacity),
								Expansion:     int64(s.PickerConfig_Infi_DefaultExpansion),
								NonScaling:    s.PickerConfig_Infi_DefaultNonScaling,
								BucketSize:    int64(s.PickerConfig_Infi_DefaultBucketSize),
								MaxIterations: int64(s.PickerConfig_Infi_DefaultMaxIterations),
							},
						}
					}
				case "bloom":
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_BloomfilterSetting{
							BloomfilterSetting: &itemfilterRPC_pb.BloomFilterSetting{
								Capacity:   int64(s.PickerConfig_Infi_DefaultCapacity),
								ErrorRate:  s.PickerConfig_Infi_DefaultErrorRate,
								Expansion:  int64(s.PickerConfig_Infi_DefaultExpansion),
								NonScaling: s.PickerConfig_Infi_DefaultNonScaling,
							},
						}
					}
				default:
					{
						filterdefaultsetting.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
							SetfilterSetting: &itemfilterRPC_pb.SetFilterSetting{},
						}
					}
				}
				s.defaultPickerFilterInfos["infi"] = &filterdefaultsetting
			}
		}
	}
	// 初始化ID生成器
	bgentype := idgener.IDGEN_SNOWFLAKE
	switch s.IDGener {
	case "uuid4":
		{
			bgentype = idgener.IDGEN_UUIDV4
		}
	case "sonyflake":
		{
			bgentype = idgener.IDGEN_SNOYFLAKE
		}
	case "snowflake":
		{
			bgentype = idgener.IDGEN_SNOWFLAKE
		}
	}
	bg, err := idgener.IDGenNameToIDGen(bgentype)
	if err != nil {
		log.Error("init blacklist idgen get error", log.Dict{"err": err.Error()})
		os.Exit(1)
	}
	s.idGener = bg

	//初始化元数据后端
	metadatastoreOpts := metadatastore.Options{
		URL:                                 s.MetaDataStore_URL,
		RedisRouteMod:                       s.MetaDataStore_RouteMod,
		ConnTimeout:                         s.MetaDataStore_ConnTimeout,
		NamespaceBase:                       s.App_Name,
		BlacklistConfig_RangeInfoKey:        s.BlacklistConfig_RangeInfoKey,
		BlacklistConfig_InfoKeySubnamespace: s.BlacklistConfig_InfoKeySubnamespace,
		RangeConfig_RangeInfoKey:            s.RangeConfig_RangeInfoKey,
		RangeConfig_InfoKeySubnamespace:     s.RangeConfig_InfoKeySubnamespace,
		RangeConfig_FilterType:              s.RangeConfig_FilterType,
		BlacklistConfig_DefaultTTLDays:      s.BlacklistConfig_DefaultTTLDays,
		RangeConfig_DefaultTTLDays:          s.RangeConfig_DefaultTTLDays,
	}
	switch s.MetaDataStore_Type {
	case "etcd":
		{
			s.metadatastore = metadatastore.NewEtcdMetadataStore()
			err := s.metadatastore.Init(metadatastore.WithOptions(&metadatastoreOpts))
			if err != nil {
				log.Error("init metadata store backend get error", log.Dict{"err": err.Error(), "MetaDataStoreType": s.MetaDataStore_Type})
				os.Exit(1)
			}
		}
	case "redis":
		{
			s.metadatastore = metadatastore.NewRedisMetadataStore()
			err := s.metadatastore.Init(metadatastore.WithOptions(&metadatastoreOpts))
			if err != nil {
				log.Error("init metadata store backend get error", log.Dict{"err": err.Error(), "MetaDataStoreType": s.MetaDataStore_Type})
				os.Exit(1)
			}
		}
	default:
		{
			metadatastoreOpts.URL = s.FilterRedis_URL
			metadatastoreOpts.RedisRouteMod = s.FilterRedis_RouteMod
			metadatastoreOpts.ConnTimeout = s.FilterRedis_ConnTimeout
			s.metadatastore = metadatastore.NewRedisMetadataStore()
			err := s.metadatastore.Init(metadatastore.WithOptions(&metadatastoreOpts))
			if err != nil {
				log.Error("init metadata store backend get error", log.Dict{"err": err.Error(), "MetaDataStoreType": s.MetaDataStore_Type})
				os.Exit(1)
			}
		}
	}
	// 初始化存储计算后端
	s.filterstore = filterstore.NewRedisFilterStore()
	err = s.filterstore.Init(
		filterstore.WithOptions(&filterstore.Options{
			RedisURL:                        s.FilterRedis_URL,
			RedisRouteMod:                   s.FilterRedis_RouteMod,
			ConnTimeout:                     s.FilterRedis_ConnTimeout,
			NamespaceBase:                   s.App_Name,
			BlacklistConfig_KeySubnamespace: s.BlacklistConfig_KeySubnamespace,
			BlacklistConfig_DefaultTTLDays:  s.BlacklistConfig_DefaultTTLDays,
			RangeConfig_KeySubnamespace:     s.RangeConfig_KeySubnamespace,
			NamespaceConfig_EntitySourceTypeSubnamespace:  s.NamespaceConfig_EntitySourceTypeSubnamespace,
			ContextConfig_KeySubnamespace:                 s.ContextConfig_KeySubnamespace,
			NamespaceConfig_PickerIntervalSubnamespace:    s.NamespaceConfig_PickerIntervalSubnamespace,
			NamespaceConfig_PickerSubnamespace:            s.NamespaceConfig_PickerSubnamespace,
			NamespaceConfig_PickerCounterSubnamespace:     s.NamespaceConfig_PickerCounterSubnamespace,
			NamespaceConfig_PickerCounterDateSubnamespace: s.NamespaceConfig_PickerCounterDateSubnamespace,
			DefaultBlacklistFilterInfo:                    defaultBlacklistFilterInfo,
			DefaultRangeFilterInfo:                        defaultRangeFilterInfo,
			DefaultContextFilterInfo:                      s.defaultContextFilterInfo,
			DefaultPickerFilterInfos:                      s.defaultPickerFilterInfos,
			PickerCounterOn:                               s.PickerCounterOn,
			PickerCounterConfig_TTLDays:                   s.PickerCounterConfig_TTLDays,
			PickerCounterConfig_RangeCounterOn:            s.PickerCounterConfig_RangeCounterOn,
			PickerCounterConfig_RangeTTLDays:              s.PickerCounterConfig_RangeTTLDays,
		}))
	if err != nil {
		log.Error("init filter backend get error", log.Dict{"err": err.Error()})
		os.Exit(1)
	}
	/**
	在这里补充其他初始化项(注意去掉注释)
	*/
	s.Run()
}

// PerformanceOpts 配置性能调优设置
func (s *Server) PerformanceOpts() {
	if s.opts == nil {
		s.opts = []grpc.ServerOption{}
	}
	if s.Max_Recv_Msg_Size != 0 {
		s.opts = append(s.opts, grpc.MaxRecvMsgSize(s.Max_Recv_Msg_Size))
	}
	if s.Max_Send_Msg_Size != 0 {
		s.opts = append(s.opts, grpc.MaxSendMsgSize(s.Max_Send_Msg_Size))
	}
	if s.Initial_Window_Size != 0 {
		s.opts = append(s.opts, grpc.InitialWindowSize(int32(s.Initial_Window_Size)))
	}
	if s.Initial_Conn_Window_Size != 0 {
		s.opts = append(s.opts, grpc.InitialConnWindowSize(int32(s.Initial_Conn_Window_Size)))
	}
	if s.Max_Concurrent_Streams != 0 {
		s.opts = append(s.opts, grpc.MaxConcurrentStreams(uint32(s.Max_Concurrent_Streams)))
	}
	if s.Max_Connection_Idle != 0 || s.Max_Connection_Age != 0 || s.Max_Connection_Age_Grace != 0 || s.Keepalive_Time != 0 || s.Keepalive_Timeout != 0 {
		kasp := keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(s.Max_Connection_Idle) * time.Second,
			MaxConnectionAge:      time.Duration(s.Max_Connection_Age) * time.Second,
			MaxConnectionAgeGrace: time.Duration(s.Max_Connection_Age_Grace) * time.Second,
			Time:                  time.Duration(s.Keepalive_Time) * time.Second,
			Timeout:               time.Duration(s.Keepalive_Timeout) * time.Second,
		}
		s.opts = append(s.opts, grpc.KeepaliveParams(kasp))
	}
	if s.Keepalive_Enforcement_Min_Time != 0 || s.Keepalive_Enforcement_Permit_Without_Stream {
		kaep := keepalive.EnforcementPolicy{
			MinTime:             time.Duration(s.Keepalive_Enforcement_Min_Time) * time.Second,
			PermitWithoutStream: s.Keepalive_Enforcement_Permit_Without_Stream,
		}
		s.opts = append(s.opts, grpc.KeepaliveEnforcementPolicy(kaep))
	}
}

// TLSOpts 配置TLS设置
func (s *Server) TLSOpts() {
	if s.opts == nil {
		s.opts = []grpc.ServerOption{}
	}
	if s.Ca_Cert_Path != "" {
		cert, err := tls.LoadX509KeyPair(s.Server_Cert_Path, s.Server_Key_Path)
		if err != nil {
			log.Error("read serv pem file error:", log.Dict{"err": err.Error(), "Cert_path": s.Server_Cert_Path, "Key_Path": s.Server_Key_Path})
			os.Exit(2)
		}
		capool := x509.NewCertPool()
		caCrt, err := os.ReadFile(s.Ca_Cert_Path)
		if err != nil {
			log.Error("read ca pem file error:", log.Dict{"err": err.Error(), "path": s.Ca_Cert_Path})
			os.Exit(2)
		}
		capool.AppendCertsFromPEM(caCrt)
		tlsconf := &tls.Config{
			RootCAs:      capool,
			ClientAuth:   tls.RequireAndVerifyClientCert, // 检验客户端证书
			Certificates: []tls.Certificate{cert},
		}
		if s.Client_Crl_Path != "" {
			clipool := x509.NewCertPool()
			cliCrt, err := os.ReadFile(s.Client_Crl_Path)
			if err != nil {
				log.Error("read pem file error:", log.Dict{"err": err.Error(), "path": s.Client_Crl_Path})
				os.Exit(2)
			}
			clipool.AppendCertsFromPEM(cliCrt)
			tlsconf.ClientCAs = clipool
		}
		creds := credentials.NewTLS(tlsconf)
		s.opts = append(s.opts, grpc.Creds(creds))
	} else {
		creds, err := credentials.NewServerTLSFromFile(s.Server_Cert_Path, s.Server_Key_Path)
		if err != nil {
			log.Error("Failed to Listen as a TLS Server", log.Dict{"error": err.Error()})
			os.Exit(2)
		}
		s.opts = append(s.opts, grpc.Creds(creds))
	}
	log.Info("grpc server will start use TLS")
}

func doNothing() {
	log.Debug("Do nothing")
}

// RegistTools 为指定服务注册各种工具
func (s *Server) RegistTools(serv *grpc.Server) func() {

	// 注册健康检查
	s.healthservice = health.NewServer()
	s.healthservice.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(serv, s.healthservice)

	// 注册反射
	reflection.Register(serv)
	// 注册grpc admin
	if s.Use_Admin {
		cleanup, err := admin.Register(serv)
		if err != nil {
			log.Warn("set grpc admin error", log.Dict{"err": err.Error()})
			return doNothing
		}
		return cleanup
	}
	return doNothing
}

// withLogger 记录gateway的工作log
func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Info("Http log", log.Dict{"code": m.Code, "duration": m.Duration, "Path": request.URL.Path})
	})
}

func handleError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	var newError runtime.HTTPStatusError
	switch err.Error() {
	case "NotFound":
		{
			newError = runtime.HTTPStatusError{
				HTTPStatus: 404,
				Err:        err,
			}
		}
	case "IDExists", "FilterSettingNotMatchFilterType", "NeedID", "NeedEntityType",
		"NeedMetaInfo", "NeedAtLeastOneCandidate", "NeedAtLeastOneEnvInfo",
		"RangeNotExist", "DateOversize", "NeedAtLeastOneEntitySource":
		{
			newError = runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
		}
	case "NotImplemented":
		{
			newError = runtime.HTTPStatusError{
				HTTPStatus: 501,
				Err:        err,
			}
		}
	case "FeatureNotOn":
		{
			newError = runtime.HTTPStatusError{
				HTTPStatus: 503,
				Err:        err,
			}
		}
	default:
		{
			newError = runtime.HTTPStatusError{
				HTTPStatus: 500,
				Err:        err,
			}
		}
	}

	// using default handler to do the rest of heavy lifting of marshaling error and adding headers
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
}

// MakeGateway 构造grpc-gateway的http服务
func (s *Server) MakeGateway(grpc_need_client_cert bool) *http.Server {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(handleError),
	)
	// setting up a dail up for gRPC service by specifying endpoint/target url
	creds := insecure.NewCredentials()
	if grpc_need_client_cert {
		if s.Ca_Cert_Path != "" {
			if s.Gateway_Client_Cert_Path != "" && s.Gateway_Client_Key_Path != "" {
				cert, err := tls.LoadX509KeyPair(s.Gateway_Client_Cert_Path, s.Gateway_Client_Key_Path)
				if err != nil {
					log.Error("read client pem file error:", log.Dict{"err": err.Error(), "Cert_path": s.Gateway_Client_Cert_Path, "Key_Path": s.Gateway_Client_Key_Path})
					os.Exit(1)
				}
				capool := x509.NewCertPool()
				caCrt, err := os.ReadFile(s.Ca_Cert_Path)
				if err != nil {
					log.Error("read ca pem file error:", log.Dict{"err": err.Error(), "path": s.Ca_Cert_Path})
					os.Exit(1)
				}
				capool.AppendCertsFromPEM(caCrt)
				tlsconf := &tls.Config{
					RootCAs:      capool,
					Certificates: []tls.Certificate{cert},
				}
				creds = credentials.NewTLS(tlsconf)

			} else {
				var err error
				creds, err = credentials.NewClientTLSFromFile(s.Ca_Cert_Path, "")
				if err != nil {
					log.Error("failed to load credentials", log.Dict{"err": err.Error()})
					os.Exit(1)
				}
			}
		}
	}
	err := itemfilterRPC_pb.RegisterITEMFILTERRPCHandlerFromEndpoint(context.Background(), mux, s.Address, []grpc.DialOption{grpc.WithTransportCredentials(creds)})
	if err != nil {
		log.Error("Failed to Register GRPCHandlerFromEndpoint", log.Dict{"error": err})
		os.Exit(1)
	}

	m := http.NewServeMux()
	m.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	m.Handle("/", withLogger(mux))
	// Creating a normal HTTP server
	server := http.Server{
		Handler: m,
	}
	return &server
}

// RunXdsServer 启动xds服务服务
func (s *Server) RunXdsServer() {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Error("Failed to Listen", log.Dict{"error": err.Error(), "address": s.Address})
		os.Exit(1)
	}
	hostinfo := strings.Split(s.Address, ":")
	if len(hostinfo) != 2 {
		log.Error("address format not ok", log.Dict{"address": s.Address})
		os.Exit(2)
	}
	port, err := strconv.Atoi(hostinfo[1])
	if err != nil {
		log.Error("address port not int", log.Dict{"address": s.Address})
		os.Exit(2)
	}
	mantenancePort := port + 1
	if s.XDS_Maintenance_Port > 0 && s.XDS_Maintenance_Port != port {
		mantenancePort = s.XDS_Maintenance_Port
	}
	mantenanceAddress := fmt.Sprintf("%s:%d", hostinfo[0], mantenancePort)
	mantenanceLis, err := net.Listen("tcp4", mantenanceAddress)
	if err != nil {
		log.Error("Mantenance Service Failed to Listen", log.Dict{"error": err.Error(), "address": mantenanceAddress})
		os.Exit(1)
	}
	maintenanceServer := grpc.NewServer()
	cleanup := s.RegistTools(maintenanceServer)
	defer cleanup()

	creds := insecure.NewCredentials()
	if s.XDS_CREDS {
		var err error
		creds, err = xdscreds.NewServerCredentials(xdscreds.ServerOptions{FallbackCreds: insecure.NewCredentials()})
		if err != nil {
			log.Error("failed to create server-side xDS credentials", log.Dict{"error": err.Error()})
			os.Exit(2)
		}
		log.Info("server will start use XDS_CREDS")
	}
	s.opts = append(s.opts, grpc.Creds(creds))

	gs := xds.NewGRPCServer(s.opts...)
	defer gs.Stop()

	// 注册服务
	itemfilterRPC_pb.RegisterITEMFILTERRPCServer(gs, s)

	if s.Use_Gateway {
		var server *http.Server
		if s.XDS_CREDS {
			server = s.MakeGateway(true)
		} else {
			server = s.MakeGateway(false)
		}
		// 设置分流,http1走restful,http2走grpc
		m := cmux.New(lis)

		// a different listener for HTTP2 since gRPC uses HTTP2
		grpcL := m.Match(cmux.HTTP2())
		// a different listener for HTTP1
		httpL := m.Match(cmux.HTTP1Fast())
		// 构造gateway用户处理http请求

		// 启动服务
		log.Info("Server Start with RESTful and Grpc api and mantenance by XDS", log.Dict{"address": s.Address, "mantenance_port": mantenancePort})
		go func() {
			if s.Gateway_Use_TLS {
				err := server.ServeTLS(httpL, s.Server_Cert_Path, s.Server_Key_Path)
				if err != nil {
					if err.Error() != "mux: server closed" {
						log.Error("Failed to Serve Http", log.Dict{"error": err})
						os.Exit(1)
					}
				}
			} else {
				err := server.Serve(httpL)
				if err != nil {
					if err.Error() != "mux: server closed" {
						log.Error("Failed to Serve Http", log.Dict{"error": err})
						os.Exit(1)
					}
				}
			}
		}()
		go func() {
			err := gs.Serve(grpcL)
			if err != nil {
				log.Error("Failed to Serve Grpc", log.Dict{"error": err})
				os.Exit(1)
			}
		}()
		go func() {
			err := maintenanceServer.Serve(mantenanceLis)
			if err != nil {
				log.Error("Failed to Start Maintenance Server", log.Dict{"error": err})
				os.Exit(1)
			}
		}()
		go func() {
			err := m.Serve()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") && err.Error() != "mux: server closed" {
					log.Error("Failed to Start cmux", log.Dict{"error": err})
					os.Exit(1)
				}
			}
		}()
		// 等待中断信号以优雅地关闭服务器（设置 3 秒的超时时间）
		quit := make(chan os.Signal, 3)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Info("Shutdown Server ...")
		m.Close()
		gs.GracefulStop()
		maintenanceServer.GracefulStop()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = server.Shutdown(ctx)
		if err != nil {
			log.Error("server shutdown error", log.Dict{"err": err.Error()})
			os.Exit(2)
		}
		log.Info("Server Shutdown")
	} else {
		log.Info("Server Start with Grpc api and mantenance by XDS", log.Dict{"address": s.Address, "mantenance_port": mantenancePort})
		go func() {
			err := gs.Serve(lis)
			if err != nil {
				log.Error("Failed to Serve Grpc", log.Dict{"error": err})
				os.Exit(1)
			}
		}()
		go func() {
			err := maintenanceServer.Serve(mantenanceLis)
			if err != nil {
				log.Error("Failed to Start Maintenance Server", log.Dict{"error": err})
				os.Exit(1)
			}
		}()
		// 等待中断信号以优雅地关闭服务器（设置 3 秒的超时时间）
		quit := make(chan os.Signal, 3)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Info("Shutdown Server ...")
		gs.GracefulStop()
		maintenanceServer.GracefulStop()
		log.Info("Server Shutdown")
	}
}

// RunCommonServer 启动常规服务服务
func (s *Server) RunCommonServer() {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Error("Failed to Listen", log.Dict{"error": err.Error(), "address": s.Address})
		os.Exit(1)
	}

	if s.Server_Cert_Path != "" && s.Server_Key_Path != "" {
		s.TLSOpts()
	}
	gs := grpc.NewServer(s.opts...)
	defer gs.Stop()
	cleanup := s.RegistTools(gs)
	defer cleanup()

	// 注册服务
	itemfilterRPC_pb.RegisterITEMFILTERRPCServer(gs, s)

	if s.Use_Gateway {
		var server *http.Server
		if s.Server_Cert_Path != "" && s.Server_Key_Path != "" {
			server = s.MakeGateway(true)
		} else {
			server = s.MakeGateway(false)
		}
		// 构造gateway用户处理http请求
		// 设置分流,http1走restful,http2走grpc
		m := cmux.New(lis)
		// a different listener for HTTP2 since gRPC uses HTTP2
		grpcL := m.Match(cmux.HTTP2())
		// a different listener for HTTP1
		httpL := m.Match(cmux.HTTP1Fast())
		// 启动服务
		log.Info("Server Start with RESTful and Grpc api", log.Dict{"address": s.Address})
		go func() {
			err := gs.Serve(grpcL)
			if err != nil {
				log.Error("Failed to Serve Grpc", log.Dict{"error": err})
				os.Exit(1)
			}

		}()
		go func() {
			if s.Gateway_Use_TLS {
				err := server.ServeTLS(httpL, s.Server_Cert_Path, s.Server_Key_Path)
				if err != nil {
					if err.Error() != "mux: server closed" {
						log.Error("Failed to Serve Http", log.Dict{"error": err})
						os.Exit(1)
					}
				}
			} else {
				err := server.Serve(httpL)
				if err != nil {
					if err.Error() != "mux: server closed" {
						log.Error("Failed to Serve Http", log.Dict{"error": err})
						os.Exit(1)
					}
				}
			}
		}()

		go func() {
			err := m.Serve()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") && err.Error() != "mux: server closed" {
					log.Error("Failed to Start cmux", log.Dict{"error": err})
					os.Exit(1)
				}
			}
		}()
		// 等待中断信号以优雅地关闭服务器（设置 3 秒的超时时间）
		quit := make(chan os.Signal, 3)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Info("Shutdown Server ...")
		m.Close()
		gs.GracefulStop()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = server.Shutdown(ctx)
		if err != nil {
			log.Error("server shutdown error", log.Dict{"err": err.Error()})
			os.Exit(2)
		}

		log.Info("Server Shutdown")

	} else {
		// 启动服务
		log.Info("Server Start with Grpc api", log.Dict{"address": s.Address})
		go func() {
			err := gs.Serve(lis)
			if err != nil {
				log.Error("Failed to Serve Grpc", log.Dict{"error": err})
				os.Exit(1)
			}
		}()
		// 等待中断信号以优雅地关闭服务器（设置 3 秒的超时时间）
		quit := make(chan os.Signal, 3)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Info("Shutdown Server ...")
		gs.GracefulStop()
		log.Info("Server Shutdown")
	}

}

// RegistInterceptor 注册拦截器
func (s *Server) RegistInterceptor() {
	if s.opts == nil {
		s.opts = []grpc.ServerOption{}
	}
	if len(s.UnaryInterceptors) > 0 {
		for _, i := range s.UnaryInterceptors {
			s.opts = append(s.opts, grpc.UnaryInterceptor(i))
		}
	}
	if len(s.StreamInterceptors) > 0 {
		for _, i := range s.StreamInterceptors {
			s.opts = append(s.opts, grpc.StreamInterceptor(i))
		}
	}
}

// Run 执行grpc服务
func (s *Server) Run() {
	// 注册grpc的配置
	s.PerformanceOpts()
	// 注册拦截器
	s.RegistInterceptor()
	if s.XDS {
		log.Info("setup as xds grpc serv")
		s.RunXdsServer()
	} else {
		log.Info("set up as common grpc serv")
		s.RunCommonServer()
	}
}

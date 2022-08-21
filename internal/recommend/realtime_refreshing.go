package recommend

const RealtimeStatus_ON = 1
const RealtimeStatus_OFF = -1

type Realtime struct {
	Enable              int32
	EnableActualTotal   int32
	ResetIfOffsetIsZero int32
	DedupCaches         map[string]*RealtimeDedupCache
}

type RealtimeDedupCache struct {
	DataType       string
	Cap            int
	TTL            int
	SharedSections []string

	Key               string
	SharedSectionKeys []string
	Data              map[string][]*RealtimeDedupCacheUnit
}

type RealtimeDedupCacheUnit struct {
	ShopId  int64
	ItemId  int64
	AdsId   int64
	KnodeId string
	VideoId int64
	DishId  int64
}

func InitRealtime(section *Section) error {
	return nil
}

func FinishRealtime(section *Section) {

}

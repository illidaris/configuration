package configuration

import (
	"math/rand"
	"sort"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
)

// IChooseFunc 是一个函数类型，用于从多个model.Instance中选择一个。
// 参数instances是可变长度的model.Instance类型切片。
// 返回值是一个model.Instance类型，表示选择的结果。
type IChooseFunc func(instances ...model.Instance) model.Instance

// InstanceSlice 是model.Instance类型切片的自定义类型。
type InstanceSlice []model.Instance

// Len 返回InstanceSlice的长度。
func (a InstanceSlice) Len() int {
	return len(a)
}

// Swap 交换InstanceSlice中索引i和j位置的元素。
func (a InstanceSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less 判断InstanceSlice中索引i位置的元素是否小于索引j位置的元素。
// 此方法用于实例的排序，基于实例的Weight属性。
func (a InstanceSlice) Less(i, j int) bool {
	return a[i].Weight < a[j].Weight
}

// defaultChoose 是一个默认的选择函数，基于提供的model.Instance实例列表进行选择。
// 参数is是可变长度的model.Instance类型切片。
// 返回值是一个model.Instance类型，表示选择的结果。
func defaultChoose(is ...model.Instance) model.Instance {
	return newChooser(is).Pick()
}

// Chooser 结构体用于存储实例列表及其累计权重，以便进行随机选择。
type Chooser struct {
	data   []model.Instance // 实例数据
	totals []int            // 实例权重的累计值
	max    int              // 权重的总和
}

// NewChooser 初始化一个新的Chooser实例，用于从提供的实例列表中进行选择。
// 参数instances是model.Instance类型的切片。
// 返回值是一个已经初始化好的Chooser实例。
func newChooser(instances []model.Instance) Chooser {
	// 初始化时对实例列表按权重进行排序
	sort.Sort(InstanceSlice(instances))
	totals := make([]int, len(instances)) // 创建累计权重数组
	runningTotal := 0                     // 运行时总权重
	for i, c := range instances {
		runningTotal += int(c.Weight) // 计算每个实例结束时的累计权重
		totals[i] = runningTotal      // 存储累计权重
	}
	return Chooser{data: instances, totals: totals, max: runningTotal}
}

// Pick 从Chooser中随机选择一个实例。
// 返回值是一个model.Instance类型，表示选择的结果。
func (chs Chooser) Pick() model.Instance {
	r := rand.Intn(chs.max) + 1         // 随机生成一个在1和最大权重之间的数
	i := sort.SearchInts(chs.totals, r) // 使用二分查找找到该权重对应的实例索引
	return chs.data[i]                  // 返回选定的实例
}

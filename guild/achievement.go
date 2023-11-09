package guild

// 成就：公会热度
func AGuildActiviy(num int) (res string) {
	switch true {
	case num > 15:
		res = "热度/门庭若市"
	case num > 8:
		res = "热度/小有热闹"
	case num > 3:
		res = "热度/结伴而行"
	default:
		res = "热度/冷冷清清"
	}
	return
}

// 成就：激励分配分散度
func AFairDistribution(per float64) (res string) {
	switch true {
	case per > 0.8:
		res = "分配/超集中"
	case per > 0.5:
		res = "分配/个别集中"
	default:
		res = ""
	}
	return
}

// 成就：文章阅读量
func AReadership(hits map[int]int) (res string) {
	switch true {
	case hits[5000] >= 1:
		res = "阅读/翻页高手"
	case hits[1000] >= 2:
		res = "阅读/风向标"
	case hits[500] >= 2:
		res = "阅读/引发兴趣"
	case hits[500] >= 1:
		res = "阅读/注意力寻求"
	default:
		res = "阅读/无人问津"
	}
	return
}

// 成就：媒体精选文章
func AMediaPicks(frontPages int) (res string) {
	switch true {
	case frontPages >= 4:
		res = "精选/群星灿烂"
	case frontPages >= 3:
		res = "精选/三星成簇"
	case frontPages >= 2:
		res = "精选/双星辉映"
	case frontPages >= 1:
		res = "精选/孤星闪烁"
	default:
		res = ""
	}
	return
}

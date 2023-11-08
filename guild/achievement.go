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

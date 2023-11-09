package schema

var Guilds = map[string]GuildMetaInfo{
	"内容公会 - 策划": GuildMetaInfo{
		NID:    "e9d5bbfe68a6457cb8ccf7e4badeaed1",
		FinNID: "328f2bfbfdbe4f9581af37f393893e36",
		Info:   "策划组致力于推动 Arweave 的创新与发展。我们汇集了视频、策展、条漫、播客等方面的创作者，通过深入了解 Arweave 技术和市场需求，将创新项目转化为可行方案！",
		Rank:   1,
	},
	"内容公会 - 翻译": GuildMetaInfo{
		NID:    "e9d5bbfe68a6457cb8ccf7e4badeaed1",
		FinNID: "e8d79c55c0394cba83664f3e5737b0bd",
		Info:   "翻译小组将重点关注 Web3、Arweave、SCP（存储共识范式）及其他相关主题的翻译研究，目标是与世界分享来自不同来源的高质量想法、深度思考和变革性技术。",
		Rank:   2,
	},
	"内容公会 - 投稿": GuildMetaInfo{
		NID:    "e9d5bbfe68a6457cb8ccf7e4badeaed1",
		FinNID: "a815dcd96395424a93d9854b4418ab03",
		Info:   "这里是 PermaDAO 文化传递的桥梁，你可以在这里投稿生态相关文章，比内容更多的是理解，比理解更多的是热爱。欢迎加入内容公会！",
		Rank:   3,
	},
	"品宣公会": GuildMetaInfo{
		NID:    "937f66c91e534c1187f0570ffe6e2b65",
		FinNID: "990db3313e42412b8c6ab07e399a2635",
		Info:   "通过多种宣传推广渠道为 Arweave 生态打响品牌，欢迎加入，这里是 PermaDAO 宣传组！",
		Rank:   4,
	},
	"活动公会": GuildMetaInfo{
		NID:    "285824ecde1349d6b23853a6a5f6685d",
		FinNID: "f2160eae42e9483882f01d3daa7090fa",
		Info:   "线上活动、嘉宾邀请，这里可以发挥你的创意与执行力，用活力来带动气氛！发挥你的想象力，为 Arweave 生态添砖加瓦，让我们一起“动”起来🏃‍♀️🏃！",
		Rank:   5,
	},
	"开发公会": GuildMetaInfo{
		NID:    "6fa2a0760efd45ff9b6dd1731443e4b2",
		FinNID: "146e1f661ed943e3a460b8cf12334b7b",
		Info:   "加入开发小组，你可以学习 Arweave 生态开发，你可以获得 Arweave 开发任务，还可以提出你的开发新点子✨！社区需要什么，我们就构建什么！",
		Rank:   6,
	},
	"PSPC - Market": GuildMetaInfo{
		NID:    "69ba28d2d17643ae9711947329138c58",
		FinNID: "a9ce0c5902b14e4891ed0fb6333a9e92",
		Info:   "Arweave 生态唯一的跨链 DEX，支持 Arweave 生态分润代币 PST 兑换！欢迎加入 Permaswap 社区运营！ 🎉",
		Rank:   7,
	},
	"PSPC - Product": GuildMetaInfo{
		NID:    "69ba28d2d17643ae9711947329138c58",
		FinNID: "27555aec8d734b6889ae1836d7a67b4a",
		Info:   "Arweave 生态唯一的跨链 DEX，支持 Arweave 生态分润代币 PST 兑换！欢迎加入产品构建！ 🎉",
		Rank:   8,
	},
	"管理公会": GuildMetaInfo{
		NID:    "4a8b5461c3b241659e95b0fb3c174250",
		FinNID: "caac7a1aefcc4ed0b02b8adbc106f021",
		Info:   "PermaDAO 管理。DAO 治理提案发起，DAO 页面的构建，DAO 工作的梳理，以及每周任务结算管理等。",
		Rank:   9,
	},
}

type GuildMetaInfo struct {
	NID    string
	FinNID string
	Info   string
	Rank   float64
}

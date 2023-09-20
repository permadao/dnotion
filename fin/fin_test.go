package fin

import "testing"

func TestDNotion(t *testing.T) {
	// dao := New(
	// 	"notion_secret",
	// 	"private_key",
	// 	"https://api.everpay.io",
	// 	[]string{
	// 		"f0fdcea614544768b5756a3fde23ff51", // content
	// 		"2f6d7e79ce8c4c3a83bdc5b4c4672998", // translation
	// 		"34f6db380ccc430086ab54054c34da44", // submission
	// 		"bd2d60313fa54c0190823942f4de4bd6", // promotion
	// 		"3a5cb7bb5d2a4990ba3d232597c7c26f", // activity
	// 		"d1aa3bbffabf41e2ab1bd14f59215e11", // admin
	// 		"117afdd004fc4890a60d586e39910883", // dev
	// 	},
	// 	[]string{
	// 		"3352862db68a494d80fd902af4f50e05", // content
	// 		"0b6288a1d3984b7a809adfb01c20a194", // translation
	// 		"ad2cf585b08843fea7cf40a682bc4529", // submission
	// 		"59e839717bab45c68bb2fa25aba9020b", // promotion
	// 		"28b4a0d7e78346f69f6a859386b38225", // activity
	// 		"6dd7a8c678ed4f42be27713cd11f147b", // admin
	// 		"4c4d512c4c9e4e09b951967752f185ea", // dev
	// 		"179790222580482f9e691ab6618f473e", // psp market
	// 		"8cd2d839823a42daa1c8ae5b104c67c7", // psp prod
	// 	},
	// 	[]string{
	// 		"328f2bfbfdbe4f9581af37f393893e36", // content
	// 		"e8d79c55c0394cba83664f3e5737b0bd", // translation
	// 		"a815dcd96395424a93d9854b4418ab03", // submission
	// 		"990db3313e42412b8c6ab07e399a2635", // promotion
	// 		"f2160eae42e9483882f01d3daa7090fa", // activity
	// 		"caac7a1aefcc4ed0b02b8adbc106f021", // admin
	// 		"146e1f661ed943e3a460b8cf12334b7b", // dev
	// 		"a9ce0c5902b14e4891ed0fb6333a9e92", // psp market
	// 		"27555aec8d734b6889ae1836d7a67b4a", // psp prod
	// 	},
	// 	"d808ed8b25d746beb267ad552b1d1bf5", // contributors
	// )

	// dao.InitContributors()

	// // 1. check counts
	// dao.CheckAllDbsCountAndID()

	// // 2. check actual usd equality with workload usd
	// dao.CheckAllWorkloadAndAmount()

	// // 3. Update workload tx to finance table
	// dao.UpdateAllWorkToFin()
	// dao.UpdateWorkToFin("ef2043f45b0144ec846f69cd035c1224", "0c8f5483e1344d919e5e5a49d6d8dabb")

	// // 4. Update all Finance transactions to progress
	// dao.UpdateAllFinToProgress("2023-07-21", "AR", 6.11, "AR", 1)
	// dao.UpdateFinToProgress("0c8f5483e1344d919e5e5a49d6d8dabb", "2023-09-01", "AR", 4.17, "AR", 1)

	// // 5. process payment by everpay
	// dao.PayAll()
	// dao.Pay("0c8f5483e1344d919e5e5a49d6d8dabb")
}

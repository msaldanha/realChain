package main

import "github.com/msaldanha/realChain/cmd"

func main() {
	command := cmd.New()
	command.Execute()

	//net, err := network.NewNetwork("127.0.0.1:4000")
	//if err != nil {
	//	panic("Failed to setup network")
	//}
	//net.InstallHandler("end.point", func(ctx *network.Context) {
	//	fmt.Printf("HANDLER: Reveived data: %s\n", string(ctx.Data))
	//})
	//net.UsePeers("127.0.0.1:4001", "127.0.0.1:4002")
	//
	//net2, err := network.NewNetwork("127.0.0.1:4001")
	//if err != nil {
	//	panic("Failed to setup network")
	//}
	//net2.UsePeers("127.0.0.1:4000", "127.0.0.1:4002")
	////net2.UsePeers("127.0.0.1:4000")
	//
	//net3, err := network.NewNetwork("127.0.0.1:4002")
	//if err != nil {
	//	panic("Failed to setup network")
	//}
	//net3.UsePeers("127.0.0.1:4000", "127.0.0.1:4001")
	//
	//go net.Run()
	//go net2.Run()
	//go net3.Run()
	//
	//for {
	//	time.Sleep(5000 * time.Millisecond)
	//	net.Broadcast("end.point", []byte("data payload from 1"))
	//	net2.Broadcast("end.point", []byte("data payload from 2"))
	//	net2.Broadcast("end.point", []byte("data payload from 3"))
	//}

}
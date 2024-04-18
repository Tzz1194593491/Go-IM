package main

//程序主入口
func main() {
	server := Server{ip: "127.0.0.1", port: 8080}
	server.Start()
}

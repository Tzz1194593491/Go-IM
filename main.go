package main

// 程序主入口
func main() {
	server := NewServer("0.0.0.0", 8888)
	server.Start()
}

package main

const listenAddr = "127.0.0.1:8080"

func main() {
	srv := NewAPIServer(listenAddr)
	srv.Run()
}

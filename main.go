package main

const listenAddr = "127.0.0.1:3000"

func main() {
	store := NewMySqlStore()
	store.Init()

	srv := NewAPIServer(store, listenAddr)
	srv.Run()
}

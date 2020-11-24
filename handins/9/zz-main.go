package main

func main() {
	storedData := "zz-init-data.json"
	GenerateAndDump(5, "zz-init-data.json")
	accs := LoadFromFile(storedData)

	// run 5 of them with different private keys and send data
	model := InitModel(&accs.PrivateKeys[0], &accs)
	StartupServer(model)
}

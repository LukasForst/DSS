package main

func main() {
	storedData := "zz-init-data.json"
	//GenerateAndDump(5, "zz-init-data.json")
	accs := LoadFromFile(storedData)

	// todo choose one
	model := InitModel(&accs.PrivateKeys[0], &accs)
	StartupServer(model)
}

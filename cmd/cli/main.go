package main

import "os"

func main() {
	if len(os.Args) < 2 {
		println("expected 'manifest', 'container' or 'token'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "manifest":
		handleManifest(os.Args[2:])
	case "container":
		handleContainer(os.Args[2:])
	case "token":
		handleToken(os.Args[2:])
	default:
		println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}

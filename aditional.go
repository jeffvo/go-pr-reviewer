package main

type unexportedType string

func ExportedFunc() unexportedType {
	return unexportedType("some string")
}

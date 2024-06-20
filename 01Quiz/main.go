package main

import (
	"fmt"
	"strings"
)
func main(){
	fmt.Println("Welcome to my quiz game!")
	var name string

	fmt.Printf("Enter your name: ")
	fmt.Scan(&name)
	fmt.Printf("Hello %v, welcome to the game!",name);

	var age uint
	fmt.Printf("\nEnter you age: ")
	fmt.Scan(&age)

	if age>=10{
		fmt.Println("Yay! You can play!")
	} else {
		fmt.Println("Oops! You can't play!")
		return
	}

	var score int=0

	fmt.Println("First Question-")
	fmt.Printf("What is better, the RTX 3080 or RTX 3090?")
	var answer1 string
	var answer2 string
	fmt.Scan(&answer1,&answer2)

	var answer string = "RTX 3090"
	myans:= answer1 +" " + answer2

	if strings.ToLower(myans) == strings.ToLower(answer){
		fmt.Println("Correct!")
		score++
	} else {
		fmt.Println("Incorrect!")
	}

	fmt.Println("Second Question-")

	var mycores int
	fmt.Printf("How many cores does the Ryzen 9 3900x have? ")
	fmt.Scan(&mycores)

	if(mycores==12){
		fmt.Println("Correct!")
		score++
	} else{
		fmt.Println("Incorrect!")
	}

	fmt.Printf("You scored %v out of 2.\n",score)
	percent:= (float64(score)/float64(2))*100
	fmt.Printf("You score %v%%",percent)
}
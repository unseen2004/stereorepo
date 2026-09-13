package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Exercise struct {
	Name string
	Dir  string
	Mode string
	Hint string
}

// counts as completed, mirroring rustlings' "I AM NOT DONE".
const notDoneMarker = "I AM NOT DONE"

var exercises = []Exercise{
	{
		Name: "hello1",
		Dir:  "hello/hello1",
		Mode: "run",
		Hint: `Take a look at the Println call in main(). It prints "Goodbye!".
Change the string so the program prints "Hello World!" instead, then
remove the "I AM NOT DONE" comment.`,
	},
	{
		Name: "variables1",
		Dir:  "variables/variables1",
		Mode: "run",
		Hint: `The program prints x, but x was never declared.

Declare it before use, for example with a short variable declaration:
    x := 5`,
	},
	{
		Name: "variables2",
		Dir:  "variables/variables2",
		Mode: "run",
		Hint: `You cannot combine ` + "`var`" + ` with ` + "`:=`" + `. ` + "`:=`" + ` declares AND assigns in one step.

Either use a full declaration:
    var x int = 5
or drop ` + "`var`" + ` entirely:
    x := 5`,
	},
	{
		Name: "variables3",
		Dir:  "variables/variables3",
		Mode: "run",
		Hint: `Constants cannot be reassigned after they are declared.

x is meant to change from 5 to 10, so it should be a variable, not a
constant. Replace ` + "`const`" + ` with ` + "`var`" + `.`,
	},
	{
		Name: "variables4",
		Dir:  "variables/variables4",
		Mode: "run",
		Hint: `The inner block declares a NEW variable x with ` + "`x := 2`" + `, which
shadows the outer one and is discarded when the block ends.

To change the OUTER x you need an assignment instead of a declaration:
    x = 2`,
	},
	{
		Name: "variables5",
		Dir:  "variables/variables5",
		Mode: "test",
		Hint: `Go variables declared without an initializer get their type's ZERO value:
    0 for numbers, false for bools, "" for strings, nil for pointers and slices.

Declare the four variables at the top of the test function with ` + "`var`" + `,
one per line, e.g.:
    var i int`,
	},
	{
		Name: "types1",
		Dir:  "types/types1",
		Mode: "run",
		Hint: `Integer division truncates: 5 / 2 is 2, not 2.5.

Make at least one operand a floating point number:
    fmt.Println(5.0 / 2.0)`,
	},
	{
		Name: "types2",
		Dir:  "types/types2",
		Mode: "run",
		Hint: `Go does not implicitly convert between types, so you cannot
concatenate a string and an int with +.

Convert n to a string first with strconv.Itoa(n), or use:
    fmt.Sprintf("The number is %d", n)`,
	},
	{
		Name: "types3",
		Dir:  "types/types3",
		Mode: "test",
		Hint: `Use explicit conversions in Go:

    float64(n)   converts an int to a float64
    int(f)       converts a float64 to an int (truncating towards zero)`,
	},
	{
		Name: "functions1",
		Dir:  "functions/functions1",
		Mode: "test",
		Hint: `The function signature already takes two ints and returns an int.
Just return the sum:

    return a + b`,
	},
	{
		Name: "functions2",
		Dir:  "functions/functions2",
		Mode: "test",
		Hint: `Go functions can return multiple values. The / and % operators give
you the quotient and remainder in one line:

    return a / b, a % b`,
	},
	{
		Name: "functions3",
		Dir:  "functions/functions3",
		Mode: "test",
		Hint: `` + "`nums ...int`" + ` collects all arguments into a slice called nums.
Loop over it with ` + "`for _, n := range nums`" + ` and add each n to total.

Because total is a NAMED return value, a bare ` + "`return`" + ` returns it.`,
	},
	{
		Name: "functions4",
		Dir:  "functions/functions4",
		Mode: "test",
		Hint: `The counter variable must live inside makeCounter so every call
to the returned function sees the same variable.

The returned closure should increment count and then return it:
    return func() int {
        count++
        return count
    }`,
	},
	{
		Name: "controlflow1",
		Dir:  "controlflow/controlflow1",
		Mode: "test",
		Hint: `Chain if / else if / else:

    if score >= 90 {
        return "A"
    } else if score >= 80 {
        return "B"
    } else if score >= 70 {
        return "C"
    }
    return "F"`,
	},
	{
		Name: "controlflow2",
		Dir:  "controlflow/controlflow2",
		Mode: "test",
		Hint: `Go's for loop:

    sum := 0
    for i := 1; i <= n; i++ {
        sum += i
    }
    return sum

When n is 0 the loop body never runs and you return 0.`,
	},
	{
		Name: "controlflow3",
		Dir:  "controlflow/controlflow3",
		Mode: "test",
		Hint: `Range over the slice and add each element to a running total:

    sum := 0
    for _, n := range nums {
        sum += n
    }
    return sum`,
	},
	{
		Name: "strings1",
		Dir:  "strings/strings1",
		Mode: "test",
		Hint: `Convert the string to a slice of runes with []rune(s) so multi-byte
characters count once each.

Then use a for loop that swaps runes[i] and runes[len(runes)-1-i] until
i reaches the middle, and return string(runes).`,
	},
	{
		Name: "strings2",
		Dir:  "strings/strings2",
		Mode: "test",
		Hint: `Use strings.ContainsRune to check each character:

    for _, r := range strings.ToLower(s) {
        if strings.ContainsRune("aeiou", r) {
            count++
        }
    }
    return count`,
	},
	{
		Name: "structs1",
		Dir:  "structs/structs1",
		Mode: "test",
		Hint: `Fill in the missing fields and construct the value:

    type person struct {
        name string
        age  int
    }

    return person{name: name, age: age}`,
	},
	{
		Name: "structs2",
		Dir:  "structs/structs2",
		Mode: "test",
		Hint: `A method is a function with a receiver between func and the name:

    func (r rectangle) area() float64 {
        return r.width * r.height
    }`,
	},
	{
		Name: "slices1",
		Dir:  "slices/slices1",
		Mode: "test",
		Hint: `Start with nil and append in a loop:

    var out []int
    for i := 1; i <= n; i++ {
        out = append(out, i)
    }
    return out

append returns a new slice, so always assign the result back.`,
	},
	{
		Name: "slices2",
		Dir:  "slices/slices2",
		Mode: "test",
		Hint: `Slices share their backing array, so to get an independent copy you
must allocate new memory:

    dst := make([]int, len(src))
    copy(dst, src)
    return dst`,
	},
	{
		Name: "maps1",
		Dir:  "maps/maps1",
		Mode: "test",
		Hint: `Create the map with make, then count:

    counts := make(map[string]int)
    for _, w := range words {
        counts[w]++
    }
    return counts

Reading a missing key gives the zero value (0), so counts[w]++ is safe.`,
	},
	{
		Name: "pointers1",
		Dir:  "pointers/pointers1",
		Mode: "test",
		Hint: `p holds the ADDRESS of an int. Dereference it with *p to read or
write the value it points to:

    *p = *p + n`,
	},
	{
		Name: "interfaces1",
		Dir:  "interfaces/interfaces1",
		Mode: "test",
		Hint: `A type satisfies an interface by implementing its methods - no
"implements" keyword needed:

    func (d dog) speak() string { return "Woof" }
    func (c cat) speak() string { return "Meow" }`,
	},
	{
		Name: "errors1",
		Dir:  "errors/errors1",
		Mode: "test",
		Hint: `Check the divisor before dividing:

    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil

The errors package is already imported for you.`,
	},
	{
		Name: "errors2",
		Dir:  "errors/errors2",
		Mode: "test",
		Hint: `Return a descriptive error for negative inputs using fmt.Errorf:

    if x < 0 {
        return 0, fmt.Errorf("cannot take the square root of %v", x)
    }`,
	},
	{
		Name: "generics1",
		Dir:  "generics/generics1",
		Mode: "test",
		Hint: `The type parameter [T int | float64] means T can be either int or
float64, and comparing with > works for both:

    if a > b {
        return a
    }
    return b`,
	},
	{
		Name: "channels1",
		Dir:  "channels/channels1",
		Mode: "test",
		Hint: `Send a value on a channel with the <- operator:

    c <- n`,
	},
	{
		Name: "channels2",
		Dir:  "channels/channels2",
		Mode: "test",
		Hint: `Start a goroutine that computes the sum and sends it on a channel,
then receive and return it:

    c := make(chan int)
    go func() {
        sum := 0
        for _, n := range nums {
            sum += n
        }
        c <- sum
    }()
    return <-c

Sending/receiving on an unbuffered channel blocks until both sides are
ready - that is what synchronizes the goroutine with main.`,
	},
}

func findExercise(name string) (*Exercise, bool) {
	for i := range exercises {
		if exercises[i].Name == name {
			return &exercises[i], true
		}
	}
	return nil, false
}

func findRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod - run golings from inside the repo")
		}
		dir = parent
	}
}

func mustRoot() string {
	root, err := findRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "golings: %v\n", err)
		os.Exit(1)
	}
	return root
}

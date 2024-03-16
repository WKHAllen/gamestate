# gamestate

A dead simple state mapping system designed to facilitate writing stateful text-based games in Go.

## Usage

The state mapping system enables users to create a series of states that the game can be in. Transitions can be defined to control which states map to each other under certain conditions. Values can be saved between states via the built in global key-value map.

```go
question := gamestate.NewGameState(func() any {
    fmt.Printf("What is your favorite programming language?\n> ")
    language, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    return strings.TrimSpace(language)
})

answerGo := gamestate.NewGameState(func() any {
    fmt.Printf("Go is a really fun language. What about Go makes it your favorite?\n> ")
    reason, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    gamestate.SetGameValue("reason", strings.TrimSpace(reason))
    return nil
})

answerRust := gamestate.NewGameState(func() any {
    fmt.Printf("Rust is a very intelligently designed language. What about Rust makes it your favorite?\n> ")
    reason, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    gamestate.SetGameValue("reason", strings.TrimSpace(reason))
    return nil
})

answerPython := gamestate.NewGameState(func() any {
    fmt.Printf("Python makes coding very simple. What about Python makes it your favorite?\n> ")
    reason, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    gamestate.SetGameValue("reason", strings.TrimSpace(reason))
    return nil
})

fallbackAnswer := gamestate.NewGameState(func() any {
    fmt.Println("I'm not familiar with that language.")
    return nil
})

end := gamestate.NewGameState(func() any {
    if reason, ok := gamestate.GetGameValue("reason"); ok {
        fmt.Println("Your reasoning for your favorite language:", reason)
    }
    fmt.Println("Thanks for playing!")
    return nil
})

question.TransitionOn(answerGo, "Go")
question.TransitionOn(answerRust, "Rust")
question.TransitionOn(answerPython, "Python")
question.Transition(fallbackAnswer)
answerGo.Transition(end)
answerRust.Transition(end)
answerPython.Transition(end)
fallbackAnswer.Transition(end)

gamestate.RunGame(question, end)
```

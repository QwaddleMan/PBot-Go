package main

import (
    "fmt"
)

type Quote struct {
    Name string
    Quote string
}

func (quote *Quote) toString() string {
    return fmt.Sprintf("\"%s\"\n- %s", quote.Name, quote.Quote);
}

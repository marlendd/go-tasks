package main

import (
	"fmt"
	"math/rand"
)

/// начало решения

// генерит случайные слова из 5 букв
// с помощью randomWord(5)
func generate(cancel <-chan struct{}) <-chan string{
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case out <- randomWord(5):
			case <-cancel: return
			}
		}
	}()
	return out
}

// выбирает слова, в которых не повторяются буквы,
// abcde - подходит
// abcda - не подходит
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case word, ok := <- in:
				if !ok {
					return
				}
				seen := map[rune]struct{}{}
				flag := false
				for _, letter := range word {
					if _, exists := seen[letter]; exists{
						flag = true
						break
					}
					seen[letter] = struct{}{}
				}
				if !flag {
					select {
					case out <- word:
					case <-cancel:
						return
					}
				}
			case <-cancel: return
			}
		}
	}()
	return out
}

// переворачивает слова
// abcde -> edcba
func reverse(cancel <-chan struct{}, in <-chan string) <-chan string{
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case word, ok := <-in:
				if !ok {
					return
				}
				reverseWord := func(s string) string {
					runeStr := []rune(s)
					reversed := runeStr
					for _, rune := range " -> " {
						reversed = append(reversed, rune)
					}
					for i := len(runeStr)-1; i >= 0; i-- {
						reversed = append(reversed, runeStr[i])
					}
					return string(reversed)
				}
				select {
				case out <- reverseWord(word):
				case <-cancel:
					return
				}
			case <-cancel: return 
			}
		}
	}()
	return out
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for c1 != nil || c2 != nil {
            select {
            case v, ok := <-c1:
                if !ok {
                    c1 = nil
                    continue
                }
                select {
                case out <- v:
                case <-cancel:
                    return
                }
            case v, ok := <-c2:
                if !ok {
                    c2 = nil
                    continue
                }
                select {
                case out <- v:
                case <-cancel:
                    return
                }
            case <-cancel:
                return
            }
        }
    }()
    return out
}
// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan string, n int) {
	for range(n) {
		select {
		case val, ok := <- in:
			if !ok {
				return
			}
			fmt.Println(val)
		case <- cancel: return
		}
	}
}

// конец решения

// генерит случайное слово из n букв
func randomWord(n int) string {
	const letters = "aeiourtnsl"
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 10)
}

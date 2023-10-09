package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//	var (
//		fileName = fmt.Sprintf("%s", "ecbank_exchange_rate_"+time.Now().Format("2006-01-02"))
//		timeout  = 30 * time.Second
//	)
//	from := flag.String("from", "", "source currency, required")
//	to := flag.String("to", "EUR", "target currency")
//	flush := flag.Bool("flush", false, "flush cache if specified")
//
//	flag.Parse()
//
//	if *flush {
//		err := os.Remove(fileName)
//		if err != nil {
//			fmt.Printf("unexpected usage of flag -flush: %v", err)
//		}
//		return
//	}
//
//	args := flag.Args()
//	if len(args) < 1 {
//		_, _ = fmt.Fprintln(os.Stderr, "missing amount to convert")
//		flag.Usage()
//		os.Exit(1)
//	}
//	value := args[0]
//	fmt.Println(*from, *to, value)
//
//	quantity, err := money.ParseDecimal(value)
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, "unable to parse value %q: %s\n", value, err.Error())
//		os.Exit(1)
//	}
//
//	fromCurrency, err := money.ParseCurrency(*from)
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, "unable to parse from %q: %s\n", *from, err.Error())
//		os.Exit(1)
//	}
//
//	toCurrency, err := money.ParseCurrency(*to)
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, "unable to parse to %q: %s\n", *to, err.Error())
//		os.Exit(1)
//	}
//
//	amount, err := money.NewAmount(quantity, fromCurrency)
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, err.Error())
//		os.Exit(1)
//	}
//	fmt.Println(amount, toCurrency)
//
//	convertedAmount, err := money.Convert(amount, toCurrency, ecbank.NewEuroCentralBank(timeout, fileName))
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, "unable to convert %s to %s: %s.\n", amount, toCurrency, err.Error())
//		os.Exit(1)
//	}
//
//	fmt.Printf("%s = %s\n", amount, convertedAmount)
//}

//const maxAttempts = 6

//func main() {
//	lgr := newLogger()
//	path := "corpus/english.txt"
//	corpus, err := gordle.ReadCorpus(path)
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, "unable to read corpus: %s", err)
//	}
//	lgr.Infof("successfully read corpus from %s", path)
//
//	// Create the game
//	g, err := gordle.New(os.Stdin, corpus, maxAttempts)
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, "unable to start game: %s", err)
//	}
//	lgr.Infof("successfully started game")
//
//	// Run the game! It will end when it's over.
//	g.Play()
//	lgr.Infof("exited game")
//}
//
//func newLogger() *pocketlog.Logger {
//	lgr := pocketlog.New(pocketlog.LevelInfo)
//	return lgr
//}

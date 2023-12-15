package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	myFigure := figure.NewColorFigure("YAVUZLAR / FUZZER ", "", "green", true)
	myFigure.Print()                                                            // burada figlet yer alıyor
	wordlistFileName := flag.String("txt", "", "Wordlist file name")            // komut
	requestCount := flag.Int("s", 0, "Number of requests to make concurrently") // satırı
	targetURL := flag.String("u", "", "Target URL")                             // parametreleri

	flag.Parse()

	if *wordlistFileName == "" || *requestCount == 0 || *targetURL == "" { // kullanıcıdan alınan
		fmt.Println("Lütfen -txt, -s ve -u parametrelerini belirtin.") // parametreleri kontrol eder
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*wordlistFileName)
	if err != nil {
		fmt.Println("Wordlist dosyasını açma hatası:", err) // txt dosyasını açıp hata kontrolü yaptığım kısım
		os.Exit(1)
	}
	defer file.Close()

	var wg sync.WaitGroup             //eş zamanlı goroutine'leri beklemek için kullanılır.
	scanner := bufio.NewScanner(file) // dosyadan satır satır okuma yapılır.

	for scanner.Scan() {
		wordlist := strings.TrimSpace(scanner.Text())

		if wordlist != "" { //Her satırdaki kelimenin başında ve sonunda boşlukları temizleyerek
			wg.Add(1)                             //wordlist değişkenine atadım
			go func(wordlist, targetURL string) { //Boş olmayan bir kelime varsa
				defer wg.Done() //sync.WaitGroup'ı artırıp yeni bir goroutine başlattım
				//Burada goroutine içinde makeRequest fonksiyonu çağrılır.
				url := fmt.Sprintf("%s/%s", targetURL, wordlist)
				makeRequest(url, targetURL)
			}(wordlist, *targetURL)
		}
	}

	wg.Wait() //sync.WaitGroup tarafından beklenen tüm goroutine'lerin bitmesini için kullandım
}

// alt kısımda ise http.Get ile url e istek atıyorum
func makeRequest(url, targtURL string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("HTTP isteği hatası: %v\n", err)
		return
	}
	defer resp.Body.Close()
	//bu kısımda ise bulunan dizinleri ekrana yazdırıyorum
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Dizin bulundu! URL: %s\n", url)
	}
}

// --help için ekrana yazılacak olan ifadeler aşağıda yer alıyor
func printHelp() {
	fmt.Println("Kullanım: -txt [wordlist.txt] -s [34] -u [http://example.com]")
	fmt.Println("-txt : wordlist.txt ; kullanıcı dosya adı için girdi parametresi.")
	fmt.Println("-s   : eş zamanlılık için istenen sayının girdi parametresi.")
	fmt.Println("-u   : dizin tarama için hedef URL.")

}

// Flag.Parse() çağrıldığında bu fonksiyonun çalışmasını sağlamak için init fonksiyonu kullandım.
func init() {
	flag.Usage = printHelp //--help komutu çalıştırıldığında yardımcı fonksiyonun çağrılmasını sağlar.
}

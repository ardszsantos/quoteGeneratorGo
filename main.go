package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Estruturas de dados
type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Banco de dados simulado
var (
	substantivos = []string{"banano pível", "the xesque", "koskóvos", "ralubo", "xampu", "goflo", "afafué",
		"BORBO DUPLO", "calabreso", "googlepax", "apapua", "djó", "xepé", "spage spé", "pagai", "jope", "mijoto", "mejoto", "todo", "golo", "melontedonpeps", "xerecs", "bananotifane",
		"bananolove", "cebruntius", "ralugbo", "nobiru", "Equu", "ralube", "Prou", "todes", "redonde", "molom", "chiq", "chambs", "pope", "popGG", "rodimunosc", "mussum", "flope", "blog", "glup", "poulos",
		"papadopoulos", "poiés", "péu"}
	verbos = []string{"explodir", "voar", "gozar", "correr", "chorar", "moggar", "bolar", "goflar", "xoflar", "flambar", "melar", "blogar", "chepar", "negoçar", "calabrear", "rebuliçar",
		"malhar", "Rapear", "molhar", "biskinhar", "Janbular", "Tar", "biskinhando", "biskinhou", "bater", "piPoLutou"}
	adjetivos = []string{"goploso", "Doroso", "Garico", "famoso", "usado", "horrível", "Fuzul", "épico", "delício", "Ovo-cuzudo", "kanye west", "cabeça de rádio", "mulher homem", "redondo", "circular", "bolotoso",
		"Made in Ovo", "Jojo Piculo", "Danonésio", "neutre", "redonde", "bolotose"}
	lugares = []string{"no banheiro", "na cozinha", "na lua", "no shopping", "em Hogwarts", "no borBolhão", "no buzão topzera", "na casa do cabrunco", "em roney", "na Chiquelândia",
		"nas rotas apicuarias", "no corre do queijo", "na ruto", "no obelisco total", "no obelisco vip", "na bolotona", "no bolotão", "no privacy do playboy carti", "no cu de todes"}
	templates = []string{
		"Hoje eu vi {substantivo} {verbo} {lugar}. Foi {adjetivo} demais!",
		"{substantivo} {verbo} {lugar} e ninguém conseguiu parar. Foi simplesmente {adjetivo}.",
		"Dizem que {substantivo} {verbo} porque {substantivo} achou {lugar} {adjetivo}.",
		"Se {substantivo} pudesse {verbo} {lugar}, o mundo seria muito mais {adjetivo}.",
		"Nada é mais {adjetivo} do que {substantivo} {verbo} enquanto {substantivo} observa {lugar}.",
		"{substantivo} estava {verbo} {lugar}, até que {substantivo} apareceu e tudo ficou {adjetivo}.",
		"Quando {substantivo} começou a {verbo} {lugar}, todos perceberam como era {adjetivo}.",
		"{substantivo} tentou {verbo}, mas {substantivo} estragou tudo de um jeito {adjetivo}.",
		"A verdade sobre {substantivo} é que ele sempre quis {verbo} {lugar}.",
		"{substantivo} é conhecido por {verbo} {lugar}, mas poucos sabem como isso é {adjetivo}.",
		"Se não fosse por {substantivo}, {substantivo} nunca teria conseguido {verbo} {lugar}.",
		"{substantivo} e {substantivo} decidiram {verbo} juntos {lugar}. Foi {adjetivo} para todos.",
		"{substantivo} achou que {lugar} era {adjetivo}, então começou a {verbo}.",
		"{substantivo} estava {verbo} {lugar} quando {substantivo} apareceu e disse que era {adjetivo}.",
		"Depois de muito tempo, {substantivo} finalmente conseguiu {verbo} {lugar}. Foi {adjetivo}.",
		"Era uma vez {substantivo} que sonhava em {verbo} {lugar}. A vida era tão {adjetivo}.",
		"As lendas dizem que {substantivo} {verbo} {lugar}, e foi por isso que tudo ficou {adjetivo}.",
		"Se {substantivo} soubesse como {verbo} {lugar}, talvez as coisas fossem mais {adjetivo}.",
		"{substantivo} não sabia que {substantivo} {verbo} {lugar} era tão {adjetivo}.",
		"Todo mundo sabe que {substantivo} só quer {verbo} {lugar} para ser {adjetivo}.",
		"{substantivo} disse: 'Nada é tão {adjetivo} quanto {verbo} {lugar}'.",
		"A profecia dizia que {substantivo} {verbo} {lugar}, e ninguém imaginava que seria tão {adjetivo}.",
		"No fim, {substantivo} percebeu que {verbo} {lugar} não era tão {adjetivo} quanto parecia.",
		"{substantivo} tentou avisar que {substantivo} ia {verbo} {lugar}, mas ninguém acreditou.",
		"{substantivo} {verbo} {lugar}, e todos acharam isso muito {adjetivo}.",
		"Quando {substantivo} começou a {verbo} {lugar}, até {substantivo} ficou surpreso. Era {adjetivo}.",
		"{substantivo} tentou {verbo} {lugar}, mas {substantivo} impediu de um jeito {adjetivo}.",
		"{substantivo} achou que {verbo} {lugar} fosse {adjetivo}, mas estava errado.",
		"{substantivo} disse que {lugar} era {adjetivo}, e decidiu {verbo} ali mesmo.",
		"{substantivo} estava {verbo} {lugar}, quando {substantivo} chegou e disse que era {adjetivo}.",
		"Todo mundo ficou surpreso quando {substantivo} começou a {verbo} {lugar} de forma tão {adjetivo}.",
		"{substantivo} tentou avisar que era {adjetivo}, mas {substantivo} já estava {verbo} {lugar}.",
	}

	allQuotes = []Quote{}
)

// Funções auxiliares
func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func generateQuote() Quote {
	template := randomChoice(templates)

	// Mapeamento de placeholders únicos para valores
	placeholders := map[string][]string{
		"{substantivo}": substantivos,
		"{verbo}":       verbos,
		"{adjetivo}":    adjetivos,
		"{lugar}":       lugares,
	}

	// Substituições já realizadas
	usedValues := map[string]map[string]bool{}

	for key := range placeholders {
		usedValues[key] = make(map[string]bool)
	}

	replacerArgs := []string{}

	// Procura os placeholders no template e gera substituições únicas
	for key, values := range placeholders {
		count := strings.Count(template, key)
		for i := 0; i < count; i++ {
			value := randomChoice(values)
			// Garante valores únicos
			for usedValues[key][value] {
				value = randomChoice(values)
			}
			usedValues[key][value] = true
			replacerArgs = append(replacerArgs, key, value)
		}
	}

	replacer := strings.NewReplacer(replacerArgs...)

	return Quote{
		ID:     rand.Intn(1000),
		Author: "Jolhambão",
		Text:   replacer.Replace(template),
	}
}

func initQuotes() {
	// Gera 10 frases iniciais
	for i := 0; i < 10; i++ {
		allQuotes = append(allQuotes, generateQuote())
	}
}

// Handlers de endpoints
func quoteHandler(w http.ResponseWriter, r *http.Request) {
	quote := generateQuote()
	response := Response{Message: "Frase gerada com sucesso", Data: quote}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func quotesHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Lista de frases disponíveis", Data: allQuotes}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Servidor está funcionando!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Configuração do servidor
func main() {
	rand.Seed(time.Now().UnixNano())
	initQuotes()

	mux := http.NewServeMux()
	mux.HandleFunc("/quote", quoteHandler)   // Gera uma frase aleatória
	mux.HandleFunc("/quotes", quotesHandler) // Retorna todas as frases
	mux.HandleFunc("/health", healthHandler) // Health check

	// Middleware: Logging and CORS
	loggedMux := loggingMiddleware(corsMiddleware(mux))

	// Inicializa o servidor
	port := "8080"
	log.Printf("Servidor iniciado na porta %s", port)
	err := http.ListenAndServe(":"+port, loggedMux)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

// Middleware para habilitar CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configuração de CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permitir qualquer origem
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Responder imediatamente para requisições OPTIONS (preflight)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Chamar o próximo handler na cadeia
		next.ServeHTTP(w, r)
	})
}

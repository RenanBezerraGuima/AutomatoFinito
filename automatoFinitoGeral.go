package main

import (
	"fmt"
	"slices"
)

type AutomatoFinito struct {
	Estados       []string
	Alfabeto      []rune                       // vetor de caracters (rune's)
	Transicoes    map[string]map[rune][]string // estadoOrigem: [símbolo: [estadoDestino]]
	EstadoInicial string
	EstadosFinais []string
	Cadeia        []rune // vetor de caracters (rune's)
}

func (AF *AutomatoFinito) adicionarEstado(estado string) {
	AF.Estados = append(AF.Estados, estado)
}

func (AF *AutomatoFinito) adicionarAlfabeto(simbolo rune) {
	AF.Alfabeto = append(AF.Alfabeto, simbolo)
}

func (AF *AutomatoFinito) adicionarTransicao(estadoOrigem string, simbolo rune, estadoDestino string) {
	if AF.Transicoes == nil {
		AF.Transicoes = make(map[string]map[rune][]string)
	}
	if AF.Transicoes[estadoOrigem] == nil {
		AF.Transicoes[estadoOrigem] = make(map[rune][]string)
	}
	AF.Transicoes[estadoOrigem][simbolo] = append(AF.Transicoes[estadoOrigem][simbolo], estadoDestino)
}

// epsilonClosure retorna o conjunto de estados alcançáveis a partir de um conjunto de estados, seguindo apenas transições épsilon.
func (AF *AutomatoFinito) epsilonClosure(estados []string) []string {
	closure := make(map[string]bool)
	for _, estado := range estados {
		closure[estado] = true
	}

	pilha := make([]string, len(estados))
	copy(pilha, estados)

	for len(pilha) > 0 {
		estadoAtual := pilha[len(pilha)-1]
		pilha = pilha[:len(pilha)-1]

		if transicoesEstado, ok := AF.Transicoes[estadoAtual]; ok {
			if destinosEpsilon, ok := transicoesEstado['ε']; ok {
				for _, destino := range destinosEpsilon {
					if !closure[destino] {
						closure[destino] = true
						pilha = append(pilha, destino)
					}
				}
			}
		}
	}

	resultado := make([]string, 0, len(closure))
	for estado := range closure {
		resultado = append(resultado, estado)
	}
	return resultado
}

func (AF *AutomatoFinito) adicionarEstadoInicial(estadoInicial string) {
	AF.EstadoInicial = estadoInicial
}

func (AF *AutomatoFinito) adicionarEstadoFinal(estadoFinal string) {
	AF.EstadosFinais = append(AF.EstadosFinais, estadoFinal)
}

func (AF *AutomatoFinito) adicionarCadeia(cadeia string) {
	AF.Cadeia = []rune(cadeia)
}

func (AF *AutomatoFinito) funcionamento() bool {
	estadosAtuais := AF.epsilonClosure([]string{AF.EstadoInicial})

	for _, simbolo := range AF.Cadeia {
		proximosEstados := make(map[string]bool)
		for _, estado := range estadosAtuais {
			if transicoesEstado, ok := AF.Transicoes[estado]; ok {
				if destinos, ok := transicoesEstado[simbolo]; ok {
					for _, destino := range destinos {
						proximosEstados[destino] = true
					}
				}
			}
		}

		if len(proximosEstados) == 0 {
			return false // Sem transições para o símbolo atual
		}

		// Converter mapa para slice para epsilonClosure
		sliceProximosEstados := make([]string, 0, len(proximosEstados))
		for estado := range proximosEstados {
			sliceProximosEstados = append(sliceProximosEstados, estado)
		}
		estadosAtuais = AF.epsilonClosure(sliceProximosEstados)
	}

	// Verificar se algum dos estados atuais é final
	for _, estado := range estadosAtuais {
		if slices.Contains(AF.EstadosFinais, estado) {
			return true
		}
	}
	return false
}

func exemplo() {
	fmt.Println("Exemplo de Automato Finito Não-Determinístico (AFN):")
	fmt.Println("Este AFN aceita cadeias que terminam com \"ab\".")
	fmt.Println("Estados: q0, q1, q2")
	fmt.Println("Alfabeto: a, b")
	fmt.Println("Transições:")
	fmt.Println("  q0, a -> q0")
	fmt.Println("  q0, b -> q0")
	fmt.Println("  q0, a -> q1")
	fmt.Println("  q1, b -> q2")
	fmt.Println("Estado Inicial: q0")
	fmt.Println("Estado Final: q2")
	fmt.Println("------------------------------------")

	AFN := AutomatoFinito{}
	AFN.adicionarEstado("q0")
	AFN.adicionarEstado("q1")
	AFN.adicionarEstado("q2")

	AFN.adicionarAlfabeto('a')
	AFN.adicionarAlfabeto('b')

	// Transições que mantêm em q0
	AFN.adicionarTransicao("q0", 'a', "q0")
	AFN.adicionarTransicao("q0", 'b', "q0")
	// Transição para o possível início do padrão "ab"
	AFN.adicionarTransicao("q0", 'a', "q1")
	// Transição que completa o padrão "ab"
	AFN.adicionarTransicao("q1", 'b', "q2")

	AFN.adicionarEstadoInicial("q0")
	AFN.adicionarEstadoFinal("q2")

	testes := []struct {
		cadeia  string
		esperado bool
	}{
		{"ab", true},
		{"aab", true},
		{"bab", true},
		{"aaab", true},
		{"b", false},
		{"a", false},
		{"aba", false},
		{"", false}, // Cadeia vazia
	}

	for _, teste := range testes {
		AFN.adicionarCadeia(teste.cadeia)
		resultado := AFN.funcionamento()
		fmt.Printf("Cadeia \"%s\" -> ", teste.cadeia)
		if resultado {
			fmt.Print("aceita")
		} else {
			fmt.Print("não aceita")
		}
		if resultado == teste.esperado {
			fmt.Println(" (Correto)")
		} else {
			fmt.Println(" (Incorreto)")
		}
	}

	// Exemplo com transição épsilon
	fmt.Println("\nExemplo de AFN com transição épsilon:")
	fmt.Println("Este AFN aceita \"a*b\" (zero ou mais 'a's seguidos por um 'b').")
	// q0 --ε--> q1 --a--> q1 --b--> q2
	// Estados: qe0, qe1, qe2
	// Alfabeto: a, b (ε é implícito)
	// Transições:
	//   qe0, ε -> qe1
	//   qe1, a -> qe1
	//   qe1, b -> qe2
	// Estado Inicial: qe0
	// Estado Final: qe2
	fmt.Println("------------------------------------")

	AFNepsilon := AutomatoFinito{}
	AFNepsilon.adicionarEstado("qe0")
	AFNepsilon.adicionarEstado("qe1")
	AFNepsilon.adicionarEstado("qe2")

	AFNepsilon.adicionarAlfabeto('a')
	AFNepsilon.adicionarAlfabeto('b')
	// Não adicionamos 'ε' ao alfabeto visível, mas usamos nas transições

	AFNepsilon.adicionarTransicao("qe0", 'ε', "qe1") // Transição épsilon
	AFNepsilon.adicionarTransicao("qe1", 'a', "qe1")
	AFNepsilon.adicionarTransicao("qe1", 'b', "qe2")

	AFNepsilon.adicionarEstadoInicial("qe0")
	AFNepsilon.adicionarEstadoFinal("qe2")

	testesEpsilon := []struct {
		cadeia  string
		esperado bool
	}{
		{"b", true},   // ε -> q1, b -> q2
		{"ab", true},  // ε -> q1, a -> q1, b -> q2
		{"aab", true}, // ε -> q1, a -> q1, a -> q1, b -> q2
		{"aaab", true},
		{"", false},   // Não consome 'b'
		{"a", false},  // Não consome 'b'
		{"ba", false}, // 'b' primeiro não é permitido pela lógica qe0->qe1
	}

	for _, teste := range testesEpsilon {
		AFNepsilon.adicionarCadeia(teste.cadeia)
		resultado := AFNepsilon.funcionamento()
		fmt.Printf("Cadeia \"%s\" -> ", teste.cadeia)
		if resultado {
			fmt.Print("aceita")
		} else {
			fmt.Print("não aceita")
		}
		if resultado == teste.esperado {
			fmt.Println(" (Correto)")
		} else {
			fmt.Println(" (Incorreto)")
		}
	}
}

func leituraEstados(AFUsuario *AutomatoFinito) {
	fmt.Println("Digite os estados (um por vez). Digite \"fim\" para encerrar:")
	for {
		var estado string
		fmt.Print("> ")
		fmt.Scan(&estado)
		if estado == "fim" {
			return
		}
		if estado == "" {
			fmt.Println("Erro: Nome do estado não pode ser vazio. Tente novamente.")
			continue // Pede novo estado
		}
		if slices.Contains(AFUsuario.Estados, estado) {
			fmt.Printf("Erro: Estado '%s' já foi adicionado. Tente outro.\n", estado)
			continue // Pede novo estado
		}
		AFUsuario.adicionarEstado(estado)
		fmt.Printf("Estado '%s' adicionado.\n", estado)
	}
}

func leituraEstadoInicial(AFUsuario *AutomatoFinito) bool {
	fmt.Print("Digite o estado inicial: ")
	var estado string
	fmt.Scan(&estado)
	if !slices.Contains(AFUsuario.Estados, estado) {
		fmt.Println("Erro: Estado inicial não está presente nos estados adicionados")
		return false
	}
	AFUsuario.adicionarEstadoInicial(estado)
	return true
}

func leituraAlfabeto(AFUsuario *AutomatoFinito) bool {
	fmt.Println("Digite o alfabeto (um símbolo por vez). Digite \"fim\" para encerrar:")
	for {
		var entrada string
		fmt.Print("> ")
		fmt.Scan(&entrada)
		if entrada == "fim" {
			return true
		}
		r := []rune(entrada)
		if len(r) != 1 {
			fmt.Println("Erro: insira exatamente um caractere.")
			continue // Pede novo símbolo
		}
		simbolo := r[0]
		if slices.Contains(AFUsuario.Alfabeto, simbolo) {
			fmt.Printf("Erro: Símbolo '%c' já foi adicionado ao alfabeto. Tente outro.\n", simbolo)
			continue // Pede novo símbolo
		}
		AFUsuario.adicionarAlfabeto(simbolo)
		fmt.Printf("Símbolo '%c' adicionado ao alfabeto.\n", simbolo)
	}
}

const epsilonRune = 'ε'

func leituraTransicoes(AFUsuario *AutomatoFinito) bool {
	fmt.Println("\n--- Adicionar Transições ---")
	fmt.Println("Para cada transição, primeiro o estado de origem e o símbolo.")
	fmt.Println("Depois, digite os estados de destino um por vez.")
	fmt.Println("Digite \"fim\" como estado de destino para finalizar a adição de destinos para a transição atual.")
	fmt.Println("Digite \"fim\" como origem para encerrar a adição de todas as transições.")

	for {
		var origem string
		fmt.Print("\nOrigem (ou \"fim\" para encerrar tudo): ")
		fmt.Scan(&origem)
		if origem == "fim" {
			return true
		}

		if !slices.Contains(AFUsuario.Estados, origem) {
			fmt.Printf("Erro: Estado de origem '%s' não existe.\n", origem)
			continue // Pede nova origem
		}

		var simboloStr string
		fmt.Printf("Símbolo para %s (ou \"eps\"/\"epsilon\" para épsilon): ", origem)
		fmt.Scan(&simboloStr)

		var simbolo rune
		isEpsilon := false
		if simboloStr == "eps" || simboloStr == "epsilon" {
			simbolo = epsilonRune
			isEpsilon = true
		} else {
			r := []rune(simboloStr)
			if len(r) != 1 {
				fmt.Println("Erro: símbolo deve ser um único caractere ou \"eps\"/\"epsilon\".")
				continue // Pede nova origem
			}
			simbolo = r[0]
			if !slices.Contains(AFUsuario.Alfabeto, simbolo) {
				fmt.Printf("Erro: Símbolo '%c' não presente no alfabeto.\n", simbolo)
				continue // Pede nova origem
			}
		}

		fmt.Printf("Adicionando destinos para (%s, %q):\n", origem, simbolo)
		for {
			var destino string
			fmt.Printf("  Destino para %s,%q (ou \"fim\" para esta transição): ", origem, simbolo)
			fmt.Scan(&destino)
			if destino == "fim" {
				break // Finaliza destinos para esta transição (origem, simbolo)
			}

			if !slices.Contains(AFUsuario.Estados, destino) {
				fmt.Printf("Erro: Estado de destino '%s' não existe. Tente novamente.\n", destino)
				continue // Pede novo destino para a mesma (origem, simbolo)
			}

			AFUsuario.adicionarTransicao(origem, simbolo, destino)
			fmt.Printf("    Adicionado: %s --%q--> %s\n", origem, simbolo, destino)
		}
		fmt.Println("Próxima transição.")
	}
}

func leituraEstadosFinais(AFUsuario *AutomatoFinito) bool {
	fmt.Println("Digite os estados finais (um por vez). Digite \"fim\" para encerrar:")
	for {
		var estado string
		fmt.Print("> ")
		fmt.Scan(&estado)
		if estado == "fim" {
			return true
		}
		if !slices.Contains(AFUsuario.Estados, estado) {
			fmt.Println("Erro: Estado final não está na lista de estados definidos.")
			continue // Pede novo estado final
		}
		if slices.Contains(AFUsuario.EstadosFinais, estado) {
			fmt.Printf("Erro: Estado '%s' já foi adicionado como final. Tente outro.\n", estado)
			continue // Pede novo estado final
		}
		AFUsuario.adicionarEstadoFinal(estado)
		fmt.Printf("Estado final '%s' adicionado.\n", estado)
	}
}

func exibicaoAutomato(AFUsuario *AutomatoFinito) {
	fmt.Println("\nAutômato criado:")
	fmt.Printf("Estado Inicial: %s\n", AFUsuario.EstadoInicial)
	fmt.Printf("Estados: %v\n", AFUsuario.Estados)
	fmt.Printf("Alfabeto: %q\n", AFUsuario.Alfabeto)
	fmt.Println("Transições:")
	for estado, m := range AFUsuario.Transicoes {
		for simbolo, destinos := range m {
			for _, destino := range destinos {
				fmt.Printf("%s,%q --> %s\n", estado, simbolo, destino)
			}
		}
	}
	fmt.Printf("Estados Finais: %v\n", AFUsuario.EstadosFinais)
}

func testeCadeiasUsuario(AFUsuario *AutomatoFinito) {
	fmt.Println("Digite a cadeia para testar (ou \"sair\" para encerrar):")
	for {
		var entrada string
		fmt.Print("> ")
		fmt.Scan(&entrada)
		if entrada == "sair" {
			return
		}
		AFUsuario.adicionarCadeia(entrada)
		if AFUsuario.funcionamento() {
			fmt.Println("Cadeia aceita")
		} else {
			fmt.Println("Cadeia não aceita")
		}
	}
}

func automatoUsuario() {
	fmt.Println("\n==== Crie seu autômato ====")
	AFUsuario := AutomatoFinito{}

	leituraEstados(&AFUsuario)
	if !leituraEstadoInicial(&AFUsuario) {
		return
	}
	if !leituraAlfabeto(&AFUsuario) {
		return
	}
	if !leituraTransicoes(&AFUsuario) {
		return
	}
	if !leituraEstadosFinais(&AFUsuario) {
		return
	}

	exibicaoAutomato(&AFUsuario)
	testeCadeiasUsuario(&AFUsuario)
}

func main() {
	for {
		fmt.Println("\nMenu Principal:")
		fmt.Println("1. Rodar Exemplo Pré-definido")
		fmt.Println("2. Criar Novo Autômato")
		fmt.Println("3. Sair")
		fmt.Print("Escolha uma opção: ")

		var escolha int
		// Tentativa de ler o inteiro. Ignora entradas não numéricas por enquanto.
		_, err := fmt.Scanln(&escolha)
		if err != nil {
			fmt.Println("Entrada inválida. Por favor, digite um número.")
			// Limpar o buffer de entrada em caso de erro, para evitar loop infinito se a entrada inválida não for consumida.
			// Esta é uma forma simples de tentar limpar. Pode não ser perfeita para todos os casos.
			var temp string
			fmt.Scanln(&temp) // Tenta consumir o resto da linha inválida.
			continue
		}

		switch escolha {
		case 1:
			exemplo()
		case 2:
			automatoUsuario()
		case 3:
			fmt.Println("Encerrando o programa.")
			return // Sai da função main e, portanto, do programa
		default:
			fmt.Println("Opção inválida, tente novamente.")
		}
		fmt.Println("------------------------------------") // Separador visual
	}
}

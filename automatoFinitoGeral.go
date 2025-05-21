package main

import (
	"fmt"
	"slices"
)

type AutomatoFinito struct {
	Estados       []string
	Alfabeto      []rune                     // vetor de caracters (rune's)
	Transicoes    map[string]map[rune]string // estadoOrigem: [símbolo: estadoDestino]
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
		AF.Transicoes = make(map[string]map[rune]string)
	}
	if AF.Transicoes[estadoOrigem] == nil {
		AF.Transicoes[estadoOrigem] = make(map[rune]string)
	}
	AF.Transicoes[estadoOrigem][simbolo] = estadoDestino
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
	estadoAtual := AF.EstadoInicial
	for _, simbolo := range AF.Cadeia {
		if proximo, ok := AF.Transicoes[estadoAtual][simbolo]; ok {
			estadoAtual = proximo
		} else {
			return false
		}
	}
	return slices.Contains(AF.EstadosFinais, estadoAtual)
}

func exemplo() {
	fmt.Println("Autômato Finito Exemplo:")
	fmt.Println("Aceita a cadeia ab e nada mais")
	fmt.Println("Linguagem = {ab}")
	fmt.Println("q0,a -> q1,b -> q2")

	AF := AutomatoFinito{}
	AF.adicionarEstado("q0")
	AF.adicionarEstado("q1")
	AF.adicionarEstado("q2")
	AF.adicionarAlfabeto('a')
	AF.adicionarAlfabeto('b')
	AF.adicionarTransicao("q0", 'a', "q1")
	AF.adicionarTransicao("q1", 'b', "q2")
	AF.adicionarEstadoInicial("q0")
	AF.adicionarEstadoFinal("q2")

	AF.adicionarCadeia("ab")
	fmt.Print("Cadeia \"ab\" -> ")
	if AF.funcionamento() {
		fmt.Println("aceita")
	} else {
		fmt.Println("não aceita")
	}

	AF.adicionarCadeia("ba")
	fmt.Print("Cadeia \"ba\" -> ")
	if AF.funcionamento() {
		fmt.Println("aceita")
	} else {
		fmt.Println("não aceita")
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
		AFUsuario.adicionarEstado(estado)
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
			return false
		}
		AFUsuario.adicionarAlfabeto(r[0])
	}
}

func leituraTransicoes(AFUsuario *AutomatoFinito) bool {
	fmt.Println("Digite as transições no formato: origem símbolo destino. Digite \"fim\" como origem para encerrar:")
	for {
		var origem, simboloStr, destino string
		fmt.Print("> ")
		fmt.Scan(&origem)
		if origem == "fim" {
			return true
		}
		fmt.Scan(&simboloStr, &destino)

		r := []rune(simboloStr)
		if len(r) != 1 {
			fmt.Println("Erro: símbolo deve ser um único caractere")
			return false
		}
		simbolo := r[0]

		if !slices.Contains(AFUsuario.Estados, origem) || !slices.Contains(AFUsuario.Estados, destino) {
			fmt.Println("Erro: Estado de origem ou destino não existentes")
			return false
		}
		if !slices.Contains(AFUsuario.Alfabeto, simbolo) {
			fmt.Println("Erro: Símbolo não presente no alfabeto")
			return false
		}

		AFUsuario.adicionarTransicao(origem, simbolo, destino)
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
			fmt.Println("Erro: Estado não presente")
			return false
		}
		AFUsuario.adicionarEstadoFinal(estado)
	}
}

func exibicaoAutomato(AFUsuario *AutomatoFinito) {
	fmt.Println("\nAutômato criado:")
	fmt.Printf("Estado Inicial: %s\n", AFUsuario.EstadoInicial)
	fmt.Printf("Estados: %v\n", AFUsuario.Estados)
	fmt.Printf("Alfabeto: %q\n", AFUsuario.Alfabeto)
	fmt.Println("Transições:")
	for estado, m := range AFUsuario.Transicoes {
		for simbolo, destino := range m {
			fmt.Printf("%s,%q --> %s\n", estado, simbolo, destino)
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
	exemplo()
	automatoUsuario()
}

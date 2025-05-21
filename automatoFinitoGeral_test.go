package main

import (
	"slices"
	"testing"
)

// slicesEqualIgnoringOrderAndDuplicates checks if two string slices are equal,
// ignoring the order of elements and treating them as sets (duplicates in input don't matter for comparison).
func slicesEqualIgnoringOrderAndDuplicates(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}

	mapA := make(map[string]bool)
	for _, item := range a {
		mapA[item] = true
	}

	mapB := make(map[string]bool)
	for _, item := range b {
		mapB[item] = true
	}

	if len(mapA) != len(mapB) {
		return false
	}

	for key := range mapA {
		if !mapB[key] {
			return false
		}
	}
	return true
}

func TestEpsilonClosure(t *testing.T) {
	tests := []struct {
		name           string
		af             AutomatoFinito
		initialStates  []string
		expectedStates []string
	}{
		{
			name: "No epsilon transitions",
			af: AutomatoFinito{
				Estados:    []string{"q0", "q1"},
				Transicoes: map[string]map[rune][]string{},
			},
			initialStates:  []string{"q0"},
			expectedStates: []string{"q0"},
		},
		{
			name: "Simple chain q0-e->q1-e->q2",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
					"q1": {'ε': {"q2"}},
				},
			},
			initialStates:  []string{"q0"},
			expectedStates: []string{"q0", "q1", "q2"},
		},
		{
			name: "Simple chain starting from middle q0-e->q1-e->q2, start q1",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
					"q1": {'ε': {"q2"}},
				},
			},
			initialStates:  []string{"q1"},
			expectedStates: []string{"q1", "q2"},
		},
		{
			name: "Branching epsilon transitions q0-e->q1, q0-e->q2",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1", "q2"}},
				},
			},
			initialStates:  []string{"q0"},
			expectedStates: []string{"q0", "q1", "q2"},
		},
		{
			name: "Cycles q0-e->q1, q1-e->q0",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
					"q1": {'ε': {"q0"}},
				},
			},
			initialStates:  []string{"q0"},
			expectedStates: []string{"q0", "q1"},
		},
		{
			name: "Cycles with a tail q0-e->q1, q1-e->q0, q1-e->q2",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
					"q1": {'ε': {"q0", "q2"}},
				},
			},
			initialStates:  []string{"q0"},
			expectedStates: []string{"q0", "q1", "q2"},
		},
		{
			name: "Disconnected components with epsilon, start q0 (reaches q1), q2 no eps",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
				},
			},
			initialStates:  []string{"q0"},
			expectedStates: []string{"q0", "q1"},
		},
		{
			name: "Disconnected components, start q2 (no eps from q2)",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
				},
			},
			initialStates:  []string{"q2"},
			expectedStates: []string{"q2"},
		},
		{
			name: "Multiple initial states with shared epsilon transitions",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2", "q3"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q2"}},
					"q1": {'ε': {"q2"}},
					"q2": {'ε': {"q3"}},
				},
			},
			initialStates:  []string{"q0", "q1"},
			expectedStates: []string{"q0", "q1", "q2", "q3"},
		},
		{
			name: "Initial states include a state reachable by epsilon from another initial",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1", "q2"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
					"q1": {'ε': {"q2"}},
				},
			},
			initialStates:  []string{"q0", "q1"},
			expectedStates: []string{"q0", "q1", "q2"},
		},
		{
			name: "No initial states",
			af: AutomatoFinito{
				Estados: []string{"q0", "q1"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}},
				},
			},
			initialStates:  []string{},
			expectedStates: []string{},
		},
		{
			name: "Epsilon to non-existent state (should be handled gracefully by AF structure, closure just explores)",
			af: AutomatoFinito{
				Estados: []string{"q0"},
				Transicoes: map[string]map[rune][]string{
					"q0": {'ε': {"q1"}}, // q1 not in AF.Estados, but epsilonClosure should still work
				},
			},
			initialStates:  []string{"q0"},
			expectedStates: []string{"q0", "q1"}, // q1 is "explored" conceptually by epsilonClosure
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure Transicoes is not nil if it was meant to be empty but addressable
			if tt.af.Transicoes == nil {
				tt.af.Transicoes = make(map[string]map[rune][]string)
			}
			// Ensure nested maps are not nil for states mentioned in initialStates
			// This mimics how adicionarTransicao would initialize them
			for _, initState := range tt.initialStates {
				if _, ok := tt.af.Transicoes[initState]; !ok && tt.af.Transicoes != nil {
					// This part is more about setting up the AF state correctly for the test
					// if initialStates are part of states that *could* have transitions.
					// If a state has no outgoing transitions, its entry in Transicoes might be nil or missing.
					// epsilonClosure should handle this fine.
				}
			}


			gotStates := tt.af.epsilonClosure(tt.initialStates)
			if !slicesEqualIgnoringOrderAndDuplicates(gotStates, tt.expectedStates) {
				t.Errorf("epsilonClosure(%v) = %v, want %v", tt.initialStates, gotStates, tt.expectedStates)
			}
		})
	}
}

func TestFuncionamentoNFA(t *testing.T) {
	// NFA1: accepts strings ending with "ab"
	nfa1 := AutomatoFinito{
		Estados:       []string{"q0", "q1", "q2"},
		Alfabeto:      []rune{'a', 'b'},
		Transicoes:    map[string]map[rune][]string{
			"q0": {
				'a': {"q0", "q1"},
				'b': {"q0"},
			},
			"q1": {
				'b': {"q2"},
			},
		},
		EstadoInicial: "q0",
		EstadosFinais: []string{"q2"},
	}

	// NFA2: accepts "a*b" using epsilon transitions (q_start --ε--> q_a_loop --a--> q_a_loop --ε--> q_b_trans --b--> q_final)
	nfa2 := AutomatoFinito{
		Estados:       []string{"q_start", "q_a_loop", "q_b_trans", "q_final"},
		Alfabeto:      []rune{'a', 'b'},
		Transicoes:    map[string]map[rune][]string{
			"q_start":  {'ε': {"q_a_loop"}},
			"q_a_loop": {
				'a': {"q_a_loop"},
				'ε': {"q_b_trans"},
			},
			"q_b_trans": {'b': {"q_final"}},
		},
		EstadoInicial: "q_start",
		EstadosFinais: []string{"q_final"},
	}
	
	// NFA3: accepts (a|b)*a --- language ends with 'a'
	nfa3 := AutomatoFinito{
		Estados:       []string{"S", "A"}, // S = initial, A = final (accepts 'a')
		Alfabeto:      []rune{'a', 'b'},
		Transicoes:    map[string]map[rune][]string{
			"S": {
				'a': {"S", "A"}, // On 'a', can stay in S or go to A
				'b': {"S"},      // On 'b', stay in S
			},
			// No transitions from A, or A could loop to itself if language was (a|b)*a(a|b)*
		},
		EstadoInicial: "S",
		EstadosFinais: []string{"A"},
	}


	tests := []struct {
		name     string
		af       *AutomatoFinito
		cadeia   string
		expected bool
	}{
		// Tests for nfa1 (ends with "ab")
		{"NFA1_ab", &nfa1, "ab", true},
		{"NFA1_aab", &nfa1, "aab", true},
		{"NFA1_bab", &nfa1, "bab", true},
		{"NFA1_cab", &nfa1, "cab", true}, // 'c' not in alphabet, handled by NFA logic (no transition)
		{"NFA1_b", &nfa1, "b", false},
		{"NFA1_a", &nfa1, "a", false},
		{"NFA1_acb_c_not_in_alphabet", &nfa1, "acb", false},
		{"NFA1_empty", &nfa1, "", false},
		{"NFA1_aaab", &nfa1, "aaab", true},
		{"NFA1_baba", &nfa1, "baba", false},

		// Tests for nfa2 (a*b with epsilon)
		{"NFA2_b", &nfa2, "b", true},     // ε -> q_a_loop, ε -> q_b_trans, b -> q_final
		{"NFA2_ab", &nfa2, "ab", true},   // ε -> q_a_loop, a -> q_a_loop, ε -> q_b_trans, b -> q_final
		{"NFA2_aab", &nfa2, "aab", true},
		{"NFA2_aaab", &nfa2, "aaab", true},
		{"NFA2_acb_c_not_in_alphabet", &nfa2, "acb", false},
		{"NFA2_ba", &nfa2, "ba", false}, // Cannot start with b effectively before a's loop
		{"NFA2_empty", &nfa2, "", false}, // Requires 'b'
		{"NFA2_a", &nfa2, "a", false},    // Requires 'b'

		// Tests for nfa3 ((a|b)*a)
		{"NFA3_a", &nfa3, "a", true},
		{"NFA3_ba", &nfa3, "ba", true},
		{"NFA3_bba", &nfa3, "bba", true},
		{"NFA3_aa", &nfa3, "aa", true},
		{"NFA3_empty", &nfa3, "", false},
		{"NFA3_b", &nfa3, "b", false},
		{"NFA3_ab", &nfa3, "ab", false}, // Ends with b
		{"NFA3_bab", &nfa3, "bab", false}, // Ends with b
		{"NFA3_bb", &nfa3, "bb", false}, // Ends with b
		{"NFA3_ca_c_not_in_alphabet", &nfa3, "ca", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.af.Transicoes == nil { // Ensure map is initialized for safety, though struct literals do this.
				tt.af.Transicoes = make(map[string]map[rune][]string)
			}
			tt.af.adicionarCadeia(tt.cadeia)
			got := tt.af.funcionamento()
			if got != tt.expected {
				t.Errorf("Automato %s com cadeia \"%s\": got %v, want %v", tt.name, tt.cadeia, got, tt.expected)
			}
		})
	}
}

func TestFuncionamentoDFA(t *testing.T) {
	dfa1 := AutomatoFinito{
		Estados:       []string{"q0", "q1", "q2"},
		Alfabeto:      []rune{'a', 'b'},
		Transicoes:    map[string]map[rune][]string{
			"q0": {'a': {"q1"}},
			"q1": {'b': {"q2"}},
		},
		EstadoInicial: "q0",
		EstadosFinais: []string{"q2"},
	}

	dfa2 := AutomatoFinito{
		Estados:       []string{"q_even", "q_odd"},
		Alfabeto:      []rune{'a'},
		Transicoes:    map[string]map[rune][]string{
			"q_even": {'a': {"q_odd"}},
			"q_odd":  {'a': {"q_even"}},
		},
		EstadoInicial: "q_even",
		EstadosFinais: []string{"q_even"},
	}

	tests := []struct {
		name     string
		af       *AutomatoFinito
		cadeia   string
		expected bool
	}{
		// Tests for dfa1 (accepts "ab" only)
		{"DFA1_ab", &dfa1, "ab", true},
		{"DFA1_a", &dfa1, "a", false},
		{"DFA1_b", &dfa1, "b", false},
		{"DFA1_aba", &dfa1, "aba", false},
		{"DFA1_empty", &dfa1, "", false},
		{"DFA1_toolong", &dfa1, "abb", false},
		{"DFA1_wrongchar", &dfa1, "ac", false}, // Symbol 'c' not in alphabet

		// Tests for dfa2 (accepts even number of 'a's)
		{"DFA2_empty", &dfa2, "", true},
		{"DFA2_a", &dfa2, "a", false},
		{"DFA2_aa", &dfa2, "aa", true},
		{"DFA2_aaa", &dfa2, "aaa", false},
		{"DFA2_aaaa", &dfa2, "aaaa", true},
		{"DFA2_b_not_in_alphabet", &dfa2, "b", false}, // 'b' not in alphabet, should fail
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure Transicoes and nested maps are initialized if they were defined as nil
			// This is more for safety in test setup, though direct struct literal usually handles it.
			if tt.af.Transicoes == nil {
				tt.af.Transicoes = make(map[string]map[rune][]string)
			}
			// It's important that tt.af.Alfabeto is correctly set for the funcionamento logic,
			// especially if it relies on checking symbol existence in the alphabet
			// (though the current NFA `funcionamento` doesn't explicitly, DFAs often do).

			tt.af.adicionarCadeia(tt.cadeia) // Set the Cadeia field
			got := tt.af.funcionamento()
			if got != tt.expected {
				t.Errorf("Automato %s com cadeia \"%s\": got %v, want %v", tt.name, tt.cadeia, got, tt.expected)
			}
		})
	}
}
